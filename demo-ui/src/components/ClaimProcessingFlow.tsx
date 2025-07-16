import React, { useState, useEffect } from 'react';
import { 
  PlayCircle, 
  Loader2, 
  CheckCircle, 
  XCircle,
  Cloud,
  Activity,
  Shield,
  Zap,
  ArrowRight
} from 'lucide-react';
import { DemoScenario } from '../types/insurance';

interface ClaimProcessingFlowProps {
  scenario: DemoScenario;
  isProcessing: boolean;
  onProcess: () => void;
}

interface ProcessingStep {
  id: string;
  name: string;
  description: string;
  icon: React.FC<any>;
  duration: number;
}

const steps: ProcessingStep[] = [
  {
    id: 'fetch',
    name: 'Fetching Weather Data',
    description: 'Collecting data from multiple sources',
    icon: Cloud,
    duration: 1000
  },
  {
    id: 'consensus',
    name: 'Running Consensus',
    description: 'Applying MAD algorithm',
    icon: Activity,
    duration: 1000
  },
  {
    id: 'validate',
    name: 'Validating Claim',
    description: 'Checking trigger conditions',
    icon: Shield,
    duration: 1000
  },
  {
    id: 'settle',
    name: 'Settlement',
    description: 'Processing payout',
    icon: Zap,
    duration: 500
  }
];

const ClaimProcessingFlow: React.FC<ClaimProcessingFlowProps> = ({ 
  scenario, 
  isProcessing, 
  onProcess 
}) => {
  const [currentStep, setCurrentStep] = useState<string | null>(null);
  const [completedSteps, setCompletedSteps] = useState<string[]>([]);
  const [result, setResult] = useState<any>(null);

  useEffect(() => {
    if (isProcessing && !currentStep) {
      processSteps();
    } else if (!isProcessing) {
      setCurrentStep(null);
      setCompletedSteps([]);
      setResult(null);
    }
  }, [isProcessing]);

  const processSteps = async () => {
    for (let i = 0; i < steps.length; i++) {
      const step = steps[i];
      setCurrentStep(step.id);
      await new Promise(resolve => setTimeout(resolve, step.duration));
      setCompletedSteps(prev => [...prev, step.id]);
    }
    
    // Set result
    setResult({
      status: scenario.expectedOutcome.status,
      payout: scenario.expectedOutcome.payout,
      reason: scenario.expectedOutcome.reason
    });
    setCurrentStep(null);
  };

  const getStepStatus = (stepId: string) => {
    if (completedSteps.includes(stepId)) return 'completed';
    if (currentStep === stepId) return 'active';
    return 'pending';
  };

  return (
    <div className="bg-dark-400/50 backdrop-blur-sm rounded-2xl p-8 mt-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h3 className="text-2xl font-bold mb-2">Claim Processing</h3>
          <p className="text-gray-400">
            Process {scenario.name} insurance claim through the AVS
          </p>
        </div>
        <button
          onClick={onProcess}
          disabled={isProcessing}
          className={`
            px-6 py-3 rounded-lg font-medium transition-all flex items-center space-x-2
            ${isProcessing 
              ? 'bg-gray-700 text-gray-400 cursor-not-allowed' 
              : 'bg-gradient-to-r from-primary-600 to-primary-700 hover:from-primary-700 hover:to-primary-800 text-white shadow-lg hover:shadow-xl'
            }
          `}
        >
          {isProcessing ? (
            <>
              <Loader2 className="w-5 h-5 animate-spin" />
              <span>Processing...</span>
            </>
          ) : (
            <>
              <PlayCircle className="w-5 h-5" />
              <span>Process Claim</span>
            </>
          )}
        </button>
      </div>

      {/* Processing Steps */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        {steps.map((step, index) => {
          const status = getStepStatus(step.id);
          const Icon = step.icon;
          
          return (
            <div key={step.id} className="relative">
              <div className={`
                p-6 rounded-xl transition-all text-center
                ${status === 'completed' ? 'bg-green-900/20 border border-green-800/50' : ''}
                ${status === 'active' ? 'bg-primary-900/20 border border-primary-600 animate-pulse' : ''}
                ${status === 'pending' ? 'bg-dark-500/30 border border-gray-800' : ''}
              `}>
                <div className={`
                  w-16 h-16 mx-auto mb-4 rounded-full flex items-center justify-center
                  ${status === 'completed' ? 'bg-green-900/50' : ''}
                  ${status === 'active' ? 'bg-primary-600/30' : ''}
                  ${status === 'pending' ? 'bg-dark-500/50' : ''}
                `}>
                  {status === 'completed' ? (
                    <CheckCircle className="w-8 h-8 text-green-400" />
                  ) : status === 'active' ? (
                    <Loader2 className="w-8 h-8 text-primary-400 animate-spin" />
                  ) : (
                    <Icon className="w-8 h-8 text-gray-500" />
                  )}
                </div>
                <h4 className="font-semibold mb-1">{step.name}</h4>
                <p className="text-xs text-gray-400">{step.description}</p>
              </div>
              
              {index < steps.length - 1 && (
                <ArrowRight className={`
                  absolute top-1/2 -right-6 transform -translate-y-1/2 w-4 h-4 hidden md:block
                  ${completedSteps.includes(steps[index + 1]?.id) ? 'text-green-500' : 'text-gray-600'}
                `} />
              )}
            </div>
          );
        })}
      </div>

      {/* Result */}
      {result && (
        <div className={`
          p-6 rounded-xl border-2 animate-in
          ${result.status === 'approved' 
            ? 'bg-green-900/20 border-green-600' 
            : 'bg-red-900/20 border-red-600'
          }
        `}>
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center space-x-3">
              {result.status === 'approved' ? (
                <CheckCircle className="w-8 h-8 text-green-400" />
              ) : (
                <XCircle className="w-8 h-8 text-red-400" />
              )}
              <div>
                <h4 className="text-xl font-bold">
                  Claim {result.status === 'approved' ? 'Approved' : 'Rejected'}
                </h4>
                <p className="text-gray-400">{result.reason}</p>
              </div>
            </div>
            {result.payout > 0 && (
              <div className="text-right">
                <p className="text-sm text-gray-400">Payout Amount</p>
                <p className="text-3xl font-bold text-green-400">
                  ${result.payout.toLocaleString()}
                </p>
              </div>
            )}
          </div>
          
          <div className="grid grid-cols-3 gap-4 mt-6 pt-6 border-t border-gray-700">
            <div className="text-center">
              <p className="text-sm text-gray-400">Processing Time</p>
              <p className="text-lg font-semibold">3.5s</p>
            </div>
            <div className="text-center">
              <p className="text-sm text-gray-400">Consensus Sources</p>
              <p className="text-lg font-semibold">3</p>
            </div>
            <div className="text-center">
              <p className="text-sm text-gray-400">Confidence</p>
              <p className="text-lg font-semibold">98.5%</p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default ClaimProcessingFlow;