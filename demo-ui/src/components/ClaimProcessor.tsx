import React, { useState, useEffect } from 'react';
import { InsurancePolicy, ClaimResult, WeatherDataPoint } from '../types/insurance';

interface ClaimProcessorProps {
  policy: InsurancePolicy | null;
  weatherData: WeatherDataPoint[];
  onProcess: () => void;
  isProcessing: boolean;
  result: ClaimResult | null;
}

interface ProcessStep {
  id: string;
  name: string;
  icon: string;
  duration: number;
}

const ClaimProcessor: React.FC<ClaimProcessorProps> = ({ 
  policy, 
  weatherData, 
  onProcess, 
  isProcessing, 
  result 
}) => {
  const [currentStep, setCurrentStep] = useState<number>(-1);

  const steps: ProcessStep[] = [
    { id: 'validate', name: 'Validate Policy', icon: '📋', duration: 500 },
    { id: 'fetch', name: 'Aggregate Weather Data', icon: '🌡️', duration: 1000 },
    { id: 'consensus', name: 'Apply Consensus', icon: '🤝', duration: 800 },
    { id: 'evaluate', name: 'Evaluate Triggers', icon: '⚡', duration: 600 },
    { id: 'sign', name: 'Sign & Verify', icon: '🔐', duration: 400 },
    { id: 'complete', name: 'Complete', icon: '✅', duration: 200 }
  ];

  useEffect(() => {
    if (isProcessing && currentStep < steps.length - 1) {
      const timer = setTimeout(() => {
        setCurrentStep(currentStep + 1);
      }, steps[Math.max(0, currentStep)]?.duration || 500);

      return () => clearTimeout(timer);
    } else if (!isProcessing) {
      setCurrentStep(-1);
    }
  }, [isProcessing, currentStep]);

  const canProcess = policy && weatherData.length > 0 && !isProcessing;

  const formatPayout = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'approved': return '#4caf50';
      case 'rejected': return '#f44336';
      case 'partial': return '#ff9800';
      default: return '#2196f3';
    }
  };

  return (
    <div className="claim-processor">
      {!result && (
        <>
          <div className="claim-controls">
            <button 
              className={`process-button ${isProcessing ? 'processing' : ''}`}
              onClick={onProcess}
              disabled={!canProcess}
            >
              {isProcessing ? 'Processing Claim...' : 'Process Insurance Claim'}
            </button>
          </div>

          {isProcessing && (
            <div className="processing-status">
              <div className="status-steps">
                {steps.map((step, index) => (
                  <div 
                    key={step.id} 
                    className={`step ${index <= currentStep ? 'active' : ''} ${index < currentStep ? 'completed' : ''}`}
                  >
                    <div className="step-icon">{step.icon}</div>
                    <span className="step-name">{step.name}</span>
                  </div>
                ))}
              </div>
              
              <div className="processing-info">
                <p>🔄 Processing claim using EigenLayer AVS consensus mechanism</p>
                <p>📊 Analyzing {weatherData.length} weather data points from {new Set(weatherData.map(d => d.source)).size} sources</p>
              </div>
            </div>
          )}

          {!policy && (
            <div className="placeholder">
              Select a policy and simulate weather data to process a claim
            </div>
          )}
        </>
      )}

      {result && (
        <div className="result-display">
          <div className="result-header">
            <div>
              <h3 className={`result-status ${result.claimStatus}`} style={{ color: getStatusColor(result.claimStatus) }}>
                Claim {result.claimStatus.toUpperCase()}
              </h3>
              <p className="claim-id">Claim ID: {result.claimId}</p>
            </div>
            <div className="processing-time">
              ⚡ Processed in {result.processingTime}ms
            </div>
          </div>

          <div className="result-details">
            <div className="payout-section">
              {result.payoutAmount > 0 ? (
                <div className="payout-info">
                  <h4>Payout Amount</h4>
                  <div className="payout-amount">{formatPayout(result.payoutAmount)}</div>
                  <p>Automated transfer initiated</p>
                </div>
              ) : (
                <div className="no-payout">
                  <h4>No Payout Required</h4>
                  <p>Weather conditions did not meet policy triggers</p>
                </div>
              )}
            </div>

            <div className="triggers-evaluation">
              <h4>Trigger Evaluation</h4>
              {result.triggeredPerils.map((trigger, index) => (
                <div key={index} className="trigger-result">
                  <span className={`trigger-status ${trigger.conditionsMet ? 'met' : 'not-met'}`}>
                    {trigger.conditionsMet ? '✅' : '❌'}
                  </span>
                  <span className="trigger-name">{trigger.peril.replace('_', ' ')}</span>
                  {trigger.conditionsMet && (
                    <span className="trigger-payout">
                      {(trigger.payoutRatio * 100).toFixed(0)}% payout triggered
                    </span>
                  )}
                </div>
              ))}
            </div>
          </div>

          <div className="verification-section">
            <h4>🔐 Cryptographic Verification</h4>
            <p className="verification-hash">Hash: {result.verificationHash}</p>
            <div className="verification-details">
              <span>✓ Consensus reached with {result.weatherData.length} data points</span>
              <span>✓ Signed by AVS operators</span>
              <span>✓ Immutable on-chain record</span>
            </div>
          </div>

          <div className="claim-benefits">
            <h4>💡 Benefits of AVS Processing</h4>
            <div className="benefits-grid">
              <div className="benefit">
                <span className="benefit-icon">⚡</span>
                <span>Instant Processing</span>
                <small>{result.processingTime}ms vs 30 days traditional</small>
              </div>
              <div className="benefit">
                <span className="benefit-icon">💰</span>
                <span>99% Lower Cost</span>
                <small>~$10 vs $1000+ traditional</small>
              </div>
              <div className="benefit">
                <span className="benefit-icon">🔒</span>
                <span>Zero Fraud Risk</span>
                <small>Cryptographic consensus</small>
              </div>
              <div className="benefit">
                <span className="benefit-icon">📊</span>
                <span>Fully Transparent</span>
                <small>All data verifiable on-chain</small>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default ClaimProcessor;