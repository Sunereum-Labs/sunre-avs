# SunRe AVS Production Deployment Guide

## Overview

This guide outlines the steps to deploy the SunRe Weather Insurance AVS from local devnet to production on Ethereum mainnet.

## Architecture Overview

```
┌─────────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│ Insurance Contract  │────▶│   SunRe AVS      │────▶│ Weather APIs    │
│ (Monitors policies) │     │ (Validates data) │     │ (3+ sources)    │
└─────────────────────┘     └──────────────────┘     └─────────────────┘
         │                           │
         ▼                           ▼
┌─────────────────────┐     ┌──────────────────┐
│  Claim Processor    │◀────│ EigenLayer Core  │
│ (Automated payouts) │     │ (BLS signatures) │
└─────────────────────┘     └──────────────────┘
```

## Deployment Phases

### Phase 1: Local DevNet Testing ✅

```bash
# Start local devnet
./scripts/start_avs_devnet.sh

# Test task submission
devkit avs call --payload <base64-task>

# Monitor logs
devkit avs logs performer
```

### Phase 2: Holesky Testnet Deployment

#### 2.1 Prerequisites
- [ ] Holesky ETH for gas
- [ ] Operator keys (ECDSA + BLS)
- [ ] RPC endpoints

#### 2.2 Contract Deployment
```bash
# Update config/contexts/testnet.yaml
cp config/contexts/devnet.yaml config/contexts/testnet.yaml
# Edit with testnet RPC URLs and addresses

# Deploy contracts
forge script script/testnet/DeployAVSL1Contracts.s.sol \
  --rpc-url $HOLESKY_RPC \
  --private-key $DEPLOYER_KEY \
  --broadcast

# Verify contracts
forge verify-contract <CONTRACT_ADDRESS> TaskAVSRegistrar \
  --chain holesky
```

#### 2.3 Operator Registration
```bash
# Register operators with EigenLayer
devkit avs register-operator \
  --config config/contexts/testnet.yaml \
  --operator-key $OPERATOR_KEY
```

#### 2.4 AVS Configuration
```yaml
# config/avs-testnet.yaml
avs:
  chain_id: 17000  # Holesky
  rpc_url: "https://holesky.infura.io/v3/YOUR_KEY"
  contracts:
    avs_registrar: "0x..."
    task_mailbox: "0x..."
  
weather_apis:
  - name: "tomorrowio"
    endpoint: "https://api.tomorrow.io/v4"
    api_key: "${TOMORROW_IO_KEY}"
  - name: "weatherapi"
    endpoint: "https://api.weatherapi.com/v1"
    api_key: "${WEATHER_API_KEY}"
  - name: "openmeteo"
    endpoint: "https://api.open-meteo.com/v1"
```

### Phase 3: Mainnet Production Deployment

#### 3.1 Security Checklist

- [ ] **Key Management**
  ```bash
  # Use AWS KMS or hardware security modules
  export OPERATOR_KEY_ARN="arn:aws:kms:..."
  ```

- [ ] **API Key Rotation**
  ```yaml
  # Store in AWS Secrets Manager
  weather_apis:
    api_key: "${aws:secretsmanager:sunre-api-keys:tomorrowio}"
  ```

- [ ] **Rate Limiting**
  ```go
  // internal/datasources/ratelimiter.go
  limiter := rate.NewLimiter(rate.Every(time.Second), 10)
  ```

- [ ] **Monitoring**
  ```yaml
  # prometheus.yaml
  scrape_configs:
    - job_name: 'avs-metrics'
      static_configs:
        - targets: ['localhost:9090']
  ```

#### 3.2 Smart Contract Integration

```solidity
// Example Insurance Contract Integration
contract WeatherInsurance {
    ITaskMailbox public taskMailbox;
    mapping(bytes32 => Policy) public policies;
    
    function requestWeatherCheck(bytes32 policyId) external {
        Policy memory policy = policies[policyId];
        
        // Create monitoring task
        bytes memory taskData = abi.encode(
            "weather_check",
            policy.location,
            policy.threshold,
            policyId
        );
        
        // Submit to AVS
        taskMailbox.submitTask(taskData);
    }
    
    // Called by AVS when conditions are met
    function triggerClaim(
        bytes32 policyId,
        uint256 temperature,
        bytes calldata proof
    ) external onlyAVS {
        require(verifyProof(proof), "Invalid proof");
        
        Policy memory policy = policies[policyId];
        if (temperature > policy.threshold) {
            processPayout(policyId, policy.coverageAmount);
        }
    }
}
```

#### 3.3 Deployment Script

```bash
#!/bin/bash
# deploy-mainnet.sh

# 1. Deploy contracts
forge script script/mainnet/DeployAVSL1Contracts.s.sol \
  --rpc-url $MAINNET_RPC \
  --private-key $DEPLOYER_KEY \
  --broadcast \
  --verify

# 2. Register with EigenLayer
devkit avs register \
  --network mainnet \
  --avs-address $AVS_ADDRESS

# 3. Start operators
docker-compose -f docker-compose.prod.yml up -d

# 4. Monitor health
curl https://api.sunre-avs.com/health
```

## Production Architecture

### High Availability Setup

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  performer-1:
    image: sunre/avs-performer:latest
    environment:
      - NODE_ENV=production
      - OPERATOR_KEY_ARN=${OPERATOR_KEY_ARN}
    deploy:
      replicas: 3
      
  aggregator:
    image: sunre/avs-aggregator:latest
    deploy:
      replicas: 2
      
  monitoring:
    image: prom/prometheus:latest
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
```

### Database Schema

```sql
-- Production database for audit trail
CREATE TABLE weather_checks (
    id UUID PRIMARY KEY,
    policy_id VARCHAR(66) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    location JSONB NOT NULL,
    consensus_temperature DECIMAL(5,2),
    confidence DECIMAL(3,2),
    data_sources JSONB,
    signature BYTEA,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_policy_timestamp ON weather_checks(policy_id, timestamp);
```

## Monitoring & Alerts

### Key Metrics

1. **AVS Health**
   - Task success rate > 99%
   - Response time < 5s
   - Operator uptime > 99.9%

2. **Data Quality**
   - Consensus confidence > 0.8
   - Active data sources >= 3
   - Outlier rate < 10%

3. **Business Metrics**
   - Claims processed/hour
   - Average claim value
   - False positive rate

### Alert Configuration

```yaml
# alerts.yaml
groups:
  - name: avs_alerts
    rules:
      - alert: LowConsensusConfidence
        expr: avs_consensus_confidence < 0.7
        for: 5m
        annotations:
          summary: "Low consensus confidence detected"
          
      - alert: DataSourceFailure
        expr: avs_active_sources < 3
        for: 1m
        annotations:
          summary: "Insufficient data sources"
```

## Disaster Recovery

### Backup Strategy

1. **State Backup**
   ```bash
   # Hourly snapshots
   0 * * * * pg_dump avs_db > /backup/avs_$(date +\%Y\%m\%d_\%H).sql
   ```

2. **Key Backup**
   - BLS keys in AWS KMS
   - ECDSA keys in hardware wallet
   - Backup operators in different regions

3. **Failover Plan**
   - Primary: US-East-1
   - Secondary: EU-West-1
   - Automatic failover via Route53

## Cost Optimization

### API Usage

```go
// Implement caching to reduce API calls
cache := ttlcache.New(
    ttlcache.WithTTL[string, float64](5 * time.Minute),
)

// Batch requests where possible
func (w *WeatherOracle) BatchFetch(locations []Location) []DataPoint {
    // Implementation
}
```

### Gas Optimization

```solidity
// Use events instead of storage where possible
event ClaimTriggered(
    bytes32 indexed policyId,
    uint256 temperature,
    uint256 payout
);

// Batch claim processing
function processClaims(bytes32[] calldata policyIds) external {
    // Process multiple claims in one transaction
}
```

## Launch Checklist

- [ ] All contracts deployed and verified
- [ ] Operators registered with sufficient stake
- [ ] Monitoring dashboards configured
- [ ] Disaster recovery tested
- [ ] Security audit completed
- [ ] API rate limits configured
- [ ] Documentation published
- [ ] Support channels established

## Support & Maintenance

- **Documentation**: https://docs.sunre-avs.com
- **Support**: support@sunre-avs.com
- **Emergency**: Use PagerDuty escalation
- **Updates**: Follow semantic versioning

## Next Steps

1. Complete security audit
2. Load testing on testnet
3. Gradual mainnet rollout
4. Partner onboarding
5. Community operator program