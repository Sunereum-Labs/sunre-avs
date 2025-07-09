# DevKit CLI Integration Guide

This guide ensures the Weather Insurance AVS works correctly with EigenLayer DevKit CLI commands.

## ✅ Verified DevKit Commands

### 1. Build the AVS
```bash
devkit avs build
```
This runs `make build` which compiles the performer binary to `./bin/performer`

### 2. Start Local DevNet
```bash
devkit avs devnet start
```
This starts:
- Local Anvil chains (L1 and L2)
- Deploys EigenLayer contracts
- Starts the AVS performer on port 8080

### 3. Submit Tasks

#### Weather Check Task
```bash
# Create payload
PAYLOAD=$(echo -n '{"type":"weather_check","location":{"latitude":40.7128,"longitude":-74.0060,"city":"New York","country":"USA"},"threshold":25.0}' | base64)

# Submit task
devkit avs call --payload "$PAYLOAD"
```

#### Insurance Claim Task
```bash
# Create payload (simplified for readability)
PAYLOAD=$(cat <<'EOF' | jq -c | base64
{
  "type": "insurance_claim",
  "claim_request": {
    "policy_id": "CROP-2024-001",
    "policy": {
      "policy_id": "CROP-2024-001",
      "policy_holder": "Green Valley Farms",
      "insurance_type": "crop",
      "location": {
        "latitude": 35.2271,
        "longitude": -80.8431,
        "city": "Charlotte",
        "country": "USA"
      },
      "coverage_amount": 100000,
      "premium": 5000,
      "start_date": "2024-06-01T00:00:00Z",
      "end_date": "2024-09-30T00:00:00Z",
      "triggers": [{
        "trigger_id": "HEAT-3DAY",
        "peril": "heat_wave",
        "conditions": {
          "temperature_max": 35,
          "consecutive_days": 3
        },
        "payout_ratio": 0.5,
        "description": "Heat protection"
      }]
    },
    "claim_date": "2024-07-15T00:00:00Z",
    "automated_check": true
  },
  "demo_mode": true,
  "demo_scenario": "heat_wave"
}
EOF
)

# Submit task
devkit avs call --payload "$PAYLOAD"
```

### 4. View Logs
```bash
# View performer logs
devkit avs logs performer

# View all logs
devkit avs logs
```

### 5. Stop DevNet
```bash
devkit avs devnet stop
```

## 📝 Payload Format

The AVS expects base64-encoded JSON payloads with a `type` field:

### Weather Check Format
```json
{
  "type": "weather_check",
  "location": {
    "latitude": number,
    "longitude": number,
    "city": string,
    "country": string
  },
  "threshold": number
}
```

### Insurance Claim Format
```json
{
  "type": "insurance_claim",
  "claim_request": {
    "policy_id": string,
    "policy": InsurancePolicy,
    "claim_date": ISO8601 string,
    "automated_check": boolean
  },
  "demo_mode": boolean (optional),
  "demo_scenario": string (optional: "heat_wave", "cold_snap", "normal")
}
```

## 🔧 Troubleshooting

### Issue: "Task validation failed"
**Solution**: Ensure your JSON is valid and base64 encoded correctly:
```bash
# Test JSON validity
echo "$JSON" | jq .

# Encode properly
echo -n "$JSON" | base64
```

### Issue: "No weather data sources"
**Solution**: The AVS uses Open-Meteo by default (no API key needed). For production, add API keys:
```bash
export OPENWEATHERMAP_API_KEY="your-key"
export WEATHERAPI_API_KEY="your-key"
devkit avs devnet start
```

### Issue: "Insufficient responses"
**Solution**: This is normal in demo mode with only Open-Meteo. The consensus requires 3+ sources. Use `demo_mode: true` for testing.

## 🧪 Quick Test Script

Run all tests:
```bash
./scripts/test_devkit_integration.sh
```

This will show you the exact commands to run for each scenario.

## 📊 Expected Responses

### Weather Check Response
```json
{
  "type": "weather_check_response",
  "temperature": 23.5,
  "meets_threshold": false,
  "confidence": 0.95,
  "data_points": 3,
  "timestamp": 1234567890
}
```

### Insurance Claim Response
```json
{
  "claim_id": "CLM-abc12345",
  "policy_id": "CROP-2024-001",
  "claim_status": "approved",
  "triggered_perils": [...],
  "payout_amount": 50000,
  "weather_data": {...},
  "verification_hash": "0x...",
  "timestamp": "2024-07-15T..."
}
```

## 🚀 Advanced Usage

### Custom Operators
```bash
# Deploy with custom operator count
devkit avs devnet start --operators 10
```

### Custom Port
```bash
# Run performer on different port
AVS_PORT=8090 devkit avs devnet start
```

### Production Deployment
```bash
# Deploy to testnet
devkit avs deploy --network holesky

# Deploy to mainnet
devkit avs deploy --network mainnet
```

## ✨ Demo Commands Summary

```bash
# 1. Build
devkit avs build

# 2. Start
devkit avs devnet start

# 3. Run demo
./scripts/insurance_demo.sh

# 4. Copy a payload and submit
devkit avs call --payload "<BASE64_PAYLOAD>"

# 5. Check logs
devkit avs logs performer
```

The AVS is fully compatible with DevKit CLI and ready for demonstration!