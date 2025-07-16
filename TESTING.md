# SunRe AVS Testing Guide

This document describes the testing strategy and available tests for the SunRe AVS.

## Running Tests

### Quick Start

```bash
# Run all unit tests
make test

# Run tests in short mode (skip slow tests)
make test-short

# Run benchmarks
make test-bench

# Run tests with coverage
make test-coverage
```

### Test Structure

The main test file is located at `cmd/main_test.go` and includes:

1. **Unit Tests** - Test individual functions and components
2. **Integration Tests** - Test the full task handling flow
3. **Benchmarks** - Performance testing for critical paths

## Test Categories

### 1. Validation Tests (`TestValidateTask`)

Tests the validation logic for different task types:
- ✅ Valid weather check requests
- ✅ Invalid latitude/longitude values
- ✅ Invalid temperature thresholds
- ✅ Valid insurance claims
- ✅ Missing or invalid policy data
- ✅ Unknown task types

### 2. Task Handling Tests (`TestHandleTask`)

Tests the complete task processing flow:
- ✅ Insurance claim processing with demo data
- ⚠️  Weather check (disabled - requires live APIs)
- ⚠️  Live weather demo (disabled - requires live APIs)

**Note**: Tests that require external API calls are commented out by default to ensure tests run reliably in CI/CD environments.

### 3. Component Tests

- `TestWeatherOracleInitialization` - Verifies oracle setup
- `TestValidateWeatherLocation` - Tests location validation logic
- `TestServerIntegration` - Tests RPC server startup
- `TestConcurrentTaskValidation` - Tests concurrent request handling

### 4. Benchmarks

Performance benchmarks for critical operations:
- `BenchmarkValidateTask` - Measures validation performance (~2.3μs/op)
- `BenchmarkHandleTask` - Measures full task handling (varies with API calls)

## Test Data

### Demo Mode

Tests use demo mode for insurance claims to avoid external dependencies:

```json
{
  "type": "insurance_claim",
  "claim_request": {
    "policy": { ... },
    "claim_date": "2024-07-15T00:00:00Z"
  },
  "demo_mode": true,
  "demo_scenario": "heat_wave"
}
```

### Location Validation

Tests cover edge cases:
- Valid coordinates: 40.7128°N, 74.0060°W (NYC)
- Edge cases: North/South poles (±90° latitude)
- Invalid: Latitudes > 90° or < -90°
- Invalid: Longitudes > 180° or < -180°

## Coverage

Run coverage analysis:

```bash
make test-coverage
# Opens coverage.html in your browser
```

Key areas covered:
- Task validation logic
- Insurance claim processing
- Location validation
- Error handling

## Integration with DevKit

The tests are designed to work with EigenLayer DevKit:

```bash
# Run tests with devkit (requires devkit installed)
devkit avs test

# Or use the Makefile target
make test-integration
```

## Continuous Integration

Tests are designed to run in CI/CD pipelines:
- No external API dependencies in default tests
- Deterministic results with demo data
- Fast execution (~1.2s for full suite)
- Clear error messages for debugging

## Adding New Tests

When adding new functionality:

1. Add validation tests in `TestValidateTask`
2. Add handling tests in `TestHandleTask` 
3. Use demo mode for predictable results
4. Add benchmarks for performance-critical code
5. Document any external dependencies

## Troubleshooting

### Tests Failing Due to API Calls

If you see errors like "too many outliers filtered", the test is trying to make real API calls. Ensure:
- Tests use `demo_mode: true` for insurance claims
- Weather check tests are commented out or use mocks
- API keys are not required for unit tests

### Benchmark Variations

Benchmark results may vary based on:
- Network latency (for API calls)
- System load
- API rate limits

For consistent benchmarks, use mock data sources.

## Best Practices

1. **Isolation**: Tests should not depend on external services
2. **Determinism**: Tests should produce consistent results
3. **Speed**: Keep unit tests fast (<100ms per test)
4. **Coverage**: Aim for >80% code coverage
5. **Documentation**: Document test scenarios and edge cases