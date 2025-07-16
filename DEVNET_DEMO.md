# SunRe AVS DevNet Demo Guide

## Overview

This guide demonstrates how to run the SunRe Weather Insurance AVS on local devnet and submit tasks that simulate production insurance contract interactions.

## Architecture

In production, the flow works as follows:

```
1. Insurance Contract → Submits monitoring task → AVS
2. AVS → Fetches weather data → Multiple APIs
3. AVS → Reaches consensus → Returns result
4. AVS → If threshold exceeded → Triggers claim
5. Insurance Contract → Processes payout → Policyholder
```

## Quick Start

### 1. Start the DevNet

```bash
# Option A: Full DevNet setup (recommended)
./scripts/start_avs_devnet.sh

# Option B: Manual setup
devkit avs devnet start
make build
./bin/performer --port 8080
```

### 2. Submit a Task

Tasks are submitted as base64-encoded JSON payloads:

```bash
# Example: Check current weather in NYC
TASK='{"type":"weather_check","location":{"latitude":40.7128,"longitude":-74.0060,"city":"New York","country":"USA"},"threshold":35.0}'
PAYLOAD=$(echo -n "$TASK" | base64)

devkit avs call --payload $PAYLOAD
```

## Task Types

### 1. Weather Monitoring Task

**Use Case**: Insurance contracts periodically check if weather conditions exceed policy thresholds.

```json
{
  "type": "weather_check",
  "location": {
    "latitude": 40.7128,
    "longitude": -74.0060,
    "city": "New York",
    "country": "USA"
  },
  "threshold": 35.0,
  "policy_id": "NYC-CROP-2024-001"
}
```

**Response**:
```json
{
  "type": "weather_check_response",
  "temperature": 32.5,
  "meets_threshold": false,
  "confidence": 0.95,
  "data_points": [
    {"source": "TomorrowIO", "temperature": 32.4},
    {"source": "WeatherAPI", "temperature": 32.6},
    {"source": "OpenMeteo", "temperature": 32.5}
  ]
}
```

### 2. Insurance Claim Verification

**Use Case**: When monitoring detects threshold breach, automatically verify and process claims.

```json
{
  "type": "insurance_claim",
  "claim_request": {
    "policy_id": "NYC-CROP-2024-001",
    "policy": {
      "insurance_type": "crop",
      "coverage_amount": 1000000,
      "triggers": [{
        "peril": "heat_wave",
        "conditions": {
          "temperature_max": 35,
          "consecutive_days": 3
        },
        "payout_ratio": 0.5
      }]
    },
    "claim_date": "2024-07-15T00:00:00Z"
  }
}
```

**Response**:
```json
{
  "claim_id": "CLM-abc123",
  "claim_status": "approved",
  "payout_amount": 500000,
  "triggered_perils": ["heat_wave"],
  "verification_hash": "0x1234...",
  "consensus_data": {
    "temperature": 38.2,
    "confidence": 0.92
  }
}
```

### 3. Live Weather Demo

**Use Case**: Real-time weather data demonstration with consensus from multiple sources.

```json
{
  "type": "live_weather_demo",
  "location": {
    "latitude": 40.7128,
    "longitude": -74.0060,
    "city": "New York",
    "country": "USA"
  }
}
```

## Production Simulation

### Automated Monitoring Script

```bash
# Start the monitoring simulation
./scripts/insurance_monitor.sh

# This simulates:
# - Periodic weather checks (every 60 seconds)
# - Automatic claim submission when thresholds exceeded
# - Multiple policies across different locations
```

### Manual Task Submission

```bash
# Use the demo script to see all task types
./scripts/submit_task_demo.sh

# Submit specific tasks
./scripts/submit_task_demo.sh | grep "devkit avs call" | bash
```

## DevNet Details

### Running Services

```bash
# Check running services
devkit avs devnet list

# Expected output:
# CONTAINER ID   NAME                 PORTS
# abc123         anvil-l1             0.0.0.0:8545->8545/tcp
# def456         anvil-l2             0.0.0.0:8546->8546/tcp
# ghi789         avs-performer        0.0.0.0:8080->8080/tcp
```

### Logs

```bash
# View performer logs
devkit avs logs performer

# View task results
tail -f performer.log

# View monitoring logs
tail -f monitoring.log
```

## Integration Points

### For Smart Contract Developers

```solidity
// Example integration in your insurance contract
interface IAVSTaskSubmitter {
    function submitTask(bytes calldata taskData) external returns (bytes32 taskId);
}

contract WeatherInsurance {
    IAVSTaskSubmitter public avs;
    
    function checkWeatherConditions(bytes32 policyId) external {
        Policy memory policy = policies[policyId];
        
        bytes memory taskData = abi.encode(
            "weather_check",
            policy.location.lat,
            policy.location.lon,
            policy.temperatureThreshold
        );
        
        bytes32 taskId = avs.submitTask(taskData);
        pendingTasks[taskId] = policyId;
    }
}
```

### For Frontend Developers

```javascript
// Submit task via API
async function submitWeatherCheck(location, threshold) {
  const task = {
    type: "weather_check",
    location: location,
    threshold: threshold
  };
  
  const payload = btoa(JSON.stringify(task));
  
  const response = await fetch('http://localhost:8080/task', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ payload })
  });
  
  return response.json();
}
```

## Troubleshooting

### AVS Not Responding

```bash
# Check if performer is running
ps aux | grep performer

# Restart performer
pkill -f performer
./bin/performer --port 8080 &

# Check logs
tail -f performer.log
```

### Task Submission Fails

```bash
# Verify task format
echo $PAYLOAD | base64 -d | jq .

# Check devnet status
devkit avs devnet list

# Restart devnet
devkit avs devnet stop
devkit avs devnet start
```

### No Weather Data

```bash
# Check API keys in main.go
grep -E "APIKey|api_key" cmd/main.go

# Test API connectivity
curl -s "https://api.tomorrow.io/v4/weather/realtime?location=40.7128,-74.0060&apikey=YOUR_KEY"
```

## Next Steps

1. **Deploy to Testnet**
   - Update RPC URLs in config
   - Register operators
   - Deploy contracts

2. **Production Setup**
   - Implement key management
   - Set up monitoring
   - Configure alerts

3. **Integration**
   - Connect insurance contracts
   - Add frontend dashboard
   - Set up automated monitoring

## Resources

- [Production Deployment Guide](PRODUCTION_DEPLOYMENT.md)
- [API Documentation](docs/API.md)
- [Smart Contract Examples](contracts/examples/)
- [DevKit Documentation](https://docs.eigenlayer.xyz/devkit)