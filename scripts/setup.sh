#!/bin/bash
# SunRe AVS - Setup Script
# Purpose: One-command setup for development environment
# Usage: ./scripts/setup.sh

set -euo pipefail

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "$PROJECT_ROOT"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "================================================"
echo "         SunRe AVS - Development Setup          "
echo "================================================"

# Check prerequisites
echo -e "\n${YELLOW}Checking prerequisites...${NC}"

check_command() {
    if command -v $1 &> /dev/null; then
        echo -e "  ${GREEN}✓${NC} $2 installed"
        return 0
    else
        echo -e "  ✗ $2 not found. Please install: $3"
        return 1
    fi
}

MISSING=0
check_command go "Go (1.23+)" "https://go.dev" || MISSING=1
check_command forge "Foundry" "curl -L https://foundry.paradigm.xyz | bash" || MISSING=1
check_command docker "Docker" "https://docker.com" || MISSING=1
check_command devkit "DevKit CLI" "curl -sSL https://install.eigenlayer.xyz | sh" || MISSING=1

if [ $MISSING -eq 1 ]; then
    echo -e "\n${YELLOW}Please install missing prerequisites and run again.${NC}"
    exit 1
fi

# Install dependencies
echo -e "\n${YELLOW}Installing dependencies...${NC}"

echo "  Installing Go dependencies..."
go mod tidy

echo "  Installing contract dependencies..."
cd contracts && forge install && cd ..

# Build everything
echo -e "\n${YELLOW}Building project...${NC}"

echo "  Building contracts..."
cd contracts && forge build && cd ..

echo "  Building performer..."
go build -o bin/sunre-avs cmd/main.go

# Setup complete
echo -e "\n${GREEN}================================================${NC}"
echo -e "${GREEN}           Setup Complete! 🚀                  ${NC}"
echo -e "${GREEN}================================================${NC}"
echo ""
echo "Next steps:"
echo "  1. Start local devnet:  ./scripts/start.sh"
echo "  2. Run tests:           ./scripts/test.sh"
echo "  3. Deploy to testnet:   ./scripts/deploy.sh testnet"
echo ""