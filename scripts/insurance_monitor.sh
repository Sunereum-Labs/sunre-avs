#!/bin/bash

# Insurance Contract Monitoring Simulation
# This demonstrates how insurance contracts would interact with the AVS

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}"
cat << "EOF"
 ___                                            
|_ _|_ __  ___ _   _ _ __ __ _ _ __   ___ ___  
 | || '_ \/ __| | | | '__/ _` | '_ \ / __/ _ \ 
 | || | | \__ \ |_| | | | (_| | | | | (_|  __/ 
|___|_| |_|___/\__,_|_|  \__,_|_| |_|\___\___| 
                                                
        Monitoring & Claim Automation Demo
EOF
echo -e "${NC}"

# Function to create monitoring task
create_monitoring_task() {
    local policy_id=$1
    local location_lat=$2
    local location_lon=$3
    local threshold=$4
    local peril=$5
    
    cat << EOF
{
  "type": "weather_check",
  "policy_id": "$policy_id",
  "location": {
    "latitude": $location_lat,
    "longitude": $location_lon,
    "city": "Monitored Location",
    "country": "USA"
  },
  "threshold": $threshold,
  "peril": "$peril",
  "monitoring_interval": 3600,
  "callback_contract": "0x1234567890123456789012345678901234567890"
}
EOF
}

# Function to create claim verification task
create_claim_task() {
    local policy_id=$1
    local claim_type=$2
    
    cat << EOF
{
  "type": "insurance_claim",
  "claim_request": {
    "policy_id": "$policy_id",
    "policy": {
      "policy_id": "$policy_id",
      "policy_holder": "Demo Policyholder",
      "insurance_type": "$claim_type",
      "location": {
        "latitude": 40.7128,
        "longitude": -74.0060,
        "city": "New York",
        "country": "USA"
      },
      "coverage_amount": 1000000,
      "premium": 50000,
      "start_date": "2024-01-01T00:00:00Z",
      "end_date": "2024-12-31T00:00:00Z",
      "triggers": [{
        "trigger_id": "AUTO-TRIGGER-001",
        "peril": "heat_wave",
        "conditions": {
          "temperature_max": 35,
          "consecutive_days": 3
        },
        "payout_ratio": 0.5,
        "description": "Automated heat wave protection"
      }]
    },
    "claim_date": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "automated_check": true
  }
}
EOF
}

# Function to submit task to AVS
submit_task() {
    local task_json=$1
    local task_type=$2
    
    echo -e "${BLUE}Submitting $task_type task to AVS...${NC}"
    
    # Encode to base64
    local payload=$(echo -n "$task_json" | base64)
    
    # Submit using devkit
    if devkit avs call --payload "$payload" 2>&1 | tee -a monitoring.log; then
        echo -e "${GREEN}✓ Task submitted successfully${NC}"
        return 0
    else
        echo -e "${RED}✗ Task submission failed${NC}"
        return 1
    fi
}

# Main monitoring loop
echo -e "${BLUE}Starting Insurance Monitoring Demo${NC}"
echo -e "${YELLOW}This simulates how insurance contracts would monitor conditions${NC}\n"

# Example policies to monitor
declare -A policies
policies["FARM-001"]="40.7128,-74.0060,35,heat_wave,crop"
policies["EVENT-001"]="34.0522,-118.2437,50,excess_rain,event"
policies["TRAVEL-001"]="51.5074,-0.1278,-10,cold_snap,travel"

echo -e "${GREEN}Monitoring ${#policies[@]} active policies:${NC}"
for policy_id in "${!policies[@]}"; do
    IFS=',' read -r lat lon threshold peril type <<< "${policies[$policy_id]}"
    echo "  • $policy_id: $type insurance, monitoring for $peril"
done
echo ""

# Monitoring interval (seconds)
INTERVAL=60

# Run monitoring
iteration=0
while true; do
    iteration=$((iteration + 1))
    echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}Monitoring Cycle #$iteration - $(date)${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════${NC}"
    
    for policy_id in "${!policies[@]}"; do
        IFS=',' read -r lat lon threshold peril type <<< "${policies[$policy_id]}"
        
        echo -e "\n${GREEN}Policy: $policy_id${NC}"
        
        # Create and submit monitoring task
        monitoring_task=$(create_monitoring_task "$policy_id" "$lat" "$lon" "$threshold" "$peril")
        
        if submit_task "$monitoring_task" "monitoring"; then
            echo -e "${YELLOW}→ AVS will check if conditions exceed threshold${NC}"
            echo -e "${YELLOW}→ If triggered, claim will be automatically processed${NC}"
            
            # Simulate random trigger (10% chance)
            if [ $((RANDOM % 10)) -eq 0 ]; then
                echo -e "\n${RED}⚠️  ALERT: Threshold exceeded for $policy_id!${NC}"
                echo -e "${YELLOW}Initiating automatic claim...${NC}"
                
                # Create and submit claim task
                claim_task=$(create_claim_task "$policy_id" "$type")
                submit_task "$claim_task" "claim verification"
                
                echo -e "${GREEN}✓ Claim submitted for automated processing${NC}"
            else
                echo -e "${GREEN}✓ Conditions normal - no action needed${NC}"
            fi
        fi
    done
    
    echo -e "\n${CYAN}Next check in $INTERVAL seconds...${NC}"
    echo -e "${YELLOW}Press Ctrl+C to stop monitoring${NC}\n"
    
    sleep $INTERVAL
done