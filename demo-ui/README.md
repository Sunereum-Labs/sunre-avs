# SunRe AVS Interactive Demo

A compelling TypeScript/React interface that showcases the power of automated weather insurance claims processing using EigenLayer AVS.

## Features

### 🎯 Interactive Demo Scenarios
- **Crop Insurance**: Heat wave protection for farmers
- **Event Insurance**: Weather cancellation coverage
- **Travel Insurance**: Flight delay compensation
- **Energy Insurance**: Renewable energy revenue protection
- **Normal Weather**: No-claim scenario demonstration

### 🌡️ Live Weather Simulation
- Multi-source weather data visualization
- Real-time temperature, humidity, wind, and precipitation metrics
- Visual alerts for extreme weather conditions
- 10-day weather pattern simulation

### ⚡ Claim Processing Visualization
- Step-by-step processing animation
- Real-time consensus mechanism demonstration
- Instant payout calculation
- Cryptographic verification display

### 📊 Key Metrics Display
- Processing time comparison (seconds vs weeks)
- Cost savings visualization (99% reduction)
- Confidence scores and data source tracking
- Fraud prevention demonstration

## Getting Started

### Prerequisites
- Node.js 16+ installed
- SunRe AVS running locally (optional for full integration)

### Installation

```bash
# Navigate to demo directory
cd demo-ui

# Install dependencies
npm install

# Start development server
npm start
```

The demo will open at `http://localhost:3000`

### Running with AVS Integration

1. Start the SunRe AVS:
```bash
# In the main project directory
devkit avs devnet start
```

2. Start the demo UI:
```bash
# In the demo-ui directory
npm start
```

The UI will automatically detect and connect to the AVS at `localhost:8080`.

## Demo Flow

### 1. Select a Scenario
Choose from 5 pre-configured insurance scenarios, each demonstrating different weather perils and outcomes.

### 2. Review Policy Details
Examine the insurance policy including:
- Coverage amount and premium
- Trigger conditions
- Payout ratios
- Policy period

### 3. Simulate Weather
Watch as the system:
- Fetches data from multiple sources
- Applies consensus algorithm
- Highlights extreme weather events

### 4. Process Claim
Experience the automated claim process:
- Policy validation
- Weather data aggregation
- Consensus application
- Trigger evaluation
- Cryptographic signing
- Instant payout decision

## Technology Stack

- **React 18**: Modern UI framework
- **TypeScript**: Type-safe development
- **CSS3**: Responsive, animated interface
- **EigenLayer AVS**: Backend integration

## Key Components

### `App.tsx`
Main application orchestrating the demo flow

### `PolicyBuilder.tsx`
Displays insurance policy details and triggers

### `WeatherSimulator.tsx`
Simulates multi-source weather data collection

### `ClaimProcessor.tsx`
Visualizes the claim processing pipeline

### `DemoScenarioSelector.tsx`
Interactive scenario selection interface

## Customization

### Adding New Scenarios
Edit `src/data/scenarios.ts` to add custom insurance products:

```typescript
{
  id: 'custom-scenario',
  name: 'Custom Insurance',
  policy: { /* policy details */ },
  weatherPattern: 'heat_wave',
  expectedOutcome: { /* expected result */ }
}
```

### Styling
Modify `src/App.css` for custom theming and animations.

## Benefits Demonstrated

1. **Speed**: 2-3 seconds vs 30-90 days
2. **Cost**: $5-10 vs $500-2000 per claim
3. **Trust**: Cryptographic consensus vs corporate promises
4. **Transparency**: All data verifiable on-chain
5. **Automation**: Zero paperwork required

## Production Deployment

```bash
# Build optimized production bundle
npm run build

# Deploy to static hosting
# The build folder contains all static assets
```

## Environment Variables

```bash
# AVS endpoint (optional)
REACT_APP_AVS_ENDPOINT=http://localhost:8080
```

## Troubleshooting

### AVS Connection Issues
- Ensure AVS is running on port 8080
- Check browser console for connection errors
- Demo works offline with simulated responses

### Performance
- Chrome/Edge recommended for best performance
- Animations may vary on different devices

---

This demo showcases how blockchain technology can revolutionize insurance, making it instant, transparent, and accessible to everyone.