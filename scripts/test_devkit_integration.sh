#!/bin/bash

echo "Testing DevKit Integration"
echo "=========================="
echo ""

# Test 1: Simple Weather Check (Original functionality)
echo "1. Testing Simple Weather Check:"
WEATHER_PAYLOAD=$(echo -n '{"type":"weather_check","location":{"latitude":40.7128,"longitude":-74.0060,"city":"New York","country":"USA"},"threshold":25.0}' | base64)
echo "devkit avs call --payload $WEATHER_PAYLOAD"
echo ""

# Test 2: Insurance Claim with Demo Mode
echo "2. Testing Insurance Claim (Demo Mode):"
INSURANCE_PAYLOAD=$(echo -n '{"type":"insurance_claim","claim_request":{"policy_id":"TEST-001","policy":{"policy_id":"TEST-001","policy_holder":"Test Farm","insurance_type":"crop","location":{"latitude":35.2271,"longitude":-80.8431,"city":"Charlotte","country":"USA"},"coverage_amount":100000,"premium":5000,"start_date":"2024-06-01T00:00:00Z","end_date":"2024-09-30T00:00:00Z","triggers":[{"trigger_id":"HEAT-TEST","peril":"heat_wave","conditions":{"temperature_max":35,"consecutive_days":3},"payout_ratio":0.5,"description":"Heat protection"}]},"claim_date":"2024-07-15T00:00:00Z","automated_check":true},"demo_mode":true,"demo_scenario":"heat_wave"}' | base64)
echo "devkit avs call --payload $INSURANCE_PAYLOAD"
echo ""

echo "Expected Response Format:"
echo "========================"
echo ""
echo "For Weather Check:"
echo '{"type":"weather_check_response","temperature":23.5,"meets_threshold":false,"confidence":0.95,"data_points":3,"timestamp":1234567890}'
echo ""
echo "For Insurance Claim:"
echo '{"claim_id":"CLM-abc123","policy_id":"TEST-001","claim_status":"approved","triggered_perils":[...],"payout_amount":50000,"weather_data":{...},"verification_hash":"0x...","timestamp":"2024-..."}'