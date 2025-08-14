# SunRe AVS Scripts

Clean, minimal scripts for development and deployment.

## Scripts

| Script | Purpose | Usage |
|--------|---------|-------|
| `setup.sh` | Install dependencies and build project | `./scripts/setup.sh` |
| `start.sh` | Start local development environment | `./scripts/start.sh [devnet\|docker]` |
| `test.sh` | Run test suite | `./scripts/test.sh [unit\|contracts\|all]` |
| `deploy.sh` | Deploy contracts to network | `./scripts/deploy.sh [local\|testnet\|mainnet]` |

## Quick Start

```bash
# First time setup
./scripts/setup.sh

# Start local development
./scripts/start.sh

# Run tests
./scripts/test.sh

# Deploy to testnet
./scripts/deploy.sh testnet
```

## Environment Variables

For testnet/mainnet deployment, set these in your `.env`:

```bash
PRIVATE_KEY_DEPLOYER=your-private-key
TESTNET_RPC_URL=https://ethereum-holesky.publicnode.com
ETHERSCAN_API_KEY=your-etherscan-key
```

All scripts follow DevKit conventions and include proper error handling.