export DOMAIN="redhat.sealights.co"
export SEALIGHTS_AGENT_TOKEN=""
export CONTAINER_ROUTE_URL="http://127.0.0.1:8080"

# Make POST request and extract testSessionId
# labId can be provided. 
TEST_SESSION_ID=$(curl -X POST "https://$DOMAIN/sl-api/v1/test-sessions" \
  -H "Authorization: Bearer $SEALIGHTS_AGENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"labId":"","testStage":"go-calc-integration","bsid":"37fd56fc-13be-4a72-a139-ae28b9df5220","sessionTimeout":10000}' | jq -r '.data.testSessionId')

# Check if TEST_SESSION_ID was retrieved successfully
if [ -n "$TEST_SESSION_ID" ]; then
  echo "Test session ID: $TEST_SESSION_ID"
  export TEST_SESSION_ID
else
  echo "Failed to retrieve test session ID"
  exit 1
fi

RESPONSE=$(curl -X GET "https://$DOMAIN/sl-api/v2/test-sessions/$TEST_SESSION_ID/exclude-tests" \
  -H "Authorization: Bearer $SEALIGHTS_AGENT_TOKEN" \
  -H "Content-Type: application/json"
)

echo $RESPONSE | jq .

# Extract all the test names into an array
mapfile -t EXCLUDED_TESTS < <(echo "$RESPONSE" | jq -r '.data.excludedTests[].testName')

GINKGO_CMD=("ginkgo" "--json-report=report.json")

for TEST in "${EXCLUDED_TESTS[@]}"; do
    GINKGO_CMD+=("--skip=$TEST")
done

cd tests/e2e && "${GINKGO_CMD[@]}"

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

echo $PROCESSED_JSON | jq .

curl -X POST "https://$DOMAIN/sl-api/v2/test-sessions/$TEST_SESSION_ID" \
  -H "Authorization: Bearer $SEALIGHTS_AGENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d "${PROCESSED_JSON}"

curl -X DELETE "https://$DOMAIN/sl-api/v1/test-sessions/$TEST_SESSION_ID" \
  -H "Authorization: Bearer $SEALIGHTS_AGENT_TOKEN" \
  -H "Content-Type: application/json"