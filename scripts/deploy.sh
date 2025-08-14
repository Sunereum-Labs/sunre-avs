#!/bin/bash
# SunRe AVS - Deploy Script
# Purpose: Deploy contracts to specified network
# Usage: ./scripts/deploy.sh [local|testnet|mainnet]

set -euo pipefail

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "$PROJECT_ROOT"

NETWORK=${1:-local}

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "================================================"
echo "         SunRe AVS - Contract Deployment       "
echo "================================================"
echo ""
echo "Target Network: $NETWORK"
echo ""

# Validate environment
check_env() {
    local var=$1
    local msg=$2
    if [ -z "${!var:-}" ]; then
        echo -e "${RED}Error: $var not set${NC}"
        echo "$msg"
        exit 1
    fi
}

deploy_contracts() {
    echo -e "${YELLOW}Deploying contracts...${NC}\n"
    
    case $NETWORK in
        local)
            echo "Using local devnet..."
            RPC_URL="http://localhost:8545"
            
            # Deploy with DevKit
            devkit avs deploy
            ;;
            
        testnet)
            echo "Deploying to Holesky testnet..."
            
            # Check required env vars
            check_env "PRIVATE_KEY_DEPLOYER" "Set your deployer private key"
            check_env "TESTNET_RPC_URL" "Set your testnet RPC URL"
            
            # Check balance
            echo "Checking deployer balance..."
            DEPLOYER=$(cast wallet address $PRIVATE_KEY_DEPLOYER)
            BALANCE=$(cast balance $DEPLOYER --rpc-url $TESTNET_RPC_URL)
            echo "  Deployer: $DEPLOYER"
            echo "  Balance: $(cast to-unit $BALANCE ether) ETH"
            
            if [ "$BALANCE" -lt "100000000000000000" ]; then
                echo -e "${RED}Insufficient balance. Need at least 0.1 ETH${NC}"
                echo "Get testnet ETH from: https://holesky-faucet.pk910.de/"
                exit 1
            fi
            
            # Deploy contracts
            cd contracts
            forge script script/DeployAVS.s.sol:DeployAVS \
                --rpc-url $TESTNET_RPC_URL \
                --private-key $PRIVATE_KEY_DEPLOYER \
                --broadcast \
                --verify \
                ${ETHERSCAN_API_KEY:+--etherscan-api-key $ETHERSCAN_API_KEY}
            cd ..
            ;;
            
        mainnet)
            echo -e "${RED}Mainnet deployment requires additional confirmation${NC}"
            echo ""
            read -p "Are you sure you want to deploy to mainnet? (yes/no): " confirm
            
            if [ "$confirm" != "yes" ]; then
                echo "Deployment cancelled"
                exit 0
            fi
            
            # Check required env vars
            check_env "PRIVATE_KEY_DEPLOYER" "Set your deployer private key"
            check_env "MAINNET_RPC_URL" "Set your mainnet RPC URL"
            
            # Deploy with extra caution
            cd contracts
            forge script script/DeployAVS.s.sol:DeployAVS \
                --rpc-url $MAINNET_RPC_URL \
                --private-key $PRIVATE_KEY_DEPLOYER \
                --broadcast \
                --verify \
                --slow \
                ${ETHERSCAN_API_KEY:+--etherscan-api-key $ETHERSCAN_API_KEY}
            cd ..
            ;;
            
        *)
            echo -e "${RED}Invalid network: $NETWORK${NC}"
            echo "Usage: $0 [local|testnet|mainnet]"
            exit 1
            ;;
    esac
}

# Deploy contracts
deploy_contracts

# Post-deployment
echo ""
echo "================================================"
echo -e "${GREEN}      Deployment Complete! 🚀                  ${NC}"
echo "================================================"
echo ""

case $NETWORK in
    local)
        echo "Next steps:"
        echo "  1. Start performer: ./bin/sunre-avs"
        echo "  2. Submit task: devkit avs call --input examples/task-weather-nyc.json"
        ;;
    testnet)
        echo "Deployment info saved to: contracts/broadcast/"
        echo ""
        echo "Next steps:"
        echo "  1. Register operators: devkit avs operator register"
        echo "  2. Start performer: ENV=production ./bin/sunre-avs"
        echo "  3. Submit task: devkit avs call --context testnet"
        ;;
    mainnet)
        echo "⚠️  MAINNET DEPLOYMENT COMPLETE"
        echo "Verify all contracts on Etherscan before proceeding!"
        ;;
esac