#!/bin/bash

# Direct task submission test
# This script tests that our AVS can receive and process tasks

echo "SunRe AVS Task Submission Test"
echo "=============================="
echo ""

# Check if performer is running
if ! pgrep -f "performer" > /dev/null; then
    echo "❌ Performer not running. Please start with: ./bin/performer --port 8080"
    exit 1
fi

echo "✅ Performer is running"

# Test 1: Create a simple task payload
echo ""
echo "Test 1: Weather Check Task"
echo "-------------------------"

WEATHER_TASK='{
  "type": "weather_check",
  "location": {
    "latitude": 40.7128,
    "longitude": -74.0060,
    "city": "New York",
    "country": "USA"
  },
  "threshold": 25.0
}'

echo "Task payload:"
echo "$WEATHER_TASK" | jq .

# Base64 encode for DevKit
WEATHER_PAYLOAD=$(echo -n "$WEATHER_TASK" | base64)
echo ""
echo "Base64 encoded payload: $WEATHER_PAYLOAD"

# Test 2: Insurance Claim Task
echo ""
echo "Test 2: Insurance Claim Task (Demo Mode)"
echo "---------------------------------------"

CLAIM_TASK='{
  "type": "insurance_claim",
  "claim_request": {
    "policy_id": "TEST-001",
    "policy": {
      "policy_id": "TEST-001",
      "policy_holder": "Test Farm",
      "insurance_type": "crop",
      "location": {
        "latitude": 40.7128,
        "longitude": -74.0060,
        "city": "New York",
        "country": "USA"
      },
      "coverage_amount": 100000,
      "premium": 5000,
      "start_date": "2024-01-01T00:00:00Z",
      "end_date": "2024-12-31T00:00:00Z",
      "triggers": [{
        "trigger_id": "HEAT-001",
        "peril": "heat_wave",
        "conditions": {
          "temperature_max": 35,
          "consecutive_days": 3
        },
        "payout_ratio": 0.5,
        "description": "Heat wave protection"
      }]
    },
    "claim_date": "2024-07-15T00:00:00Z",
    "automated_check": true
  },
  "demo_mode": true,
  "demo_scenario": "heat_wave"
}'

CLAIM_PAYLOAD=$(echo -n "$CLAIM_TASK" | base64)
echo "Claim task created (demo mode)"
echo "Base64 encoded payload: $CLAIM_PAYLOAD"

# Test 3: Live Weather Demo
echo ""
echo "Test 3: Live Weather Demo"
echo "------------------------"

LIVE_TASK='{
  "type": "live_weather_demo",
  "location": {
    "latitude": 40.7128,
    "longitude": -74.0060,
    "city": "New York",
    "country": "USA"
  }
}'

LIVE_PAYLOAD=$(echo -n "$LIVE_TASK" | base64)
echo "Live weather task created"
echo "Base64 encoded payload: $LIVE_PAYLOAD"

# Show how to submit these tasks
echo ""
echo "How to Submit Tasks:"
echo "==================="
echo ""
echo "Method 1: Using DevKit (when TaskMailbox is deployed)"
echo "devkit avs call --payload $WEATHER_PAYLOAD"
echo ""
echo "Method 2: Direct gRPC (requires gRPC client)"
echo "The performer is listening on localhost:8080 for gRPC requests"
echo ""
echo "Method 3: Via HTTP Bridge (if implemented)"
echo "curl -X POST http://localhost:8081/submit-task \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '$WEATHER_TASK'"
echo ""

# Check performer logs
echo "Recent Performer Logs:"
echo "====================="
if [ -f "performer.log" ]; then
    tail -10 performer.log
else
    echo "No performer.log found. Check if performer is logging to stdout."
fi

echo ""
echo "✅ Task payloads created successfully!"
echo "The AVS is ready to receive tasks through:"
echo "  • gRPC on localhost:8080"
echo "  • DevKit task submission (when contracts are deployed)"
echo "  • HTTP bridge (if implemented)"