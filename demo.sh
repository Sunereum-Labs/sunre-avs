#!/bin/bash

# SunRe AVS - Unified Demo Launcher
# This script provides all demo options in one place

set -e

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

# ASCII Art Banner
show_banner() {
    echo -e "${CYAN}"
    cat << "EOF"
   _____ __  ___   ______     ___   _   _______
  / ___// / / / | / / __ \   /   | | | / / ___/
  \__ \/ / / /  |/ / /_/ /  / /| | | |/ /\__ \ 
 ___/ / /_/ / /|  / _, _/  / ___ | |   /___/ / 
/____/\____/_/ |_/_/ |_|  /_/  |_| |___//____/  
                                                
    Decentralized Weather Insurance Platform
           Powered by EigenLayer AVS
EOF
    echo -e "${NC}"
}

# Function to check prerequisites
check_prerequisites() {
    local missing=0
    
    echo -e "${BLUE}Checking prerequisites...${NC}"
    
    if ! command -v devkit >/dev/null 2>&1; then
        echo -e "${RED}✗ devkit not found${NC}"
        echo "  Install: https://github.com/Layr-Labs/eigenlayer-devkit"
        missing=1
    else
        echo -e "${GREEN}✓ devkit${NC}"
    fi
    
    if ! command -v docker >/dev/null 2>&1; then
        echo -e "${RED}✗ docker not found${NC}"
        missing=1
    else
        echo -e "${GREEN}✓ docker${NC}"
    fi
    
    if ! command -v go >/dev/null 2>&1; then
        echo -e "${RED}✗ go not found${NC}"
        missing=1
    else
        echo -e "${GREEN}✓ go${NC}"
    fi
    
    if ! command -v node >/dev/null 2>&1; then
        echo -e "${RED}✗ node not found${NC}"
        missing=1
    else
        echo -e "${GREEN}✓ node${NC}"
    fi
    
    return $missing
}

# Function to show menu
show_menu() {
    echo -e "${BLUE}Choose a demo option:${NC}"
    echo ""
    echo "  1) Full DevNet Demo (AVS + Task Submission)"
    echo "  2) UI Demo Only (Mock Mode)"
    echo "  3) Insurance Monitoring Simulation"
    echo "  4) Submit Custom Task"
    echo "  5) View Documentation"
    echo "  6) Clean Up Everything"
    echo "  0) Exit"
    echo ""
}

# Full DevNet Demo
run_devnet_demo() {
    echo -e "${BLUE}Starting Full DevNet Demo...${NC}"
    ./scripts/start_avs_devnet.sh
}

# UI Demo
run_ui_demo() {
    echo -e "${BLUE}Starting UI Demo...${NC}"
    echo -e "${YELLOW}This runs in mock mode without requiring the AVS backend${NC}"
    
    cd demo-ui
    if [ ! -d "node_modules" ]; then
        echo "Installing dependencies..."
        npm install
    fi
    
    echo -e "${GREEN}Starting demo UI on http://localhost:3000${NC}"
    npm start
}

# Insurance Monitoring
run_monitoring() {
    echo -e "${BLUE}Starting Insurance Monitoring Simulation...${NC}"
    echo -e "${YELLOW}Make sure DevNet is running first (Option 1)${NC}"
    echo ""
    read -p "Is DevNet running? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        ./scripts/insurance_monitor.sh
    else
        echo -e "${YELLOW}Please start DevNet first using option 1${NC}"
    fi
}

# Submit Task
submit_task() {
    echo -e "${BLUE}Task Submission Options:${NC}"
    ./scripts/submit_task_demo.sh
    
    echo ""
    read -p "Submit a demo task? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        TASK='{"type":"live_weather_demo","location":{"latitude":40.7128,"longitude":-74.0060,"city":"New York","country":"USA"}}'
        PAYLOAD=$(echo -n "$TASK" | base64)
        echo -e "${YELLOW}Submitting live weather demo task...${NC}"
        devkit avs call --payload "$PAYLOAD"
    fi
}

# View Documentation
view_docs() {
    echo -e "${BLUE}Available Documentation:${NC}"
    echo ""
    echo "1) README.md - Project overview"
    echo "2) DEVNET_DEMO.md - DevNet usage guide"
    echo "3) PRODUCTION_DEPLOYMENT.md - Production deployment"
    echo "4) TESTING.md - Testing guide"
    echo ""
    read -p "Which document? (1-4) " choice
    
    case $choice in
        1) less README.md ;;
        2) less DEVNET_DEMO.md ;;
        3) less PRODUCTION_DEPLOYMENT.md ;;
        4) less TESTING.md ;;
        *) echo "Invalid choice" ;;
    esac
}

# Clean up
cleanup() {
    echo -e "${YELLOW}Cleaning up...${NC}"
    
    # Stop DevNet
    devkit avs devnet stop 2>/dev/null || true
    
    # Kill processes
    pkill -f performer 2>/dev/null || true
    pkill -f "npm start" 2>/dev/null || true
    
    # Clean build artifacts
    rm -rf bin/
    rm -f performer.log monitoring.log task_result.log
    
    echo -e "${GREEN}✓ Cleanup complete${NC}"
}

# Main loop
main() {
    show_banner
    
    if ! check_prerequisites; then
        echo -e "${RED}Please install missing prerequisites first${NC}"
        exit 1
    fi
    
    while true; do
        echo ""
        show_menu
        read -p "Enter choice: " choice
        
        case $choice in
            1) run_devnet_demo ;;
            2) run_ui_demo ;;
            3) run_monitoring ;;
            4) submit_task ;;
            5) view_docs ;;
            6) cleanup ;;
            0) echo "Goodbye!"; exit 0 ;;
            *) echo -e "${RED}Invalid choice${NC}" ;;
        esac
        
        echo ""
        read -p "Press Enter to continue..."
    done
}

# Run main
main