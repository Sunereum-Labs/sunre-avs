#!/bin/bash

# Test script for Temperature Oracle AVS

echo "Temperature Oracle AVS Test Script"
echo "================================="

# Example task payload for New York with 25°C threshold
TASK_PAYLOAD=$(cat <<EOF
{
  "location": {
    "latitude": 40.7128,
    "longitude": -74.0060,
    "city": "New York",
    "country": "USA"
  },
  "threshold": 25.0
}
EOF
)

echo "Task Payload:"
echo "$TASK_PAYLOAD"
echo ""

# Encode the payload as base64 for the devkit call
ENCODED_PAYLOAD=$(echo -n "$TASK_PAYLOAD" | base64)

echo "To test the Temperature Oracle AVS with DevKit:"
echo ""
echo "1. Start the devnet:"
echo "   devkit avs devnet start"
echo ""
echo "2. Create a task (after devnet is running):"
echo "   devkit avs call --payload '$ENCODED_PAYLOAD'"
echo ""
echo "3. The response will contain:"
echo "   - Temperature: Current consensus temperature"
echo "   - MeetsThreshold: Whether temperature >= threshold"
echo "   - Confidence: Confidence level (0-1)"
echo "   - DataPoints: Number of data sources used"
echo ""
echo "Example for other cities:"
echo "- London: lat=51.5074, lon=-0.1278"
echo "- Tokyo: lat=35.6762, lon=139.6503"
echo "- Sydney: lat=-33.8688, lon=151.2093"