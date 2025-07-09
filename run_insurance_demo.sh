#!/bin/bash

# Color codes for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}Weather Insurance AVS Demo${NC}"
echo -e "${BLUE}================================${NC}"
echo ""

# Function to run command and show status
run_cmd() {
    echo -e "${YELLOW}Running: $1${NC}"
    eval $1
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Success${NC}"
    else
        echo -e "${RED}✗ Failed${NC}"
        exit 1
    fi
    echo ""
}

# Check if devkit is installed
if ! command -v devkit &> /dev/null; then
    echo -e "${RED}Error: 'devkit' command not found${NC}"
    echo "Please install EigenLayer DevKit first"
    exit 1
fi

# Step 1: Build
echo -e "${BLUE}Step 1: Building AVS${NC}"
run_cmd "devkit avs build"

# Step 2: Start DevNet
echo -e "${BLUE}Step 2: Starting DevNet${NC}"
echo -e "${YELLOW}This may take a minute...${NC}"
run_cmd "devkit avs devnet start"

# Wait for services to be ready
echo -e "${YELLOW}Waiting for services to start...${NC}"
sleep 10

# Step 3: Show demo options
echo -e "${BLUE}Step 3: Choose a Demo Scenario${NC}"
echo ""
echo "1) Crop Insurance - Heat Wave (APPROVED claim)"
echo "2) Event Insurance - Normal Weather (REJECTED claim)"
echo "3) Travel Insurance - Cold Snap (APPROVED claim)"
echo "4) Simple Weather Check"
echo ""
read -p "Select option (1-4): " choice

case $choice in
    1)
        echo -e "${GREEN}Running Crop Insurance Demo...${NC}"
        PAYLOAD='eyJ0eXBlIjoiaW5zdXJhbmNlX2NsYWltIiwiY2xhaW1fcmVxdWVzdCI6eyJwb2xpY3lfaWQiOiJDUk9QLTIwMjQtMDAxIiwicG9saWN5Ijp7InBvbGljeV9pZCI6IkNST1AtMjAyNC0wMDEiLCJwb2xpY3lfaG9sZGVyIjoiR3JlZW4gVmFsbGV5IEZhcm1zIiwiaW5zdXJhbmNlX3R5cGUiOiJjcm9wIiwibG9jYXRpb24iOnsibGF0aXR1ZGUiOjM1LjIyNzEsImxvbmdpdHVkZSI6LTgwLjg0MzEsImNpdHkiOiJDaGFybG90dGUiLCJjb3VudHJ5IjoiVVNBIn0sImNvdmVyYWdlX2Ftb3VudCI6MTAwMDAwLCJwcmVtaXVtIjo1MDAwLCJzdGFydF9kYXRlIjoiMjAyNC0wNi0wMVQwMDowMDowMFoiLCJlbmRfZGF0ZSI6IjIwMjQtMDktMzBUMDA6MDA6MDBaIiwidHJpZ2dlcnMiOlt7InRyaWdnZXJfaWQiOiJIRUFULTNEQVktMzVDIiwicGVyaWwiOiJoZWF0X3dhdmUiLCJjb25kaXRpb25zIjp7InRlbXBlcmF0dXJlX21heCI6MzUsImNvbnNlY3V0aXZlX2RheXMiOjMsInRpbWVfd2luZG93Ijp7InN0YXJ0X21vbnRoIjo2LCJlbmRfbW9udGgiOjh9fSwicGF5b3V0X3JhdGlvIjowLjUsImRlc2NyaXB0aW9uIjoiNTAlIHBheW91dCBpZiB0ZW1wZXJhdHVyZSBleGNlZWRzIDM1wrBDIGZvciAzIGNvbnNlY3V0aXZlIGRheXMifSx7InRyaWdnZXJfaWQiOiJFWFRSRU1FLUhFQVQtNDBDIiwicGVyaWwiOiJoZWF0X3dhdmUiLCJjb25kaXRpb25zIjp7InRlbXBlcmF0dXJlX21heCI6NDAsImNvbnNlY3V0aXZlX2RheXMiOjJ9LCJwYXlvdXRfcmF0aW8iOjEuMCwiZGVzY3JpcHRpb24iOiJGdWxsIHBheW91dCBmb3IgZXh0cmVtZSBoZWF0ICg0MMKwQysgZm9yIDIgZGF5cykifV19LCJjbGFpbV9kYXRlIjoiMjAyNC0wNy0xNVQwMDowMDowMFoiLCJhdXRvbWF0ZWRfY2hlY2siOnRydWV9LCJkZW1vX21vZGUiOnRydWUsImRlbW9fc2NlbmFyaW8iOiJoZWF0X3dhdmUifQ=='
        echo -e "${YELLOW}Expected: Approved claim with $50,000 payout${NC}"
        ;;
    2)
        echo -e "${GREEN}Running Event Insurance Demo...${NC}"
        PAYLOAD='eyJ0eXBlIjoiaW5zdXJhbmNlX2NsYWltIiwiY2xhaW1fcmVxdWVzdCI6eyJwb2xpY3lfaWQiOiJFVkVOVC0yMDI0LU1VU0lDIiwicG9saWN5Ijp7InBvbGljeV9pZCI6IkVWRU5ULTIwMjQtTVVTSUMiLCJwb2xpY3lfaG9sZGVyIjoiU3VtbWVyIE11c2ljIEZlc3RpdmFsIExMQyIsImluc3VyYW5jZV90eXBlIjoiZXZlbnQiLCJsb2NhdGlvbiI6eyJsYXRpdHVkZSI6NDAuNzEyOCwibG9uZ2l0dWRlIjotNzQuMDA2MCwiY2l0eSI6Ik5ldyBZb3JrIiwiY291bnRyeSI6IlVTQSJ9LCJjb3ZlcmFnZV9hbW91bnQiOjUwMDAwMCwicHJlbWl1bSI6MjUwMDAsInN0YXJ0X2RhdGUiOiIyMDI0LTA4LTE1VDAwOjAwOjAwWiIsImVuZF9kYXRlIjoiMjAyNC0wOC0xN1QwMDowMDowMFoiLCJ0cmlnZ2VycyI6W3sidHJpZ2dlcl9pZCI6IlJBSU4tNTBNTSIsInBlcmlsIjoiZXhjZXNzX3JhaW4iLCJjb25kaXRpb25zIjp7InByZWNpcGl0YXRpb25fbWluIjo1MCwidGltZV93aW5kb3ciOnsic3RhcnRfaG91ciI6OCwiZW5kX2hvdXIiOjIwfX0sInBheW91dF9yYXRpbyI6MS4wLCJkZXNjcmlwdGlvbiI6IkZ1bGwgcGF5b3V0IGlmIHJhaW4gZXhjZWVkcyA1MG1tIGR1cmluZyBldmVudCBob3VycyJ9LHsidHJpZ2dlcl9pZCI6IldJTkQtNjBLTUgiLCJwZXJpbCI6ImhpZ2hfd2luZCIsImNvbmRpdGlvbnMiOnsid2luZF9zcGVlZF9taW4iOjYwfSwicGF5b3V0X3JhdGlvIjoxLjAsImRlc2NyaXB0aW9uIjoiRnVsbCBwYXlvdXQgZm9yIGRhbmdlcm91cyB3aW5kIGNvbmRpdGlvbnMifV19LCJjbGFpbV9kYXRlIjoiMjAyNC0wOC0xNlQwMDowMDowMFoiLCJhdXRvbWF0ZWRfY2hlY2siOnRydWV9LCJkZW1vX21vZGUiOnRydWUsImRlbW9fc2NlbmFyaW8iOiJub3JtYWwifQ=='
        echo -e "${YELLOW}Expected: Rejected claim (no weather triggers met)${NC}"
        ;;
    3)
        echo -e "${GREEN}Running Travel Insurance Demo...${NC}"
        PAYLOAD='eyJ0eXBlIjoiaW5zdXJhbmNlX2NsYWltIiwiY2xhaW1fcmVxdWVzdCI6eyJwb2xpY3lfaWQiOiJUUkFWRUwtMjAyNC0wMTIzIiwicG9saWN5Ijp7InBvbGljeV9pZCI6IlRSQVZFTC0yMDI0LTAxMjMiLCJwb2xpY3lfaG9sZGVyIjoiSm9obiBEb2UiLCJpbnN1cmFuY2VfdHlwZSI6InRyYXZlbCIsImxvY2F0aW9uIjp7ImxhdGl0dWRlIjo0MS45NzQyLCJsb25naXR1ZGUiOi04Ny45MDczLCJjaXR5IjoiQ2hpY2FnbyBPJ0hhcmUiLCJjb3VudHJ5IjoiVVNBIn0sImNvdmVyYWdlX2Ftb3VudCI6MTAwMCwicHJlbWl1bSI6NTAsInN0YXJ0X2RhdGUiOiIyMDI0LTEyLTIwVDAwOjAwOjAwWiIsImVuZF9kYXRlIjoiMjAyNC0xMi0yNVQwMDowMDowMFoiLCJ0cmlnZ2VycyI6W3sidHJpZ2dlcl9pZCI6IkNPTEQtREVMQVkiLCJwZXJpbCI6ImNvbGRfc25hcCIsImNvbmRpdGlvbnMiOnsidGVtcGVyYXR1cmVfbWF4IjotMTB9LCJwYXlvdXRfcmF0aW8iOjAuMiwiZGVzY3JpcHRpb24iOiJEYWlseSBjb21wZW5zYXRpb24gZm9yIGV4dHJlbWUgY29sZCBkZWxheXMifV19LCJjbGFpbV9kYXRlIjoiMjAyNC0xMi0yMlQwMDowMDowMFoiLCJhdXRvbWF0ZWRfY2hlY2siOnRydWV9LCJkZW1vX21vZGUiOnRydWUsImRlbW9fc2NlbmFyaW8iOiJjb2xkX3NuYXAifQ=='
        echo -e "${YELLOW}Expected: Approved claim with $200 daily compensation${NC}"
        ;;
    4)
        echo -e "${GREEN}Running Simple Weather Check...${NC}"
        PAYLOAD='eyJ0eXBlIjoid2VhdGhlcl9jaGVjayIsImxvY2F0aW9uIjp7ImxhdGl0dWRlIjo0MC43MTI4LCJsb25naXR1ZGUiOi03NC4wMDYwLCJjaXR5IjoiTmV3IFlvcmsiLCJjb3VudHJ5IjoiVVNBIn0sInRocmVzaG9sZCI6MjUuMH0='
        echo -e "${YELLOW}Expected: Current temperature and threshold check${NC}"
        ;;
    *)
        echo -e "${RED}Invalid option${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${BLUE}Submitting task to AVS...${NC}"
echo ""

# Submit the task
devkit avs call --payload "$PAYLOAD"

echo ""
echo -e "${GREEN}Demo completed!${NC}"
echo ""
echo -e "${BLUE}To view logs:${NC}"
echo "devkit avs logs performer"
echo ""
echo -e "${BLUE}To stop the devnet:${NC}"
echo "devkit avs devnet stop"
echo ""