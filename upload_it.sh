#!/bin/bash
set -e

# Variables
export NAMESPACE="sealights-demo"
export DEPLOYMENT_NAME="sealights-go-demo"
export PORT=8080
export DOMAIN="redhat.sealights.co"
export SEALIGHTS_AGENT_TOKEN="${SEALIGHTS_AGENT_TOKEN:-""}"
export KUBECONFIG="${KUBECONFIG:-""}"
export BUILD_SESSION_ID="${BUILD_SESSION_ID:-""}"

# Ensure the namespace exists
echo "INFO: Ensuring namespace $NAMESPACE exists..."
oc get namespace "$NAMESPACE" >/dev/null 2>&1 || oc create namespace "$NAMESPACE"

# Create Deployment
echo "INFO: Creating deployment..."
oc apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: $DEPLOYMENT_NAME
  namespace: $NAMESPACE
spec:
  replicas: 1
  selector:
    matchLabels:
      app: $DEPLOYMENT_NAME
  template:
    metadata:
      labels:
        app: $DEPLOYMENT_NAME
    spec:
      containers:
      - name: $DEPLOYMENT_NAME
        image: $IMAGE
        env:
          - name: SEALIGHTS_TOKEN
            value: $SEALIGHTS_AGENT_TOKEN
        ports:
        - containerPort: $PORT
EOF

oc get deployment -n "$NAMESPACE" "$DEPLOYMENT_NAME" -o yaml

# Check if Deployment was created
echo "INFO: Verifying deployment status..."
TIMEOUT=120
timeout $TIMEOUT oc rollout status deployment/$DEPLOYMENT_NAME -n $NAMESPACE || {
  echo "ERROR: Deployment failed or timed out after $TIMEOUT seconds!"
  POD_NAME=$(oc get pods -n $NAMESPACE -l app=$DEPLOYMENT_NAME -o jsonpath='{.items[0].metadata.name}')
  echo "INFO: Fetching logs for pod $POD_NAME..."
  oc logs -n $NAMESPACE $POD_NAME || echo "ERROR: Failed to fetch pod logs."
  exit 1
}

# Create a Service
echo "INFO: Creating service..."
oc expose deployment "$DEPLOYMENT_NAME" \
  --port="$PORT" \
  --target-port="$PORT" \
  --name="$DEPLOYMENT_NAME" \
  -n "$NAMESPACE" || { echo "ERROR: Service creation failed!"; exit 1; }

# Create Route
echo "INFO: Creating route..."
oc create route edge --service="$DEPLOYMENT_NAME" --insecure-policy=Redirect -n "$NAMESPACE" || { echo "ERROR: Route creation failed!"; exit 1; }

# Fetch Route URL
export CONTAINER_ROUTE_URL="https://$(oc get route "$DEPLOYMENT_NAME" -n "$NAMESPACE" -o jsonpath='{.spec.host}')"
echo "Application is accessible at: $CONTAINER_ROUTE_URL"

# Install Ginkgo
go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Create a Sealights test session
echo "INFO: Creating Sealights test session..."
TEST_SESSION_ID=$(curl -X POST "https://$DOMAIN/sl-api/v1/test-sessions" \
  -H "Authorization: Bearer $SEALIGHTS_AGENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"labId":"","testStage":"go-calc-integration","bsid":"'${BUILD_SESSION_ID}'","sessionTimeout":10000}' | jq -r '.data.testSessionId')

if [ -n "$TEST_SESSION_ID" ]; then
  echo "Test session ID: $TEST_SESSION_ID"
  export TEST_SESSION_ID
else
  echo "Failed to retrieve test session ID"
  exit 1
fi

# Fetch excluded tests
RESPONSE=$(curl -X GET "https://$DOMAIN/sl-api/v2/test-sessions/$TEST_SESSION_ID/exclude-tests" \
  -H "Authorization: Bearer $SEALIGHTS_AGENT_TOKEN" \
  -H "Content-Type: application/json")

echo "$RESPONSE" | jq .

# Extract excluded tests
mapfile -t EXCLUDED_TESTS < <(echo "$RESPONSE" | jq -r '.data.excludedTests[].testName')

# Prepare Ginkgo command with excluded tests
GINKGO_CMD=("ginkgo" "--json-report=report.json")
for TEST in "${EXCLUDED_TESTS[@]}"; do
    GINKGO_CMD+=("--skip=$TEST")
done

# Run tests
cd tests/e2e && "${GINKGO_CMD[@]}"

# Process test report
PROCESSED_JSON=$(
  cat "report.json" | jq -c '.[] | .SpecReports[]' | while IFS= read -r line; do
    name=$(echo "$line" | jq -r '.LeafNodeText')
    start_raw=$(echo "$line" | jq -r '.StartTime')
    end_raw=$(echo "$line" | jq -r '.EndTime')
    status=$(echo "$line" | jq -r '.State')

    # Process start time with `date`
    start=$(date --date="$start_raw" +%s%3N)

    # Check if end_raw is empty or equals "0001-01-01T00:00:00Z", and use the current time
    if [ -z "$end_raw" ] || [ "$end_raw" == "0001-01-01T00:00:00Z" ]; then
      end=$(date +%s%3N)
    else
      end=$(date --date="$end_raw" +%s%3N)
    fi

    echo "{\"name\": \"$name\", \"start\": $start, \"end\": $end, \"status\": \"$status\"}"
  done | jq -s '.'
)

echo "$PROCESSED_JSON" | jq .

# Send test results to Sealights
curl -X POST "https://$DOMAIN/sl-api/v2/test-sessions/$TEST_SESSION_ID" \
  -H "Authorization: Bearer $SEALIGHTS_AGENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d "${PROCESSED_JSON}"

# Delete the test session
curl -X DELETE "https://$DOMAIN/sl-api/v1/test-sessions/$TEST_SESSION_ID" \
  -H "Authorization: Bearer $SEALIGHTS_AGENT_TOKEN" \
  -H "Content-Type: application/json"

echo "INFO: Script completed successfully."
