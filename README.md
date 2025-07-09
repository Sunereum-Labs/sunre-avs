# SunRe AVS

An EigenLayer Actively Validated Service (AVS) that enables automated, trust-minimized weather insurance claims processing using decentralized oracle consensus.

## Overview

This AVS demonstrates how parametric insurance can be revolutionized using blockchain technology. By leveraging EigenLayer's security model and decentralized weather data consensus, insurance claims are processed automatically when predefined weather conditions are met

## Key Features

- **Automated Claims Processing**: Smart contract-triggered claims based on weather parameters
- **Multi-Source Consensus**: MAD (Median Absolute Deviation) algorithm ensures data accuracy
- **Instant Settlement**: Claims processed in minutes instead of weeks
- **Cryptographic Verification**: Every decision is provable and auditable
- **Zero Fraud Risk**: Consensus mechanism eliminates false weather claims
- **Flexible Insurance Types**: Supports crop, event, travel, and custom parametric policies

## Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│ Insurance Smart │     │ EigenLayer AVS   │     │ Weather APIs    │
│ Contract        │────▶│ (Hourglass)      │────▶│ (5+ sources)    │
└─────────────────┘     └──────────────────┘     └─────────────────┘
                               │
                               ▼
                        ┌──────────────────┐
                        │ Consensus Engine │
                        │ (MAD Algorithm)  │
                        └──────────────────┘
                               │
                               ▼
                        ┌──────────────────┐
                        │ Automated Payout │
                        └──────────────────┘
```

## Quick Start

### Prerequisites

- Go 1.23.6 or higher
- [EigenLayer DevKit](https://github.com/Layr-Labs/eigenlayer-devkit) installed
- Docker (for containerized deployment)

### Installation

```bash
# Clone the repository
git clone https://github.com/Layr-Labs/hourglass-avs-template.git
cd hourglass-avs-template

# Build the AVS
devkit avs build

# Start local development environment
devkit avs devnet start
```

### Running the Demo

```bash
# Interactive demo with multiple insurance scenarios
./run_insurance_demo.sh

# Or submit individual tasks
devkit avs call --payload <BASE64_PAYLOAD>
```

## Usage Examples

### 1. Crop Insurance Claim

Protect farmers from heat damage with automatic payouts:

```json
{
  "type": "insurance_claim",
  "claim_request": {
    "policy": {
      "insurance_type": "crop",
      "location": {"latitude": 35.2271, "longitude": -80.8431},
      "triggers": [{
        "peril": "heat_wave",
        "conditions": {"temperature_max": 35, "consecutive_days": 3},
        "payout_ratio": 0.5
      }]
    }
  }
}
```

### 2. Event Cancellation Insurance

Automatic refunds for weather-cancelled events:

```json
{
  "type": "insurance_claim",
  "claim_request": {
    "policy": {
      "insurance_type": "event",
      "triggers": [{
        "peril": "excess_rain",
        "conditions": {"precipitation_min": 50},
        "payout_ratio": 1.0
      }]
    }
  }
}
```

## Supported Insurance Types

| Type | Use Case | Example Triggers |
|------|----------|-----------------|
| **Crop** | Agricultural protection | Heat waves, frost, drought |
| **Event** | Outdoor event cancellation | Rain, wind, extreme weather |
| **Travel** | Flight delay compensation | Extreme temperatures, storms |
| **Property** | Weather damage coverage | Hail, flooding, hurricanes |
| **Energy** | Renewable energy protection | Low wind/solar periods |

## Technical Specifications

### Consensus Mechanism

- **Algorithm**: Median Absolute Deviation (MAD)
- **Minimum Sources**: 3 weather APIs required
- **Outlier Threshold**: 2.5 standard deviations
- **Confidence Scoring**: Weighted by source reliability

### Performance

- **Claim Processing Time**: 2-3 minutes
- **Throughput**: 100+ claims per minute
- **Uptime**: 99.9% availability target

### Security

- **EigenLayer Security**: Leverages restaked ETH
- **BLS Signatures**: Aggregated operator signatures
- **Slashing Conditions**: Malicious behavior penalized
- **Verification**: All decisions cryptographically provable

## Development

### Project Structure

```
├── cmd/                    # Main AVS performer implementation
├── contracts/              # Smart contracts (L1 & L2)
├── internal/
│   ├── aggregator/        # Task distribution and collection
│   ├── consensus/         # MAD consensus algorithm
│   ├── datasources/       # Weather API integrations
│   ├── executor/          # Parallel task execution
│   ├── insurance/         # Claims processing logic
│   └── types/            # Shared data structures
├── scripts/              # Demo and deployment scripts
└── docs/                 # Additional documentation
```

### Adding Weather Sources

To add a new weather API:

1. Implement the `WeatherDataSource` interface in `internal/datasources/`
2. Add configuration in `config/config.yaml`
3. Register in `DataSourceManager`

### Creating Custom Insurance Products

1. Define triggers in `types.InsuranceTrigger`
2. Implement evaluation logic in `ClaimsProcessor`
3. Add demo scenarios for testing

## API Reference

### Task Types

#### Weather Check
```typescript
{
  type: "weather_check",
  location: Location,
  threshold: number
}
```

#### Insurance Claim
```typescript
{
  type: "insurance_claim",
  claim_request: InsuranceClaimRequest,
  demo_mode?: boolean,
  demo_scenario?: string
}
```

### Response Format

```typescript
{
  claim_id: string,
  claim_status: "approved" | "rejected" | "partial" | "investigate",
  payout_amount: number,
  triggered_perils: TriggeredPeril[],
  weather_data: WeatherAssessment,
  verification_hash: string
}
```

## Testing

```bash
# Run unit tests
make test

# Run integration tests with local devnet
devkit avs devnet start
./scripts/test_devkit_integration.sh
```

## Deployment

### Local Development
```bash
devkit avs devnet start
```

### Testnet (Holesky)
```bash
devkit avs deploy --network holesky
```

### Production
```bash
devkit avs deploy --network mainnet
```

## Configuration

Configuration is managed through:
- `config/config.yaml` - Application settings
- `config/contexts/devnet.yaml` - Network configurations
- Environment variables for API keys

### Environment Variables

```bash
# Weather API Keys (optional - Open-Meteo works without key)
export OPENWEATHERMAP_API_KEY="your-key"
export WEATHERAPI_API_KEY="your-key"
export TOMORROWIO_API_KEY="your-key"
export VISUALCROSSING_API_KEY="your-key"
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Process

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Add tests
5. Submit a pull request

## Resources

- [Documentation](docs/)
- [Insurance Demo Guide](INSURANCE_DEMO_README.md)
- [DevKit Usage](DEVKIT_USAGE.md)
- [EigenLayer Docs](https://docs.eigenlayer.xyz)
- [Hourglass Framework](https://github.com/Layr-Labs/hourglass-monorepo)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Built with:
- [EigenLayer](https://eigenlayer.xyz) - Restaking infrastructure
- [Hourglass](https://github.com/Layr-Labs/hourglass-monorepo) - AVS framework
- Weather data providers: Open-Meteo, OpenWeatherMap, and others

---

**Note**: This is a demonstration AVS showcasing parametric insurance capabilities. For production use, ensure proper licensing, regulatory compliance, and thorough testing.