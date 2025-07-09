# Temperature Oracle AVS

A decentralized temperature verification oracle system built following the Refiant AVS (Actively Validated Service) architecture. The oracle verifies if a location's temperature meets a specified threshold using multiple weather data sources with MAD (Median Absolute Deviation) consensus.

## Features

- **Multi-source verification**: Integrates 5 weather APIs for redundancy
- **MAD consensus algorithm**: Filters outliers and reaches reliable consensus
- **Rate limiting**: Token bucket algorithm for API rate management
- **Caching**: Reduces API calls with intelligent caching
- **Parallel execution**: Efficient concurrent data fetching
- **Prometheus metrics**: Built-in monitoring and observability
- **Production-ready**: Comprehensive error handling and logging

## Architecture

The system follows the Refiant AVS architecture with these components:

- **Aggregator**: Manages task distribution and response collection
- **Executor**: Fetches temperature data from assigned sources
- **Consensus Engine**: Implements MAD algorithm for reliable consensus
- **Data Sources**: Weather API integrations with rate limiting

## Installation

```bash
# Clone the repository
cd temperature-oracle

# Install dependencies
go mod tidy

# Set up API keys (optional - Open-Meteo works without key)
export OPENWEATHERMAP_API_KEY="your-key"
export WEATHERAPI_API_KEY="your-key"
export TOMORROWIO_API_KEY="your-key"
export VISUALCROSSING_API_KEY="your-key"
```

## Usage

Basic usage:
```bash
go run main.go --location "New York" --threshold 25.0
```

With custom configuration:
```bash
go run main.go --config config/config.yaml --location "London" --threshold 20.0 --log-level debug
```

Supported locations:
- City names: "New York", "London", "Tokyo", "Paris", "Sydney", "San Francisco", "Singapore", "Dubai"
- Coordinates: "40.7128,-74.0060"

## Configuration

The `config/config.yaml` file controls:
- Minimum operators and consensus threshold
- API rate limits and endpoints
- MAD threshold and cache TTL
- Response timeouts

## Example Output

```
=== Temperature Verification Result ===
Location: New York (40.71, -74.01)
Consensus Temperature: 23.45°C
Threshold: 25.00°C
Meets Threshold: false
Confidence: 92.50%
Data Sources Used: 4

Data Points:
  - OpenMeteo: 23.20°C (confidence: 0.95)
  - WeatherAPI: 23.50°C (confidence: 0.93)
  - OpenWeatherMap: 23.40°C (confidence: 0.95)
  - TomorrowIO: 23.70°C (confidence: 0.92)

Task ID: task_1234567890_1
Aggregated Signature: a1b2c3d4e5f6...
=====================================
```

## Metrics

Prometheus metrics available at `http://localhost:8080/metrics`:
- `temperature_oracle_tasks_processed_total`: Task completion counter
- `temperature_oracle_task_duration_seconds`: Phase-wise latency histogram
- `temperature_oracle_consensus_temperature`: Current consensus temperature gauge

## Performance

Meets Refiant AVS latency requirements (<2 minutes total):
- Event detection: 15-30s
- Task distribution: 5-10s
- Parallel data fetching: 30-60s
- Aggregation & verification: 10-15s

## Development

Run tests:
```bash
go test ./...
```

Build binary:
```bash
go build -o temperature-oracle
```

Docker support:
```bash
docker build -t temperature-oracle .
docker run -e OPENMETEO_API_KEY="" temperature-oracle --location "Tokyo" --threshold 30
```