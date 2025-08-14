# SunRe AVS - Parametric Weather Insurance Service powered by EigenLayer

An Actively Validated Service (AVS) for parametric weather insurance built with EigenLayer DevKit. 

## What This Is

SunRe AVS is a weather verification service that enables truly verifiable parametric insurance. The AVS:
- **Verifies weather data** from multiple sources for insurance claims
- **Achieves consensus** among operators using BLS signatures
- **Processes claims automatically** based on weather triggers

## Example Use Case: Crop Insurance

Imagine a farmer in Iowa purchases parametric crop insurance for drought protection:

1. **Policy Creation**: Farmer buys insurance that pays out if rainfall < 200mm during growing season
2. **Weather Monitoring**: SunRe AVS continuously monitors weather conditions at the farm's location
3. **Trigger Event**: A drought occurs with only 150mm rainfall recorded
4. **Verification Process**:
   - Multiple operators fetch weather data from different sources
   - Each operator signs the weather data with their BLS signature
   - Consensus is reached when 67% of stake agrees on the conditions
5. **Automatic Payout**: Smart contract automatically releases insurance payout to farmer

No claims adjusters, no paperwork, no delays - just transparent, verifiable weather data triggering automatic payments.

### Other Real-World Applications

- **Flight Delay Insurance**: Automatic payouts when flights delayed due to weather
- **Event Cancellation**: Venues get compensated when outdoor events are rained out  
- **Energy Trading**: Solar/wind farms hedge against low production days
- **Tourism Protection**: Hotels compensated for excessive rain during peak season
- **Construction Delays**: Contractors protected against weather-related delays

## Quick Start

### Prerequisites
```bash
# Install DevKit CLI
curl -sSL https://install.eigenlayer.xyz | sh

# Install Foundry for smart contracts
curl -L https://foundry.paradigm.xyz | bash
foundryup

# Install Go 1.23+
# Visit https://go.dev/dl/
```

### One-Command Setup
```bash
# Clone and setup
git clone https://github.com/Sunereum-Labs/sunre-avs.git
cd sunre-avs

# Install and build everything
./scripts/setup.sh

# Start local development
./scripts/start.sh
```

### Manual Setup

#### 1. Install Dependencies
```bash
# Install Go dependencies
go mod tidy

# Install contract dependencies
cd contracts && forge install
```

#### 2. Build Everything
```bash
# Build contracts
make build

# Build Go binary
go build -o bin/sunre-avs cmd/main.go
```

#### 3. Run Local Development
```bash
# Start local devnet
make devnet

# In another terminal, deploy contracts
make deploy

# Start the AVS performer
./bin/sunre-avs
```

#### 4. Submit a Test Task
```bash
# Submit weather verification task
devkit avs call --input examples/task-weather-nyc.json
```

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Client    │────▶│ TaskAVSReg   │────▶│  Operators  │
│  Insurance  │     │     (L1)     │     │  (3+ nodes) │
└─────────────┘     └──────────────┘     └─────────────┘
                            │                     │
                            ▼                     ▼
                    ┌──────────────┐     ┌─────────────┐
                    │ AVSTaskHook  │────▶│   Weather   │
                    │     (L2)     │     │   Sources   │
                    └──────────────┘     └─────────────┘
                            │
                            ▼
                    ┌──────────────┐
                    │ BN254Verifier│
                    │  (Consensus) │
                    └──────────────┘
```

### Core Components

1. **L1 Contract - TaskAVSRegistrar** (`contracts/src/l1-contracts/`)
   - Manages operator registration and stake
   - Inherits from DevKit's `TaskAVSRegistrarBase`
   - Handles task submission and validation

2. **L2 Contracts** (`contracts/src/l2-contracts/`)
   - **AVSTaskHook**: Validates tasks and manages payments
   - **BN254CertificateVerifier**: Verifies BLS signatures for consensus

3. **Performer** (`cmd/main.go`)
   - Processes weather verification tasks
   - Fetches data from Open-Meteo API (no key required)
   - Implements caching and fallback mechanisms
   - Provides health and metrics endpoints

## 🌡️ Weather Verification Flow

1. **Task Submission**: Insurance contract submits weather verification request
2. **Validation**: AVSTaskHook validates the request and payment
3. **Data Fetching**: Operators fetch weather data from multiple sources
4. **Consensus**: Operators reach consensus on weather conditions
5. **Verification**: BLS signatures are aggregated and verified
6. **Result**: Verified weather data returned for insurance payout

## 🔧 Configuration

### Environment Variables (.env)
```bash
# Network
ENV=development
PERFORMER_PORT=8080
PERFORMER_TIMEOUT=5s
HEALTH_PORT=8081

# Operator (for production)
OPERATOR_ID=your-operator-id
OPERATOR_KEY=your-operator-key

# Optional: Weather API keys for additional sources
OPENMETEO_API_KEY=optional
```

### DevKit Configuration (`config/devkit.yaml`)
```yaml
project:
  name: "sunre"
  version: "1.0.0"
  
avs:
  min_operators: 3
  consensus_threshold: 0.67
  
networks:
  devnet:
    rpc_url: "http://localhost:8545"
  testnet:
    rpc_url: "https://ethereum-holesky.publicnode.com"
```

## 🧪 Testing

### Run Tests
```bash
# Go tests
go test ./...

# Contract tests
cd contracts && forge test

# All tests
./scripts/test.sh all
```

### Test Task Payloads

Example task payload (`examples/task-weather-nyc.json`):
```json
{
  "location": {
    "latitude": 40.7128,
    "longitude": -74.0060,
    "city": "New York"
  },
  "timestamp": 1704067200,
  "policy_id": "POL-NYC-2024-001"
}
```

Submit test tasks:
```bash
# New York weather (real-time data)
devkit avs call --input examples/task-weather-nyc.json

# Miami weather (hurricane-prone region)
devkit avs call --input examples/task-weather-miami.json

# London weather (frequent rain events)
devkit avs call --input examples/task-weather-london.json
```

#### Testing with Real Weather Events:

To test the system with actual weather conditions:

1. **Current Weather**: Set timestamp to current time for real-time verification
2. **Historical Events**: Use past timestamps to verify known weather events
3. **Extreme Conditions**: Test during storms, heatwaves, or other notable weather

Example - Testing a known rainfall event:
```bash
# Create a task for a specific date/location
echo '{
  "location": {"latitude": 25.7617, "longitude": -80.1918, "city": "Miami"},
  "timestamp": 1693526400,  # Hurricane Idalia - Sept 2023
  "policy_id": "TEST-HURRICANE-001"
}' > test-hurricane.json

devkit avs call --input test-hurricane.json
```

The system will fetch historical weather data and verify the extreme conditions that occurred during Hurricane Idalia.

## Deployment

### Weather API Configuration for Production

For production deployments, operators should configure multiple weather data sources for redundancy and accuracy. While the AVS works with the free Open-Meteo API by default, production deployments should add premium sources:

#### Required Configuration in `.env`:

```bash
# Primary source (free, no key required)
# Open-Meteo API is used by default

# Optional premium sources for production reliability
WEATHERAPI_KEY=your_weatherapi_com_key        # weatherapi.com
OPENWEATHER_API_KEY=your_openweather_key       # openweathermap.org
TOMORROW_IO_KEY=your_tomorrow_io_key           # tomorrow.io
WEATHER_GOV_KEY=your_weather_gov_key           # weather.gov (US only)
```

#### Adding Weather Sources in `cmd/main.go`:

```go
// Add new weather source
if apiKey := os.Getenv("WEATHERAPI_KEY"); apiKey != "" {
    weatherClient.AddSource("weatherapi", apiKey, 2) // priority 2
}
```

**Note**: Each operator can use different weather sources. The consensus mechanism ensures accuracy even if operators use different APIs.

#### How Consensus Works with Multiple Sources:

1. **Data Collection**: Each operator fetches data from their configured sources
2. **Aggregation**: Operator calculates median values from their sources  
3. **Submission**: Operator signs and submits their weather data
4. **Consensus**: System accepts data when 67% of stake agrees (within 10% deviation)
5. **Slashing**: Operators submitting outlier data (>10% deviation) may be slashed

This design ensures no single weather API can manipulate the system, providing truly decentralized weather verification.

### Deploy to Testnet (Holesky)

1. **Setup Environment**
```bash
# Create env file with your keys
cp .env.example .env
# Edit .env with your PRIVATE_KEY_DEPLOYER and ETHERSCAN_API_KEY
# Add weather API keys for production sources (optional for testing)
```

2. **Get Testnet ETH**
```bash
# Get Holesky ETH from faucet
# https://holesky-faucet.pk910.de/
```

3. **Deploy Contracts**
```bash
./scripts/deploy.sh testnet
```

4. **Register Operators**
```bash
# Register as operator
devkit avs operator register --context testnet
```

5. **Start Services**
```bash
# Start performer with testnet config
ENV=production PERFORMER_PORT=8080 ./bin/sunre-avs
```

### Docker Deployment

```bash
# Build and run with Docker Compose
docker-compose up

# Or build manually
docker build -t sunre-avs .
docker run -p 8080:8080 -p 8081:8081 sunre-avs
```

## 📊 Monitoring

### Health Endpoints
- **Health Check**: `http://localhost:8081/health`
- **Metrics**: `http://localhost:8081/metrics`

### Metrics Tracked
- Tasks processed/succeeded/failed
- Average latency
- Weather API response times
- Cache hit rates

### Example Health Response
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "1.0.0"
}
```

## 🛠️ DevKit Commands

```bash
# Build AVS
devkit avs build

# Deploy contracts
devkit avs deploy

# Submit task
devkit avs call --input task.json

# Check status
devkit avs status

# View logs
devkit avs logs

# Switch context
devkit avs context set testnet
```

## 📁 Project Structure

```
sunre-avs/
├── cmd/
│   ├── main.go              # Main performer implementation
│   └── main_test.go         # Tests
├── contracts/
│   ├── src/
│   │   ├── l1-contracts/    # L1 contracts
│   │   │   └── TaskAVSRegistrar.sol
│   │   └── l2-contracts/    # L2 contracts
│   │       ├── AVSTaskHook.sol
│   │       └── BN254CertificateVerifier.sol
│   └── script/              # Deployment scripts
├── config/
│   ├── config.yaml          # Main configuration
│   ├── devkit.yaml          # DevKit configuration
│   └── contexts/
│       └── testnet.yaml     # Testnet configuration
├── examples/                # Task examples
│   ├── task-weather-nyc.json
│   ├── task-weather-miami.json
│   └── task-weather-london.json
├── scripts/
│   ├── setup.sh             # Install and build
│   ├── start.sh             # Start services
│   ├── test.sh              # Run tests
│   └── deploy.sh            # Deploy contracts
├── docker-compose.yml       # Docker orchestration
├── Dockerfile              # Container image
├── Makefile               # Build commands
└── README.md              # This file
```

## Troubleshooting

### Common Issues

**Issue**: Contract deployment fails
```bash
# Solution: Ensure you have enough testnet ETH
# Check balance:
cast balance YOUR_ADDRESS --rpc-url https://ethereum-holesky.publicnode.com
```

**Issue**: Weather API timeout
```bash
# Solution: Check network connectivity
# Test API directly:
curl "https://api.open-meteo.com/v1/forecast?latitude=40.7&longitude=-74.0&current=temperature_2m"
```

**Issue**: Operator registration fails
```bash
# Solution: Ensure minimum stake requirement
# Check requirement:
devkit avs info --context testnet
```

## 📚 Resources

- [EigenLayer Documentation](https://docs.eigenlayer.xyz)
- [DevKit Documentation](https://docs.eigenlayer.xyz/devkit)
- [Hourglass Architecture](https://docs.eigenlayer.xyz/devkit/architecture)
- [Open-Meteo API](https://open-meteo.com/en/docs)

## 🤝 Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Submit a pull request

## License

BUSL-1.1 - Business Source License 1.1

## Acknowledgments

Built with [EigenLayer DevKit]([https://eigenlayer.xyz](https://github.com/Layr-Labs/devkit-cli)) - making AVS development accessible to everyone.

---
