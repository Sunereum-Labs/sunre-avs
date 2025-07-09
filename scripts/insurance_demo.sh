#!/bin/bash

# Weather Insurance AVS Demo Script

echo "Weather Insurance AVS Demo"
echo "========================="
echo ""
echo "This demo showcases automated insurance claim processing using weather data consensus"
echo ""

# Demo 1: Crop Insurance - Heat Wave
echo "1. CROP INSURANCE - Heat Wave Protection"
echo "----------------------------------------"
CROP_HEAT_PAYLOAD=$(cat <<'EOF'
{
  "type": "insurance_claim",
  "claim_request": {
    "policy_id": "CROP-2024-001",
    "policy": {
      "policy_id": "CROP-2024-001",
      "policy_holder": "Green Valley Farms",
      "insurance_type": "crop",
      "location": {
        "latitude": 35.2271,
        "longitude": -80.8431,
        "city": "Charlotte",
        "country": "USA"
      },
      "coverage_amount": 100000,
      "premium": 5000,
      "start_date": "2024-06-01T00:00:00Z",
      "end_date": "2024-09-30T00:00:00Z",
      "triggers": [
        {
          "trigger_id": "HEAT-3DAY-35C",
          "peril": "heat_wave",
          "conditions": {
            "temperature_max": 35,
            "consecutive_days": 3,
            "time_window": {
              "start_month": 6,
              "end_month": 8
            }
          },
          "payout_ratio": 0.5,
          "description": "50% payout if temperature exceeds 35°C for 3 consecutive days"
        },
        {
          "trigger_id": "EXTREME-HEAT-40C",
          "peril": "heat_wave",
          "conditions": {
            "temperature_max": 40,
            "consecutive_days": 2
          },
          "payout_ratio": 1.0,
          "description": "Full payout for extreme heat (40°C+ for 2 days)"
        }
      ]
    },
    "claim_date": "2024-07-15T00:00:00Z",
    "automated_check": true
  },
  "demo_mode": true,
  "demo_scenario": "heat_wave"
}
EOF
)

echo "Payload: Crop insurance claim for heat wave damage"
echo "Expected: Approved claim with 50% payout ($50,000)"
echo ""
echo "Base64 encoded payload:"
echo "$CROP_HEAT_PAYLOAD" | base64
echo ""

# Demo 2: Event Insurance - Rain Cancellation
echo "2. EVENT INSURANCE - Weather Cancellation"
echo "-----------------------------------------"
EVENT_RAIN_PAYLOAD=$(cat <<'EOF'
{
  "type": "insurance_claim",
  "claim_request": {
    "policy_id": "EVENT-2024-MUSIC",
    "policy": {
      "policy_id": "EVENT-2024-MUSIC",
      "policy_holder": "Summer Music Festival LLC",
      "insurance_type": "event",
      "location": {
        "latitude": 40.7128,
        "longitude": -74.0060,
        "city": "New York",
        "country": "USA"
      },
      "coverage_amount": 500000,
      "premium": 25000,
      "start_date": "2024-08-15T00:00:00Z",
      "end_date": "2024-08-17T00:00:00Z",
      "triggers": [
        {
          "trigger_id": "RAIN-50MM",
          "peril": "excess_rain",
          "conditions": {
            "precipitation_min": 50,
            "time_window": {
              "start_hour": 8,
              "end_hour": 20
            }
          },
          "payout_ratio": 1.0,
          "description": "Full payout if rain exceeds 50mm during event hours"
        },
        {
          "trigger_id": "WIND-60KMH",
          "peril": "high_wind",
          "conditions": {
            "wind_speed_min": 60
          },
          "payout_ratio": 1.0,
          "description": "Full payout for dangerous wind conditions"
        }
      ]
    },
    "claim_date": "2024-08-16T00:00:00Z",
    "automated_check": true
  },
  "demo_mode": true,
  "demo_scenario": "normal"
}
EOF
)

echo "Payload: Event insurance claim for weather conditions"
echo "Expected: Rejected claim (normal weather, no triggers met)"
echo ""
echo "Base64 encoded payload:"
echo "$EVENT_RAIN_PAYLOAD" | base64
echo ""

# Demo 3: Travel Insurance - Cold Weather Delays
echo "3. TRAVEL INSURANCE - Flight Delay Protection"
echo "---------------------------------------------"
TRAVEL_COLD_PAYLOAD=$(cat <<'EOF'
{
  "type": "insurance_claim",
  "claim_request": {
    "policy_id": "TRAVEL-2024-0123",
    "policy": {
      "policy_id": "TRAVEL-2024-0123",
      "policy_holder": "John Doe",
      "insurance_type": "travel",
      "location": {
        "latitude": 41.9742,
        "longitude": -87.9073,
        "city": "Chicago O'Hare",
        "country": "USA"
      },
      "coverage_amount": 1000,
      "premium": 50,
      "start_date": "2024-12-20T00:00:00Z",
      "end_date": "2024-12-25T00:00:00Z",
      "triggers": [
        {
          "trigger_id": "COLD-DELAY",
          "peril": "cold_snap",
          "conditions": {
            "temperature_max": -10
          },
          "payout_ratio": 0.2,
          "description": "Daily compensation for extreme cold delays"
        }
      ]
    },
    "claim_date": "2024-12-22T00:00:00Z",
    "automated_check": true
  },
  "demo_mode": true,
  "demo_scenario": "cold_snap"
}
EOF
)

echo "Payload: Travel insurance for cold weather delays"
echo "Expected: Approved claim with 20% payout ($200 daily compensation)"
echo ""
echo "Base64 encoded payload:"
echo "$TRAVEL_COLD_PAYLOAD" | base64
echo ""

echo "HOW TO RUN THESE DEMOS:"
echo "======================="
echo ""
echo "1. Start the AVS devnet:"
echo "   devkit avs devnet start"
echo ""
echo "2. Submit a claim (copy the base64 payload above):"
echo "   devkit avs call --payload <BASE64_PAYLOAD>"
echo ""
echo "3. The AVS will:"
echo "   - Validate the insurance policy"
echo "   - Fetch weather data from multiple sources"
echo "   - Apply MAD consensus algorithm"
echo "   - Evaluate trigger conditions"
echo "   - Calculate payouts automatically"
echo "   - Return a verifiable claim decision"
echo ""
echo "Response includes:"
echo "- claim_id: Unique claim identifier"
echo "- claim_status: approved/rejected/partial/investigate"
echo "- payout_amount: Calculated compensation"
echo "- triggered_perils: Which conditions were met"
echo "- weather_data: Consensus weather assessment"
echo "- verification_hash: Cryptographic proof"