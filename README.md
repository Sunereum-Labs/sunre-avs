# SunRe AVS - Decentralized Weather Insurance Platform

<p align="center">
  <img src="https://img.shields.io/badge/EigenLayer-AVS-blue" alt="EigenLayer AVS">
  <img src="https://img.shields.io/badge/Consensus-BLS-green" alt="BLS Consensus">
  <img src="https://img.shields.io/badge/Weather-Insurance-orange" alt="Weather Insurance">
</p>

## Overview

SunRe AVS is a decentralized weather insurance platform built on EigenLayer's AVS (Actively Validated Service) architecture. It enables parametric insurance products that automatically process claims based on weather data consensus from multiple oracles.

### Key Features

- 🌡️ **Multi-Source Weather Consensus** - Aggregates data from 3+ weather APIs using MAD algorithm
- 🔐 **BLS Signature Aggregation** - Cryptographic consensus among EigenLayer operators
- 📊 **Parametric Insurance** - Automatic claim processing when weather triggers are met
- ⚡ **Task-Based Architecture** - Insurance contracts submit monitoring tasks to AVS
- 🚀 **Production Ready** - Full DevKit integration for testnet/mainnet deployment

## 🚀 Quick Start (60 seconds)

```bash
# Run the complete end-to-end demo
./run_demo.sh

# Verify AVS is processing tasks
./prove_avs.sh

# Access the demo at: http://localhost:3000
```

### Alternative Demo Options
```bash
# Interactive launcher with menu
./demo.sh

# Just the UI (mock mode)
cd demo-ui && npm start
```

## Architecture

```
┌─────────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│ Insurance Contract  │────▶│   SunRe AVS      │────▶│ Weather APIs    │
│ (Submits Tasks)     │     │ (Consensus)      │     │ (3+ sources)    │
└─────────────────────┘     └──────────────────┘     └─────────────────┘
         │                           │
         ▼                           ▼
┌─────────────────────┐     ┌──────────────────┐
│  Claim Processor    │◀────│ EigenLayer Core  │
│ (Auto Payouts)      │     │ (BLS Signatures) │
└─────────────────────┘     └──────────────────┘
```

## How It Works

1. **Task Submission**: Insurance contracts submit monitoring tasks to the AVS
2. **Data Collection**: AVS operators fetch weather data from multiple sources
3. **Consensus**: MAD algorithm filters outliers and reaches agreement
4. **Verification**: BLS signatures provide cryptographic proof
5. **Claim Processing**: Smart contracts automatically process payouts

## Supported Insurance Types

| Type | Use Case | Example Triggers |
|------|----------|-----------------|
| 🌾 **Crop** | Agricultural protection | Heat waves, frost, drought |
| 🎪 **Event** | Outdoor event cancellation | Rain, wind, extreme weather |
| ✈️ **Travel** | Flight delay compensation | Extreme temperatures, storms |
| 🏢 **Property** | Weather damage coverage | Hail, flooding, hurricanes |

## Task Types

### Weather Monitoring Task
```json
{
  "type": "weather_check",
  "location": {"latitude": 40.7128, "longitude": -74.0060},
  "threshold": 35.0,
  "policy_id": "POLICY-001"
}
```

### Insurance Claim Task
```json
{
  "type": "insurance_claim",
  "claim_request": {
    "policy_id": "POLICY-001",
    "policy": {
      "insurance_type": "crop",
      "triggers": [{
        "peril": "heat_wave",
        "conditions": {"temperature_max": 35},
        "payout_ratio": 0.5
      }]
    }
  }
}
```

## Development

### Prerequisites

- Go 1.21+
- Node.js 16+
- Docker
- [EigenLayer DevKit](https://github.com/Layr-Labs/eigenlayer-devkit)

### Build & Test

```bash
# Build AVS performer
make build

# Run unit tests
make test

# Start local DevNet
devkit avs devnet start

# Submit a task
devkit avs call --payload <base64-encoded-task>
```

### Project Structure

```
├── cmd/                    # Main AVS performer
├── internal/               # Core business logic
│   ├── consensus/         # MAD consensus algorithm
│   ├── weather/           # Weather data sources
│   └── insurance/         # Claim processing
├── contracts/             # Smart contracts
├── demo-ui/               # React demo interface
├── scripts/               # Deployment scripts
└── config/                # Network configurations
```

## Weather Data Sources

The AVS uses multiple weather APIs with built-in rate limiting:

- **Tomorrow.io** - High precision weather data (API key included for demo)
- **WeatherAPI.com** - Global coverage (API key included for demo)
- **Open-Meteo** - Open source fallback (no key required)

## Production Deployment

See [PRODUCTION_DEPLOYMENT.md](PRODUCTION_DEPLOYMENT.md) for:
- Holesky testnet deployment
- Mainnet security checklist
- Monitoring and alerts
- Disaster recovery

## Documentation

- 📘 [DevNet Demo Guide](DEVNET_DEMO.md) - Local development with task examples
- 🚀 [Production Deployment](PRODUCTION_DEPLOYMENT.md) - Mainnet deployment guide
- 🧪 [Testing Guide](TESTING.md) - Test coverage and strategies

## Demo UI

The demo UI provides an interactive interface to:
- View system architecture and consensus process
- Test different insurance scenarios
- Monitor live weather data from NYC
- Submit and track insurance claims

Access at http://localhost:3000 after running `./demo.sh`

## Smart Contract Integration

```solidity
// Example: Insurance contract submitting monitoring task
interface IAVSTaskSubmitter {
    function submitTask(bytes calldata taskData) external returns (bytes32);
}

contract WeatherInsurance {
    function monitorWeather(bytes32 policyId) external {
        bytes memory task = abi.encode(
            "weather_check",
            location,
            threshold,
            policyId
        );
        avs.submitTask(task);
    }
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing`)
3. Commit your changes (`git commit -m 'Add feature'`)
4. Push to the branch (`git push origin feature/amazing`)
5. Open a Pull Request

## Security

- Smart contracts audited by [Pending]
- Bug bounty: security@sunre-avs.com
- Consensus algorithm prevents manipulation
- All decisions cryptographically verifiable

## License

MIT License - see [LICENSE](LICENSE) file

## Acknowledgments

- Built on [EigenLayer](https://eigenlayer.xyz/) infrastructure
- Uses [Hourglass](https://github.com/Layr-Labs/hourglass) framework
- Weather data from Tomorrow.io, WeatherAPI.com, and Open-Meteo

---

<p align="center">
  Built with ❤️ for transparent, automated insurance
</p>