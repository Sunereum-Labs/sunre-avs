#!/bin/bash

# SunRe AVS DevNet Setup Script
# This script starts the local devnet and demonstrates task submission

set -e

# Color codes for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}"
cat << "EOF"
   _____ __  ___   ______     ___   _   _______
  / ___// / / / | / / __ \   /   | | | / / ___/
  \__ \/ / / /  |/ / /_/ /  / /| | | |/ /\__ \ 
 ___/ / /_/ / /|  / _, _/  / ___ | |   /___/ / 
/____/\____/_/ |_/_/ |_|  /_/  |_| |___//____/  
                                                
    AVS DevNet Setup & Task Submission Demo
EOF
echo -e "${NC}"

# Function to cleanup on exit
cleanup() {
    echo -e "\n${YELLOW}Cleaning up...${NC}"
    devkit avs devnet stop 2>/dev/null || true
    pkill -f performer 2>/dev/null || true
    echo -e "${GREEN}✓ Cleanup complete${NC}"
}

trap cleanup EXIT

# Step 1: Stop any existing devnet
echo -e "${BLUE}Step 1: Cleaning up any existing devnet...${NC}"
devkit avs devnet stop 2>/dev/null || true
sleep 2

# Step 2: Start the devnet
echo -e "\n${BLUE}Step 2: Starting local devnet...${NC}"
echo -e "${YELLOW}This will start local blockchain nodes and deploy contracts${NC}"

devkit avs devnet start &
DEVNET_PID=$!

# Wait for devnet to be ready
echo -n "Waiting for devnet to start"
for i in {1..60}; do
    if devkit avs devnet list 2>/dev/null | grep -q "anvil"; then
        echo -e "\n${GREEN}✓ DevNet is running${NC}"
        break
    fi
    echo -n "."
    sleep 2
    
    if [ $i -eq 60 ]; then
        echo -e "\n${RED}✗ DevNet failed to start${NC}"
        exit 1
    fi
done

# Step 3: Deploy contracts (if needed)
echo -e "\n${BLUE}Step 3: Deploying contracts...${NC}"
if devkit avs devnet deploy-contracts; then
    echo -e "${GREEN}✓ Contracts deployed${NC}"
else
    echo -e "${YELLOW}! Contracts may already be deployed${NC}"
fi

# Step 4: Build and start the AVS performer
echo -e "\n${BLUE}Step 4: Starting AVS performer...${NC}"
if [ ! -f "./bin/performer" ]; then
    echo "Building performer..."
    make build
fi

./bin/performer --port 8080 > performer.log 2>&1 &
PERFORMER_PID=$!
echo -e "${GREEN}✓ Performer started (PID: $PERFORMER_PID)${NC}"

# Wait for performer to be ready
sleep 5

# Step 5: Show devnet information
echo -e "\n${BLUE}Step 5: DevNet Information${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
devkit avs devnet list
echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"

# Step 6: Create example task payloads
echo -e "\n${BLUE}Step 6: Example Task Submission${NC}"
echo -e "${YELLOW}In a production scenario, insurance contracts would submit these tasks${NC}\n"

# Weather monitoring task
WEATHER_TASK=$(echo -n '{
  "type": "weather_check",
  "location": {
    "latitude": 40.7128,
    "longitude": -74.0060,
    "city": "New York",
    "country": "USA"
  },
  "threshold": 35.0,
  "policy_id": "PROD-POLICY-001",
  "monitoring_interval": 3600
}' | base64)

# Insurance claim verification task
CLAIM_TASK=$(echo -n '{
  "type": "insurance_claim",
  "claim_request": {
    "policy_id": "PROD-POLICY-001",
    "policy": {
      "policy_id": "PROD-POLICY-001",
      "policy_holder": "NYC Farm Co.",
      "insurance_type": "crop",
      "location": {
        "latitude": 40.7128,
        "longitude": -74.0060,
        "city": "New York",
        "country": "USA"
      },
      "coverage_amount": 500000,
      "premium": 25000,
      "start_date": "2024-01-01T00:00:00Z",
      "end_date": "2024-12-31T00:00:00Z",
      "triggers": [{
        "trigger_id": "HEAT-3DAY",
        "peril": "heat_wave",
        "conditions": {
          "temperature_max": 35,
          "consecutive_days": 3
        },
        "payout_ratio": 0.5,
        "description": "Heat wave protection"
      }]
    },
    "claim_date": "2024-07-15T00:00:00Z",
    "automated_check": true
  }
}' | base64)

echo -e "${GREEN}Example 1: Weather Monitoring Task${NC}"
echo "This task would be submitted periodically by the insurance contract:"
echo -e "${CYAN}devkit avs call --payload $WEATHER_TASK${NC}"
echo ""

echo -e "${GREEN}Example 2: Insurance Claim Verification${NC}"
echo "This task is submitted when a claim needs verification:"
echo -e "${CYAN}devkit avs call --payload $CLAIM_TASK${NC}"
echo ""

# Step 7: Submit a test task
echo -e "${BLUE}Step 7: Submitting test weather monitoring task...${NC}"
if devkit avs call --payload "$WEATHER_TASK" 2>&1 | tee task_result.log; then
    echo -e "${GREEN}✓ Task submitted successfully${NC}"
    echo "Check task_result.log for details"
else
    echo -e "${YELLOW}! Task submission failed - this is normal if the AVS is still initializing${NC}"
fi

# Step 8: Production deployment notes
echo -e "\n${BLUE}Production Deployment Steps:${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
echo ""
echo -e "${GREEN}1. Testnet Deployment:${NC}"
echo "   • Deploy contracts to Holesky testnet"
echo "   • Register operators with testnet EigenLayer"
echo "   • Configure real API endpoints and keys"
echo "   • Set up monitoring infrastructure"
echo ""
echo -e "${GREEN}2. Smart Contract Integration:${NC}"
echo "   • Insurance contracts call AVS for weather verification"
echo "   • AVS monitors conditions based on policy triggers"
echo "   • Automatic claim processing when conditions met"
echo "   • BLS signature aggregation for consensus"
echo ""
echo -e "${GREEN}3. Production Checklist:${NC}"
echo "   • Secure key management (AWS KMS/HSM)"
echo "   • Rate limiting for API calls"
echo "   • Redundant data sources (minimum 5)"
echo "   • Monitoring and alerting"
echo "   • Disaster recovery plan"
echo ""

echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✨ DevNet is ready for testing!${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
echo ""
echo -e "Logs available at:"
echo -e "  • Performer: ${CYAN}performer.log${NC}"
echo -e "  • DevKit: ${CYAN}devkit avs logs performer${NC}"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop the devnet${NC}"

# Keep running
wait