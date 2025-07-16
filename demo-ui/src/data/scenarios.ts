import { DemoScenario } from '../types/insurance';

export const demoScenarios: DemoScenario[] = [
  {
    id: 'crop-heat-wave',
    name: 'Crop Heat Protection',
    description: 'Protect crops from extreme heat damage with automatic payouts',
    policy: {
      policyId: 'CROP-2024-001',
      policyHolder: 'Green Valley Farms',
      insuranceType: 'crop',
      location: {
        latitude: 35.2271,
        longitude: -80.8431,
        city: 'Charlotte',
        country: 'USA'
      },
      coverageAmount: 100000,
      premium: 5000,
      startDate: new Date('2024-06-01'),
      endDate: new Date('2024-09-30'),
      triggers: [
        {
          triggerId: 'HEAT-3DAY-35C',
          peril: 'heat_wave',
          conditions: {
            temperatureMax: 35,
            consecutiveDays: 3,
            timeWindow: {
              startMonth: 6,
              endMonth: 8
            }
          },
          payoutRatio: 0.5,
          description: '50% payout if temperature exceeds 35°C for 3 consecutive days'
        },
        {
          triggerId: 'EXTREME-HEAT-40C',
          peril: 'heat_wave',
          conditions: {
            temperatureMax: 40,
            consecutiveDays: 2
          },
          payoutRatio: 1.0,
          description: 'Full payout for extreme heat (40°C+ for 2 days)'
        }
      ]
    },
    weatherPattern: 'heat_wave',
    expectedOutcome: {
      status: 'approved',
      payout: 50000,
      reason: 'Temperature exceeded 35°C for 5 consecutive days, triggering 50% payout'
    }
  },
  {
    id: 'event-rain',
    name: 'Event Cancellation',
    description: 'Insurance for outdoor events against weather disruptions',
    policy: {
      policyId: 'EVENT-2024-MUSIC',
      policyHolder: 'Summer Music Festival LLC',
      insuranceType: 'event',
      location: {
        latitude: 40.7128,
        longitude: -74.0060,
        city: 'New York',
        country: 'USA'
      },
      coverageAmount: 500000,
      premium: 25000,
      startDate: new Date('2024-08-15'),
      endDate: new Date('2024-08-17'),
      triggers: [
        {
          triggerId: 'RAIN-50MM',
          peril: 'excess_rain',
          conditions: {
            precipitationMin: 50,
            timeWindow: {
              startMonth: 8,
              endMonth: 8,
              startHour: 8,
              endHour: 20
            }
          },
          payoutRatio: 1.0,
          description: 'Full payout if rain exceeds 50mm during event hours'
        },
        {
          triggerId: 'WIND-60KMH',
          peril: 'high_wind',
          conditions: {
            windSpeedMin: 60
          },
          payoutRatio: 1.0,
          description: 'Full payout for dangerous wind conditions'
        }
      ]
    },
    weatherPattern: 'storm',
    expectedOutcome: {
      status: 'approved',
      payout: 500000,
      reason: 'Severe storm with 75mm rain and 80km/h winds forced event cancellation'
    }
  },
  {
    id: 'travel-cold',
    name: 'Travel Delay Protection',
    description: 'Compensation for weather-related travel delays',
    policy: {
      policyId: 'TRAVEL-2024-0123',
      policyHolder: 'John Doe',
      insuranceType: 'travel',
      location: {
        latitude: 41.9742,
        longitude: -87.9073,
        city: 'Chicago O\'Hare',
        country: 'USA'
      },
      coverageAmount: 1000,
      premium: 50,
      startDate: new Date('2024-12-20'),
      endDate: new Date('2024-12-25'),
      triggers: [
        {
          triggerId: 'COLD-DELAY',
          peril: 'cold_snap',
          conditions: {
            temperatureMax: -10
          },
          payoutRatio: 0.2,
          description: 'Daily compensation (20%) for extreme cold delays'
        }
      ]
    },
    weatherPattern: 'cold_snap',
    expectedOutcome: {
      status: 'approved',
      payout: 200,
      reason: 'Temperature dropped to -15°C causing flight delays, daily compensation triggered'
    }
  },
  {
    id: 'crop-normal',
    name: 'No Claim Scenario',
    description: 'Example of normal weather conditions with no payout',
    policy: {
      policyId: 'CROP-2024-002',
      policyHolder: 'Sunny Acres Farm',
      insuranceType: 'crop',
      location: {
        latitude: 37.7749,
        longitude: -122.4194,
        city: 'San Francisco',
        country: 'USA'
      },
      coverageAmount: 75000,
      premium: 3500,
      startDate: new Date('2024-05-01'),
      endDate: new Date('2024-10-31'),
      triggers: [
        {
          triggerId: 'FROST-0C',
          peril: 'frost',
          conditions: {
            temperatureMax: 0,
            consecutiveDays: 2
          },
          payoutRatio: 0.75,
          description: '75% payout for frost damage (below 0°C for 2 days)'
        }
      ]
    },
    weatherPattern: 'normal',
    expectedOutcome: {
      status: 'rejected',
      payout: 0,
      reason: 'Weather conditions remained within normal range, no triggers activated'
    }
  },
  {
    id: 'energy-wind',
    name: 'Renewable Energy',
    description: 'Protection for wind farm revenue during low wind periods',
    policy: {
      policyId: 'ENERGY-2024-WIND',
      policyHolder: 'GreenPower Wind Farm',
      insuranceType: 'energy',
      location: {
        latitude: 52.3676,
        longitude: 4.9041,
        city: 'North Sea',
        country: 'Netherlands'
      },
      coverageAmount: 250000,
      premium: 15000,
      startDate: new Date('2024-01-01'),
      endDate: new Date('2024-12-31'),
      triggers: [
        {
          triggerId: 'LOW-WIND-5MS',
          peril: 'low_wind',
          conditions: {
            windSpeedMax: 5,
            consecutiveDays: 5
          },
          payoutRatio: 0.3,
          description: '30% payout for low wind speed (< 5 m/s) for 5+ days'
        }
      ]
    },
    weatherPattern: 'normal',
    expectedOutcome: {
      status: 'partial',
      payout: 75000,
      reason: 'Wind speeds below threshold for 6 days, partial compensation activated'
    }
  }
];