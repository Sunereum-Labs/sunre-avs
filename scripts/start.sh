#!/bin/bash
# SunRe AVS - Start Script
# Purpose: Start local development environment
# Usage: ./scripts/start.sh [devnet|docker]

set -euo pipefail

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "$PROJECT_ROOT"

MODE=${1:-devnet}

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "================================================"
echo "         SunRe AVS - Starting Services          "
echo "================================================"

case $MODE in
    devnet)
        echo -e "\n${YELLOW}Starting DevKit devnet...${NC}"
        
        # Start devnet
        devkit avs devnet start &
        DEVNET_PID=$!
        
        echo "Waiting for devnet to be ready..."
        sleep 10
        
        # Deploy contracts
        echo -e "\n${YELLOW}Deploying contracts...${NC}"
        devkit avs deploy
        
        # Start performer
        echo -e "\n${YELLOW}Starting performer...${NC}"
        ./bin/sunre-avs &
        PERFORMER_PID=$!
        
        echo -e "\n${GREEN}Services started successfully!${NC}"
        echo ""
        echo "  Devnet RPC:     http://localhost:8545"
        echo "  Performer:      http://localhost:8080"
        echo "  Health:         http://localhost:8081/health"
        echo ""
        echo "Submit a test task:"
        echo "  devkit avs call --input examples/task-weather-nyc.json"
        echo ""
        echo "Press Ctrl+C to stop all services"
        
        # Wait for interrupt
        trap "kill $DEVNET_PID $PERFORMER_PID 2>/dev/null; devkit avs devnet stop" EXIT
        wait
        ;;
        
    docker)
        echo -e "\n${YELLOW}Starting with Docker Compose...${NC}"
        
        docker-compose up -d
        
        echo -e "\n${GREEN}Services started successfully!${NC}"
        echo ""
        echo "View logs:"
        echo "  docker-compose logs -f"
        echo ""
        echo "Stop services:"
        echo "  docker-compose down"
        ;;
        
    *)
        echo "Usage: $0 [devnet|docker]"
        echo "  devnet - Start local DevKit devnet (default)"
        echo "  docker - Start with Docker Compose"
        exit 1
        ;;
esac