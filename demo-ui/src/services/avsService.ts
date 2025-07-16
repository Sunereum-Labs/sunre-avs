import { AVSTaskRequest, ClaimResult } from '../types/insurance';

const AVS_ENDPOINT = process.env.REACT_APP_AVS_ENDPOINT || 'http://localhost:8081';

export async function submitToAVS(request: AVSTaskRequest): Promise<ClaimResult> {
  try {
    // Encode the request payload as base64
    const payload = btoa(JSON.stringify(request));
    
    // Submit to AVS endpoint
    const response = await fetch(`${AVS_ENDPOINT}/task`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        payload: payload
      })
    });

    if (!response.ok) {
      throw new Error(`AVS request failed: ${response.statusText}`);
    }

    const result = await response.json();
    
    // Parse the AVS response
    if (result.Result) {
      const claimResult = JSON.parse(atob(result.Result));
      return {
        ...claimResult,
        timestamp: new Date(claimResult.timestamp)
      };
    }

    throw new Error('Invalid AVS response format');
  } catch (error) {
    console.error('AVS submission error:', error);
    
    // Return a mock response for demo purposes if AVS is not available
    return generateMockResponse(request);
  }
}

function generateMockResponse(request: AVSTaskRequest): ClaimResult {
  const policy = request.claimRequest.policy;
  const scenario = request.demoScenario || 'normal';
  
  let claimStatus: 'approved' | 'rejected' | 'partial' = 'rejected';
  let payoutAmount = 0;
  let triggeredPerils: any[] = [];

  // Simulate different outcomes based on scenario
  if (scenario === 'heat_wave' && policy.insuranceType === 'crop') {
    claimStatus = 'approved';
    payoutAmount = policy.coverageAmount * 0.5;
    triggeredPerils = [{
      peril: 'heat_wave',
      conditionsMet: true,
      payoutRatio: 0.5
    }];
  } else if (scenario === 'storm' && policy.insuranceType === 'event') {
    claimStatus = 'approved';
    payoutAmount = policy.coverageAmount;
    triggeredPerils = [{
      peril: 'excess_rain',
      conditionsMet: true,
      payoutRatio: 1.0
    }];
  } else if (scenario === 'cold_snap' && policy.insuranceType === 'travel') {
    claimStatus = 'approved';
    payoutAmount = policy.coverageAmount * 0.2;
    triggeredPerils = [{
      peril: 'cold_snap',
      conditionsMet: true,
      payoutRatio: 0.2
    }];
  }

  // Generate mock weather data
  const mockWeatherData = Array.from({ length: 15 }, (_, i) => ({
    source: ['OpenMeteo', 'WeatherAPI', 'VisualCrossing'][i % 3],
    temperature: scenario === 'heat_wave' ? 38 + Math.random() * 3 : 
                 scenario === 'cold_snap' ? -15 + Math.random() * 3 : 
                 22 + Math.random() * 3,
    timestamp: new Date(),
    confidence: 0.85 + Math.random() * 0.15
  }));

  return {
    claimId: `CLM-${Date.now().toString(36)}`,
    policyId: policy.policyId,
    claimStatus,
    triggeredPerils,
    payoutAmount,
    weatherData: mockWeatherData,
    verificationHash: `0x${Array.from({ length: 64 }, () => Math.floor(Math.random() * 16).toString(16)).join('')}`,
    timestamp: new Date(),
    processingTime: 2000 + Math.random() * 1000
  };
}