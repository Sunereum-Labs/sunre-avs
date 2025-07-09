# Temperature Oracle AVS

This AVS implements a decentralized temperature verification oracle that checks if a location's temperature meets a specified threshold using consensus from multiple weather data sources.

## Architecture

The Temperature Oracle AVS follows the Hourglass framework architecture:

### Components

1. **Performer (cmd/main.go)**: The main AVS service that handles task requests
   - Validates incoming temperature verification tasks
   - Orchestrates the temperature data collection
   - Returns consensus results

2. **Internal Modules**:
   - **aggregator**: Manages task distribution and response collection
   - **consensus**: Implements MAD (Median Absolute Deviation) algorithm
   - **datasources**: Weather API integrations with rate limiting
   - **executor**: Parallel data fetching from weather sources
   - **types**: Shared data structures

### Task Flow

1. **Task Request**: DevKit sends a task with location and threshold
2. **Validation**: Performer validates coordinates and threshold
3. **Distribution**: Task distributed to simulated operators
4. **Execution**: Operators fetch data from weather APIs
5. **Consensus**: MAD algorithm filters outliers and reaches consensus
6. **Response**: Returns temperature, threshold result, and confidence

## Usage with DevKit

### 1. Build the AVS
```bash
devkit avs build
```

### 2. Start the DevNet
```bash
devkit avs devnet start
```

### 3. Create a Temperature Verification Task

Example task payload:
```json
{
  "location": {
    "latitude": 40.7128,
    "longitude": -74.0060,
    "city": "New York",
    "country": "USA"
  },
  "threshold": 25.0
}
```

Submit the task:
```bash
# Encode payload as base64
PAYLOAD=$(echo -n '{"location":{"latitude":40.7128,"longitude":-74.0060,"city":"New York","country":"USA"},"threshold":25.0}' | base64)

# Call the AVS
devkit avs call --payload "$PAYLOAD"
```

### 4. Response Format

The AVS returns:
```json
{
  "temperature": 23.45,
  "meets_threshold": false,
  "confidence": 0.925,
  "data_points": 4,
  "timestamp": 1736449200
}
```

## Weather Data Sources

Currently configured to use Open-Meteo (no API key required). Can be extended with:
- OpenWeatherMap
- WeatherAPI.com
- Tomorrow.io
- Visual Crossing

## Configuration

The AVS uses hardcoded configuration for demo purposes:
- Minimum operators: 3
- Response timeout: 60s
- Consensus threshold: 67%
- MAD threshold: 2.5
- Cache TTL: 5 minutes

## Testing

Run the test script:
```bash
./scripts/test_temperature_task.sh
```

This will show example payloads for different cities and how to submit tasks.

## Production Considerations

For production deployment:
1. Add real API keys for multiple weather sources
2. Implement proper operator management
3. Add persistent storage for task history
4. Configure monitoring and alerts
5. Implement proper fee mechanisms