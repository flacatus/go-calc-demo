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

# Run tests
cd tests/e2e && ginkgo --json-report="report.json"
