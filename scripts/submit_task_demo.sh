#!/bin/bash

# Simple task submission demo
# This shows how to submit tasks to the AVS

echo "SunRe AVS Task Submission Demo"
echo "=============================="
echo ""

# Function to encode JSON to base64
encode_task() {
    echo -n "$1" | base64
}

# 1. Weather Check Task (Regular monitoring)
echo "1. Weather Monitoring Task"
echo "   Use case: Insurance contract checking if conditions exceed threshold"
echo ""

WEATHER_TASK_JSON='{
  "type": "weather_check",
  "location": {
    "latitude": 40.7128,
    "longitude": -74.0060,
    "city": "New York",
    "country": "USA"
  },
  "threshold": 35.0,
  "policy_details": {
    "policy_id": "NYC-CROP-2024-001",
    "insured_amount": 1000000,
    "trigger_temperature": 35
  }
}'

WEATHER_TASK_PAYLOAD=$(encode_task "$WEATHER_TASK_JSON")

echo "Encoded payload:"
echo "$WEATHER_TASK_PAYLOAD"
echo ""
echo "Submit with:"
echo "devkit avs call --payload $WEATHER_TASK_PAYLOAD"
echo ""

# 2. Insurance Claim Task (Triggered by conditions)
echo "2. Insurance Claim Verification Task"
echo "   Use case: Automated claim when weather conditions are met"
echo ""

CLAIM_TASK_JSON='{
  "type": "insurance_claim",
  "claim_request": {
    "policy_id": "NYC-CROP-2024-001",
    "policy": {
      "policy_id": "NYC-CROP-2024-001",
      "policy_holder": "Manhattan Urban Farm",
      "insurance_type": "crop",
      "location": {
        "latitude": 40.7128,
        "longitude": -74.0060,
        "city": "New York",
        "country": "USA"
      },
      "coverage_amount": 1000000,
      "premium": 50000,
      "start_date": "2024-01-01T00:00:00Z",
      "end_date": "2024-12-31T00:00:00Z",
      "triggers": [{
        "trigger_id": "HEAT-NYC-001",
        "peril": "heat_wave",
        "conditions": {
          "temperature_max": 35,
          "consecutive_days": 3
        },
        "payout_ratio": 0.5,
        "description": "Heat wave crop protection"
      }]
    },
    "claim_date": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
    "automated_check": true
  }
}'

CLAIM_TASK_PAYLOAD=$(encode_task "$CLAIM_TASK_JSON")

echo "Encoded payload:"
echo "$CLAIM_TASK_PAYLOAD"
echo ""
echo "Submit with:"
echo "devkit avs call --payload $CLAIM_TASK_PAYLOAD"
echo ""

# 3. Live Demo Task
echo "3. Live Weather Demo Task"
echo "   Use case: Real-time weather data with consensus from 3 sources"
echo ""

DEMO_TASK_JSON='{
  "type": "live_weather_demo",
  "location": {
    "latitude": 40.7128,
    "longitude": -74.0060,
    "city": "New York",
    "country": "USA"
  }
}'

DEMO_TASK_PAYLOAD=$(encode_task "$DEMO_TASK_JSON")

echo "Encoded payload:"
echo "$DEMO_TASK_PAYLOAD"
echo ""
echo "Submit with:"
echo "devkit avs call --payload $DEMO_TASK_PAYLOAD"
echo ""

# Test submission
echo "Testing task submission..."
echo "=========================="

# Check if AVS is running
if curl -s http://localhost:8080 >/dev/null 2>&1; then
    echo "✓ AVS performer is running"
    
    # Submit demo task
    echo ""
    echo "Submitting live weather demo task..."
    if devkit avs call --payload "$DEMO_TASK_PAYLOAD" 2>&1; then
        echo "✓ Task submitted successfully!"
    else
        echo "✗ Task submission failed"
        echo "Make sure devnet is running: ./scripts/start_avs_devnet.sh"
    fi
else
    echo "✗ AVS performer not running"
    echo ""
    echo "To start the AVS:"
    echo "1. Start devnet: devkit avs devnet start"
    echo "2. Start performer: ./bin/performer --port 8080"
    echo "Or use: ./scripts/start_avs_devnet.sh"
fi