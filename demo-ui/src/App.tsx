import React, { useState, useEffect } from 'react';
import { 
  CloudRain, 
  Activity, 
  Shield, 
  Zap, 
  TrendingUp,
  Globe,
  CheckCircle,
  AlertCircle,
  Loader2,
  BarChart3,
  Thermometer,
  Wind,
  Droplets,
  Sun,
  MapPin
} from 'lucide-react';
import './App.css';

// Components
import AVSStatus from './components/AVSStatus';
import LiveWeatherCard from './components/LiveWeatherCard';
import InsuranceScenarioCard from './components/InsuranceScenarioCard';
import ClaimProcessingFlow from './components/ClaimProcessingFlow';
import ConsensusVisualization from './components/ConsensusVisualization';

// Types
import { AVSConnection } from './types/insurance';

// Data
import { demoScenarios } from './data/scenarios';

interface LiveWeatherData {
  temperature: number;
  trend: 'warming' | 'cooling' | 'stable';
  sources: Array<{
    name: string;
    temperature: number;
    confidence: number;
  }>;
  consensus: {
    algorithm: string;
    confidence: number;
    totalSources: number;
  };
  historicalData: {
    temperature4hAgo: number;
    change: number;
  };
  recommendation: {
    type: string;
    scenario: string;
    description: string;
    coverage: string;
    premium: string;
  };
}

function App() {
  const [activeTab, setActiveTab] = useState<'overview' | 'scenarios' | 'live'>('overview');
  const [selectedScenario, setSelectedScenario] = useState<any>(null);
  const [isProcessing, setIsProcessing] = useState(false);
  const [liveWeatherData, setLiveWeatherData] = useState<LiveWeatherData | null>(null);
  const [isCheckingWeather, setIsCheckingWeather] = useState(false);
  const [consensusDemo, setConsensusDemo] = useState<any>(null);
  const [avsConnection, setAvsConnection] = useState<AVSConnection>({
    endpoint: 'http://localhost:8081',
    status: 'disconnected',
    networkType: 'devnet'
  });

  // Check AVS connection
  useEffect(() => {
    const checkConnection = async () => {
      try {
        const response = await fetch(`${avsConnection.endpoint}/health`);
        if (response.ok) {
          setAvsConnection(prev => ({ ...prev, status: 'connected' }));
        }
      } catch (error) {
        setAvsConnection(prev => ({ ...prev, status: 'disconnected' }));
      }
    };

    checkConnection();
    const interval = setInterval(checkConnection, 5000);
    return () => clearInterval(interval);
  }, [avsConnection.endpoint]);

  // Trigger consensus check (simulating insurance policy check)
  const triggerConsensusCheck = async () => {
    setIsCheckingWeather(true);
    setConsensusDemo({ status: 'processing', phase: 'fetching' });
    
    try {
      // Simulate progressive data gathering
      const simulateProgress = async () => {
        // Phase 1: Fetching from multiple sources
        await new Promise(resolve => setTimeout(resolve, 1000));
        setConsensusDemo({
          status: 'processing',
          phase: 'fetching',
          message: 'Fetching weather data from multiple sources...',
          sources_checked: ['Tomorrow.io', 'WeatherAPI', 'Open-Meteo']
        });

        // Phase 2: Gathering data
        await new Promise(resolve => setTimeout(resolve, 1500));
        const weatherData = liveWeatherData || {
          temperature: 36.5,
          sources: [
            { name: 'Tomorrow.io', temperature: 36.8, confidence: 0.95 },
            { name: 'WeatherAPI', temperature: 36.2, confidence: 0.92 },
            { name: 'Open-Meteo', temperature: 36.5, confidence: 0.89 }
          ]
        };
        
        setConsensusDemo({
          status: 'processing',
          phase: 'consensus',
          message: 'Applying consensus algorithm...',
          sources_data: weatherData.sources,
          consensus_temperature: weatherData.temperature
        });

        // Phase 3: Policy evaluation
        await new Promise(resolve => setTimeout(resolve, 1000));
        const policyTriggerTemp = 35;
        const meetsConditions = weatherData.temperature > policyTriggerTemp;
        
        setConsensusDemo({
          status: 'processing',
          phase: 'evaluating',
          message: 'Evaluating insurance policy conditions...',
          policy_conditions: {
            trigger_temperature: policyTriggerTemp,
            actual_temperature: weatherData.temperature,
            condition_met: meetsConditions
          }
        });

        // Phase 4: Final result
        await new Promise(resolve => setTimeout(resolve, 1000));
        setConsensusDemo({
          status: 'completed',
          claim_status: meetsConditions ? 'approved' : 'rejected',
          claim_id: `CLM-${Date.now()}`,
          payout_amount: meetsConditions ? 50000 : 0,
          weather_data: {
            consensus_temperature: weatherData.temperature,
            data_sources: weatherData.sources.length,
            consensus_method: 'MAD (Median Absolute Deviation)',
            sources_detail: weatherData.sources
          },
          policy_evaluation: {
            policy_id: 'NYC-DEMO-001',
            trigger_type: 'Heat Wave',
            trigger_condition: `Temperature > ${policyTriggerTemp}°C`,
            actual_value: `${weatherData.temperature}°C`,
            condition_met: meetsConditions,
            payout_ratio: 0.5,
            coverage_amount: 100000
          }
        });
      };

      // First try actual API
      const payload = {
        type: 'insurance_claim',
        claim_request: {
          policy_id: 'NYC-DEMO-001',
          policy: {
            policy_id: 'NYC-DEMO-001',
            policy_holder: 'Demo User',
            insurance_type: 'event',
            location: {
              latitude: 40.7128,
              longitude: -74.0060,
              city: 'New York',
              country: 'USA'
            },
            coverage_amount: 100000,
            premium: 2000,
            triggers: [{
              trigger_id: 'HEAT-NYC-001',
              peril: 'heat_wave',
              conditions: {
                temperature_max: 35,
                consecutive_days: 1
              },
              payout_ratio: 0.5,
              description: 'Heat wave event protection'
            }]
          },
          claim_date: new Date().toISOString(),
          automated_check: true
        },
        demo_mode: true
      };

      const response = await fetch(`${avsConnection.endpoint}/task`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ payload: btoa(JSON.stringify(payload)) }),
      });

      if (response.ok) {
        const data = await response.json();
        const result = JSON.parse(atob(data.Result || data.result));
        
        // If backend fails, run simulation
        if (result.claim_status === 'rejected' && result.weather_data?.data_sources === 0) {
          await simulateProgress();
        } else {
          setConsensusDemo(result);
        }
      } else {
        // Run simulation on non-200 response
        await simulateProgress();
      }
    } catch (error) {
      console.error('Consensus check failed:', error);
      setConsensusDemo({ status: 'error', message: 'Connection failed' });
    } finally {
      setIsCheckingWeather(false);
    }
  };

  // Fetch live weather data
  const fetchLiveWeather = async () => {
    try {
      const payload = {
        type: 'live_weather_demo',
        location: {
          latitude: 40.7128,
          longitude: -74.0060,
          city: 'New York',
          country: 'USA'
        }
      };

      const response = await fetch(`${avsConnection.endpoint}/task`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ payload: btoa(JSON.stringify(payload)) }),
      });

      if (response.ok) {
        const data = await response.json();
        const result = JSON.parse(atob(data.Result || data.result));
        
        // Transform the data to our format
        setLiveWeatherData({
          temperature: result.current_weather.temperature,
          trend: result.historical_weather.trend,
          sources: Object.entries(result.current_weather.data_sources).map(([name, data]: [string, any]) => ({
            name,
            temperature: data.current_temperature,
            confidence: data.confidence
          })),
          consensus: {
            algorithm: result.consensus_details.algorithm,
            confidence: result.current_weather.consensus_confidence,
            totalSources: result.consensus_details.total_sources
          },
          historicalData: {
            temperature4hAgo: result.historical_weather.temperature_4h_ago,
            change: result.historical_weather.temperature_change
          },
          recommendation: result.insurance_recommendation
        });
      }
    } catch (error) {
      console.error('Failed to fetch live weather:', error);
    }
  };

  useEffect(() => {
    if (activeTab === 'live' || activeTab === 'overview') {
      fetchLiveWeather();
      const interval = setInterval(fetchLiveWeather, 30000);
      return () => clearInterval(interval);
    }
  }, [activeTab]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-dark-400 to-dark-500 text-white">
      {/* Header */}
      <header className="border-b border-gray-800 backdrop-blur-md bg-dark-500/50 sticky top-0 z-50">
        <div className="container mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="p-2 bg-gradient-to-br from-primary-500 to-primary-600 rounded-lg">
                <CloudRain className="w-6 h-6" />
              </div>
              <div>
                <h1 className="text-2xl font-bold">SunRe AVS</h1>
                <p className="text-sm text-gray-400">Decentralized Weather Insurance</p>
              </div>
            </div>
            <AVSStatus connection={avsConnection} />
          </div>
        </div>
      </header>

      {/* Navigation Tabs */}
      <div className="container mx-auto px-6 pt-8">
        <div className="flex space-x-1 bg-dark-400/50 p-1 rounded-lg backdrop-blur-sm max-w-md">
          <button
            onClick={() => setActiveTab('overview')}
            className={`flex-1 px-4 py-2 rounded-md font-medium transition-all ${
              activeTab === 'overview'
                ? 'bg-primary-600 text-white shadow-lg'
                : 'text-gray-400 hover:text-white hover:bg-dark-300/50'
            }`}
          >
            Overview
          </button>
          <button
            onClick={() => setActiveTab('scenarios')}
            className={`flex-1 px-4 py-2 rounded-md font-medium transition-all ${
              activeTab === 'scenarios'
                ? 'bg-primary-600 text-white shadow-lg'
                : 'text-gray-400 hover:text-white hover:bg-dark-300/50'
            }`}
          >
            Demo Scenarios
          </button>
          <button
            onClick={() => setActiveTab('live')}
            className={`flex-1 px-4 py-2 rounded-md font-medium transition-all relative ${
              activeTab === 'live'
                ? 'bg-primary-600 text-white shadow-lg'
                : 'text-gray-400 hover:text-white hover:bg-dark-300/50'
            }`}
          >
            Live NYC
            <span className="absolute -top-1 -right-1 w-2 h-2 bg-green-500 rounded-full animate-pulse" />
          </button>
        </div>
      </div>

      {/* Main Content */}
      <main className="container mx-auto px-6 py-8">
        {activeTab === 'overview' && (
          <div className="space-y-8">
            {/* Hero Section */}
            <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-primary-600 to-primary-800 p-8 md:p-12">
              <div className="relative z-10">
                <h2 className="text-4xl font-bold mb-4">
                  Automated Weather Insurance on EigenLayer
                </h2>
                <p className="text-xl text-primary-100 max-w-2xl mb-8">
                  Experience instant, trust-minimized insurance claims powered by decentralized weather oracles 
                  and cryptographic consensus.
                </p>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                  <div className="flex items-center space-x-3">
                    <div className="p-3 bg-white/20 rounded-lg backdrop-blur-sm">
                      <Zap className="w-6 h-6" />
                    </div>
                    <div>
                      <p className="font-semibold">2-3 min</p>
                      <p className="text-sm text-primary-200">Claim Processing</p>
                    </div>
                  </div>
                  <div className="flex items-center space-x-3">
                    <div className="p-3 bg-white/20 rounded-lg backdrop-blur-sm">
                      <Shield className="w-6 h-6" />
                    </div>
                    <div>
                      <p className="font-semibold">Zero Fraud</p>
                      <p className="text-sm text-primary-200">Consensus Validation</p>
                    </div>
                  </div>
                  <div className="flex items-center space-x-3">
                    <div className="p-3 bg-white/20 rounded-lg backdrop-blur-sm">
                      <Activity className="w-6 h-6" />
                    </div>
                    <div>
                      <p className="font-semibold">3+ Sources</p>
                      <p className="text-sm text-primary-200">Weather APIs</p>
                    </div>
                  </div>
                </div>
              </div>
              <div className="absolute top-0 right-0 -mt-20 -mr-20 w-80 h-80 bg-white/5 rounded-full blur-3xl" />
              <div className="absolute bottom-0 left-0 -mb-20 -ml-20 w-80 h-80 bg-white/5 rounded-full blur-3xl" />
            </div>

            {/* Live Weather Preview */}
            {liveWeatherData && (
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <LiveWeatherCard data={liveWeatherData} />
                <ConsensusVisualization data={liveWeatherData} />
              </div>
            )}

            {/* How It Works */}
            <div className="bg-dark-400/50 backdrop-blur-sm rounded-2xl p-8">
              <h3 className="text-2xl font-bold mb-6 flex items-center">
                <BarChart3 className="w-6 h-6 mr-2 text-primary-500" />
                How It Works
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                {[
                  { icon: Globe, title: 'Multi-Source Data', desc: 'Fetch weather from 3+ APIs' },
                  { icon: Activity, title: 'MAD Consensus', desc: 'Statistical validation algorithm' },
                  { icon: Shield, title: 'Smart Contract', desc: 'Automated trigger evaluation' },
                  { icon: Zap, title: 'Instant Payout', desc: 'Cryptographic settlement' }
                ].map((step, idx) => (
                  <div key={idx} className="text-center">
                    <div className="mx-auto w-16 h-16 bg-primary-600/20 rounded-full flex items-center justify-center mb-4">
                      <step.icon className="w-8 h-8 text-primary-500" />
                    </div>
                    <h4 className="font-semibold mb-2">{step.title}</h4>
                    <p className="text-sm text-gray-400">{step.desc}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {activeTab === 'scenarios' && (
          <div className="space-y-6">
            <div className="text-center mb-8">
              <h2 className="text-3xl font-bold mb-4">Interactive Demo Scenarios</h2>
              <p className="text-gray-400 max-w-2xl mx-auto">
                Select a scenario to see how parametric insurance automatically processes claims 
                based on real weather conditions.
              </p>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {demoScenarios.map((scenario) => (
                <InsuranceScenarioCard
                  key={scenario.id}
                  scenario={scenario}
                  isSelected={selectedScenario?.id === scenario.id}
                  onSelect={() => setSelectedScenario(scenario)}
                />
              ))}
            </div>

            {selectedScenario && (
              <ClaimProcessingFlow
                scenario={selectedScenario}
                isProcessing={isProcessing}
                onProcess={() => {
                  setIsProcessing(true);
                  setTimeout(() => setIsProcessing(false), 3000);
                }}
              />
            )}
          </div>
        )}

        {activeTab === 'live' && (
          <div className="space-y-6">
            <div className="text-center mb-8">
              <h2 className="text-3xl font-bold mb-4">Live Weather Data - New York City</h2>
              <p className="text-gray-400 max-w-2xl mx-auto">
                Real-time weather consensus from multiple sources, updated every 30 seconds.
              </p>
            </div>

            {liveWeatherData ? (
              <>
                {/* Demo Action Button */}
                <div className="bg-gradient-to-r from-primary-600/20 to-purple-600/20 backdrop-blur-sm rounded-2xl p-6 border border-primary-600/30">
                  <div className="flex flex-col md:flex-row items-center justify-between gap-4">
                    <div className="text-center md:text-left">
                      <h3 className="text-xl font-semibold mb-2">Simulate Insurance Policy Check</h3>
                      <p className="text-gray-400">
                        Click to see how the AVS processes weather data consensus when an insurance policy needs verification
                      </p>
                    </div>
                    <button
                      onClick={triggerConsensusCheck}
                      disabled={isCheckingWeather}
                      className="px-6 py-3 bg-gradient-to-r from-primary-600 to-purple-600 hover:from-primary-700 hover:to-purple-700 rounded-lg font-semibold transition-all flex items-center space-x-2 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {isCheckingWeather ? (
                        <>
                          <Loader2 className="w-5 h-5 animate-spin" />
                          <span>Processing...</span>
                        </>
                      ) : (
                        <>
                          <Zap className="w-5 h-5" />
                          <span>Trigger Consensus Check</span>
                        </>
                      )}
                    </button>
                  </div>
                  
                  {/* Consensus Demo Results */}
                  {consensusDemo && (
                    <div className="mt-6 p-4 bg-dark-500/50 rounded-lg">
                      {consensusDemo.status === 'processing' ? (
                        <div className="space-y-4">
                          <div className="text-center">
                            <Loader2 className="w-8 h-8 animate-spin mx-auto mb-2 text-primary-500" />
                            <p className="text-gray-300 font-medium">{consensusDemo.message || 'Processing claim...'}</p>
                          </div>
                          
                          {/* Phase 1: Fetching */}
                          {consensusDemo.phase === 'fetching' && consensusDemo.sources_checked && (
                            <div className="space-y-2">
                              <p className="text-sm text-gray-400">Checking weather sources:</p>
                              {consensusDemo.sources_checked.map((source: string, idx: number) => (
                                <div key={source} className="flex items-center space-x-2 text-sm">
                                  <CheckCircle className="w-4 h-4 text-green-500" />
                                  <span className="text-gray-300">{source}</span>
                                </div>
                              ))}
                            </div>
                          )}
                          
                          {/* Phase 2: Consensus */}
                          {consensusDemo.phase === 'consensus' && consensusDemo.sources_data && (
                            <div className="space-y-3">
                              <p className="text-sm text-gray-400">Weather data collected:</p>
                              <div className="grid grid-cols-3 gap-2">
                                {consensusDemo.sources_data.map((source: any) => (
                                  <div key={source.name} className="bg-dark-600/50 rounded p-2 text-center">
                                    <p className="text-xs text-gray-400">{source.name}</p>
                                    <p className="font-semibold">{source.temperature.toFixed(1)}°C</p>
                                    <p className="text-xs text-gray-500">{(source.confidence * 100).toFixed(0)}% conf</p>
                                  </div>
                                ))}
                              </div>
                              <div className="bg-primary-900/20 border border-primary-700/50 rounded p-2">
                                <p className="text-sm text-primary-300">
                                  Consensus Temperature: <span className="font-semibold">{consensusDemo.consensus_temperature?.toFixed(1)}°C</span>
                                </p>
                              </div>
                            </div>
                          )}
                          
                          {/* Phase 3: Evaluating */}
                          {consensusDemo.phase === 'evaluating' && consensusDemo.policy_conditions && (
                            <div className="space-y-2">
                              <p className="text-sm text-gray-400">Evaluating policy conditions:</p>
                              <div className="bg-dark-600/50 rounded p-3 space-y-2">
                                <div className="flex justify-between text-sm">
                                  <span className="text-gray-400">Trigger condition:</span>
                                  <span className="text-gray-300">Temperature &gt; {consensusDemo.policy_conditions.trigger_temperature}°C</span>
                                </div>
                                <div className="flex justify-between text-sm">
                                  <span className="text-gray-400">Actual temperature:</span>
                                  <span className="font-semibold text-white">{consensusDemo.policy_conditions.actual_temperature}°C</span>
                                </div>
                                <div className="flex justify-between text-sm">
                                  <span className="text-gray-400">Condition met:</span>
                                  <span className={`font-semibold ${consensusDemo.policy_conditions.condition_met ? 'text-green-400' : 'text-red-400'}`}>
                                    {consensusDemo.policy_conditions.condition_met ? 'YES' : 'NO'}
                                  </span>
                                </div>
                              </div>
                            </div>
                          )}
                        </div>
                      ) : consensusDemo.status === 'error' ? (
                        <div className="text-center py-4 text-red-400">
                          <AlertCircle className="w-8 h-8 mx-auto mb-2" />
                          <p>{consensusDemo.message}</p>
                        </div>
                      ) : consensusDemo.status === 'completed' ? (
                        <div className="space-y-4">
                          {/* Final Result */}
                          <div className="bg-gradient-to-r from-primary-900/20 to-purple-900/20 backdrop-blur-sm rounded p-4 border border-primary-700/50">
                            <div className="flex items-center justify-between mb-3">
                              <span className="text-lg font-semibold">Claim Decision</span>
                              <span className={`px-3 py-1 rounded-full text-sm font-semibold ${
                                consensusDemo.claim_status === 'approved' 
                                  ? 'bg-green-900/50 text-green-300 border border-green-700/50' 
                                  : 'bg-red-900/50 text-red-300 border border-red-700/50'
                              }`}>
                                {consensusDemo.claim_status?.toUpperCase()}
                              </span>
                            </div>
                            
                            {consensusDemo.payout_amount > 0 && (
                              <div className="flex items-center justify-between text-lg">
                                <span className="text-gray-300">Payout Amount:</span>
                                <span className="font-bold text-green-400">${consensusDemo.payout_amount.toLocaleString()}</span>
                              </div>
                            )}
                          </div>
                          
                          {/* Weather Data Summary */}
                          {consensusDemo.weather_data && (
                            <div className="bg-dark-600/50 rounded p-4 space-y-3">
                              <h4 className="text-sm font-semibold text-gray-300">Weather Consensus Summary</h4>
                              <div className="grid grid-cols-2 gap-4 text-sm">
                                <div>
                                  <span className="text-gray-400">Consensus Temp:</span>
                                  <p className="font-semibold">{consensusDemo.weather_data.consensus_temperature?.toFixed(1)}°C</p>
                                </div>
                                <div>
                                  <span className="text-gray-400">Data Sources:</span>
                                  <p className="font-semibold">{consensusDemo.weather_data.data_sources}</p>
                                </div>
                                <div className="col-span-2">
                                  <span className="text-gray-400">Algorithm:</span>
                                  <p className="font-semibold">{consensusDemo.weather_data.consensus_method}</p>
                                </div>
                              </div>
                              
                              {consensusDemo.weather_data.sources_detail && (
                                <div className="pt-3 border-t border-gray-700">
                                  <p className="text-xs text-gray-400 mb-2">Individual source readings:</p>
                                  <div className="space-y-1">
                                    {consensusDemo.weather_data.sources_detail.map((source: any) => (
                                      <div key={source.name} className="flex justify-between text-xs">
                                        <span className="text-gray-500">{source.name}:</span>
                                        <span className="text-gray-400">{source.temperature.toFixed(1)}°C ({(source.confidence * 100).toFixed(0)}%)</span>
                                      </div>
                                    ))}
                                  </div>
                                </div>
                              )}
                            </div>
                          )}
                          
                          {/* Policy Evaluation Details */}
                          {consensusDemo.policy_evaluation && (
                            <div className="bg-dark-600/50 rounded p-4 space-y-2 text-sm">
                              <h4 className="font-semibold text-gray-300">Policy Evaluation</h4>
                              <div className="space-y-1 text-gray-400">
                                <div className="flex justify-between">
                                  <span>Policy ID:</span>
                                  <span className="font-mono text-xs">{consensusDemo.policy_evaluation.policy_id}</span>
                                </div>
                                <div className="flex justify-between">
                                  <span>Trigger Type:</span>
                                  <span>{consensusDemo.policy_evaluation.trigger_type}</span>
                                </div>
                                <div className="flex justify-between">
                                  <span>Condition:</span>
                                  <span>{consensusDemo.policy_evaluation.trigger_condition}</span>
                                </div>
                                <div className="flex justify-between">
                                  <span>Actual Value:</span>
                                  <span className="font-semibold text-white">{consensusDemo.policy_evaluation.actual_value}</span>
                                </div>
                                <div className="flex justify-between">
                                  <span>Coverage:</span>
                                  <span>${consensusDemo.policy_evaluation.coverage_amount?.toLocaleString()}</span>
                                </div>
                                <div className="flex justify-between">
                                  <span>Payout Ratio:</span>
                                  <span>{(consensusDemo.policy_evaluation.payout_ratio * 100).toFixed(0)}%</span>
                                </div>
                              </div>
                            </div>
                          )}
                          
                          {consensusDemo.claim_id && (
                            <div className="pt-2 border-t border-gray-700">
                              <p className="text-xs text-gray-400">Claim ID: <span className="font-mono">{consensusDemo.claim_id}</span></p>
                            </div>
                          )}
                        </div>
                      ) : (
                        <div className="space-y-4">
                          <div className="flex items-center justify-between">
                            <span className="text-gray-400">Claim Status:</span>
                            <span className={`font-semibold ${consensusDemo.claim_status === 'approved' ? 'text-green-400' : 'text-red-400'}`}>
                              {consensusDemo.claim_status?.toUpperCase() || 'PENDING'}
                            </span>
                          </div>
                        </div>
                      )}
                    </div>
                  )}
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                  {/* Current Temperature */}
                  <div className="lg:col-span-2">
                    <LiveWeatherCard data={liveWeatherData} detailed />
                  </div>

                  {/* Insurance Recommendation */}
                  <div className="bg-gradient-to-br from-purple-900/20 to-pink-900/20 backdrop-blur-sm rounded-2xl p-6 border border-purple-800/50">
                    <h3 className="text-lg font-semibold mb-4 flex items-center">
                      <Shield className="w-5 h-5 mr-2 text-purple-400" />
                      Dynamic Insurance Match
                    </h3>
                    <div className="space-y-3">
                      <div className="bg-dark-500/50 rounded-lg p-4">
                        <p className="text-sm text-gray-400 mb-1">Recommended Product</p>
                        <p className="font-semibold">{liveWeatherData.recommendation.scenario}</p>
                      </div>
                      <div className="bg-dark-500/50 rounded-lg p-4">
                        <p className="text-sm text-gray-400 mb-1">Coverage</p>
                        <p className="font-semibold text-green-400">{liveWeatherData.recommendation.coverage}</p>
                      </div>
                      <div className="bg-dark-500/50 rounded-lg p-4">
                        <p className="text-sm text-gray-400 mb-1">Premium</p>
                        <p className="font-semibold">{liveWeatherData.recommendation.premium}</p>
                      </div>
                      <p className="text-sm text-gray-300 mt-4">
                        {liveWeatherData.recommendation.description}
                      </p>
                    </div>
                  </div>
                </div>

                {/* Consensus Details */}
                <ConsensusVisualization data={liveWeatherData} detailed />

                {/* Data Sources */}
                <div className="bg-dark-400/50 backdrop-blur-sm rounded-2xl p-6">
                  <h3 className="text-xl font-semibold mb-4">Weather Data Sources</h3>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {liveWeatherData.sources.map((source) => (
                      <div key={source.name} className="bg-dark-500/50 rounded-lg p-4 border border-gray-800">
                        <div className="flex items-center justify-between mb-2">
                          <span className="font-medium">{source.name}</span>
                          <CheckCircle className="w-5 h-5 text-green-500" />
                        </div>
                        <div className="space-y-1">
                          <div className="flex justify-between">
                            <span className="text-sm text-gray-400">Temperature</span>
                            <span className="font-semibold">{source.temperature.toFixed(1)}°C</span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-sm text-gray-400">Confidence</span>
                            <span className="font-semibold">{(source.confidence * 100).toFixed(0)}%</span>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </>
            ) : (
              <div className="flex items-center justify-center h-64">
                <Loader2 className="w-8 h-8 animate-spin text-primary-500" />
              </div>
            )}
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="border-t border-gray-800 mt-16">
        <div className="container mx-auto px-6 py-8">
          <div className="flex flex-col md:flex-row items-center justify-between">
            <div className="text-gray-400 text-sm">
              Powered by EigenLayer AVS | Built with Hourglass Framework
            </div>
            <div className="flex items-center space-x-6 mt-4 md:mt-0">
              <a href="#" className="text-gray-400 hover:text-white transition-colors">Documentation</a>
              <a href="#" className="text-gray-400 hover:text-white transition-colors">GitHub</a>
              <a href="#" className="text-gray-400 hover:text-white transition-colors">API</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}

export default App;