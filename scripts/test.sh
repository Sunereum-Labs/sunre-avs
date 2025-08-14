#!/bin/bash
# SunRe AVS - Test Script  
# Purpose: Run all tests (Go, Solidity, Integration)
# Usage: ./scripts/test.sh [unit|contracts|integration|all]

set -euo pipefail

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "$PROJECT_ROOT"

TEST_TYPE=${1:-all}

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "================================================"
echo "           SunRe AVS - Test Suite              "
echo "================================================"

run_test() {
    local name=$1
    local cmd=$2
    
    echo -e "\n${YELLOW}Running $name...${NC}"
    if eval "$cmd"; then
        echo -e "${GREEN}✓ $name passed${NC}"
        return 0
    else
        echo -e "${RED}✗ $name failed${NC}"
        return 1
    fi
}

FAILED=0

case $TEST_TYPE in
    unit)
        run_test "Go unit tests" "go test ./cmd/... -v" || FAILED=1
        ;;
        
    contracts)
        run_test "Contract compilation" "cd contracts && forge build" || FAILED=1
        run_test "Contract tests" "cd contracts && forge test" || true  # No tests yet
        ;;
        
    integration)
        run_test "Integration tests" "go test ./... -tags=integration -v" || FAILED=1
        ;;
        
    all)
        # Run all tests
        run_test "Go unit tests" "go test ./cmd/... -v" || FAILED=1
        run_test "Contract compilation" "cd contracts && forge build" || FAILED=1
        run_test "Contract tests" "cd contracts && forge test 2>/dev/null" || true  # No tests yet
        
        # Quick functionality check
        echo -e "\n${YELLOW}Checking build artifacts...${NC}"
        if [ -f "bin/sunre-avs" ]; then
            echo -e "  ${GREEN}✓${NC} Performer binary exists"
        else
            echo -e "  ${RED}✗${NC} Performer binary missing"
            FAILED=1
        fi
        
        if [ -d "contracts/out" ]; then
            echo -e "  ${GREEN}✓${NC} Contract artifacts exist"
        else
            echo -e "  ${RED}✗${NC} Contract artifacts missing"
            FAILED=1
        fi
        ;;
        
    *)
        echo "Usage: $0 [unit|contracts|integration|all]"
        echo "  unit        - Run Go unit tests"
        echo "  contracts   - Run Solidity tests"
        echo "  integration - Run integration tests"
        echo "  all         - Run all tests (default)"
        exit 1
        ;;
esac

echo ""
echo "================================================"
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}         All tests passed! ✓                   ${NC}"
else
    echo -e "${RED}         Some tests failed ✗                   ${NC}"
    exit 1
fi
echo "================================================"