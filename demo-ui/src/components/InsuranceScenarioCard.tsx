import React from 'react';
import { 
  Zap, 
  CloudRain, 
  Sun, 
  Snowflake, 
  Wind,
  CheckCircle,
  XCircle,
  AlertCircle
} from 'lucide-react';
import { DemoScenario } from '../types/insurance';

interface InsuranceScenarioCardProps {
  scenario: DemoScenario;
  isSelected: boolean;
  onSelect: () => void;
}

const InsuranceScenarioCard: React.FC<InsuranceScenarioCardProps> = ({ 
  scenario, 
  isSelected, 
  onSelect 
}) => {
  const getIcon = () => {
    switch (scenario.weatherPattern) {
      case 'heat_wave':
        return <Sun className="w-8 h-8 text-orange-400" />;
      case 'storm':
        return <CloudRain className="w-8 h-8 text-blue-400" />;
      case 'cold_snap':
        return <Snowflake className="w-8 h-8 text-cyan-400" />;
      case 'drought':
        return <Wind className="w-8 h-8 text-yellow-400" />;
      default:
        return <Zap className="w-8 h-8 text-gray-400" />;
    }
  };

  const getStatusIcon = () => {
    switch (scenario.expectedOutcome.status) {
      case 'approved':
        return <CheckCircle className="w-5 h-5 text-green-500" />;
      case 'rejected':
        return <XCircle className="w-5 h-5 text-red-500" />;
      default:
        return <AlertCircle className="w-5 h-5 text-yellow-500" />;
    }
  };

  const getStatusColor = () => {
    switch (scenario.expectedOutcome.status) {
      case 'approved':
        return 'text-green-400 bg-green-900/20 border-green-800/50';
      case 'rejected':
        return 'text-red-400 bg-red-900/20 border-red-800/50';
      default:
        return 'text-yellow-400 bg-yellow-900/20 border-yellow-800/50';
    }
  };

  return (
    <div
      onClick={onSelect}
      className={`
        relative overflow-hidden rounded-2xl p-6 cursor-pointer transition-all duration-300
        ${isSelected 
          ? 'bg-gradient-to-br from-primary-900/30 to-primary-800/20 border-2 border-primary-600 shadow-lg shadow-primary-600/20' 
          : 'bg-dark-400/50 border border-gray-800 hover:border-gray-600 hover:shadow-lg'
        }
      `}
    >
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-start space-x-4">
          <div className={`p-3 rounded-lg ${isSelected ? 'bg-primary-600/20' : 'bg-dark-500/50'}`}>
            {getIcon()}
          </div>
          <div>
            <h3 className="text-lg font-semibold mb-1">{scenario.name}</h3>
            <p className="text-sm text-gray-400">{scenario.policy.insuranceType}</p>
          </div>
        </div>
        {isSelected && (
          <div className="absolute top-4 right-4 w-3 h-3 bg-primary-500 rounded-full animate-pulse" />
        )}
      </div>

      <p className="text-gray-300 text-sm mb-4">{scenario.description}</p>

      <div className="space-y-3">
        <div className="flex items-center justify-between text-sm">
          <span className="text-gray-400">Coverage</span>
          <span className="font-semibold">${scenario.policy.coverageAmount.toLocaleString()}</span>
        </div>
        <div className="flex items-center justify-between text-sm">
          <span className="text-gray-400">Expected Outcome</span>
          <div className="flex items-center space-x-2">
            {getStatusIcon()}
            <span className={`font-medium capitalize`}>
              {scenario.expectedOutcome.status}
            </span>
          </div>
        </div>
        {scenario.expectedOutcome.payout > 0 && (
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-400">Expected Payout</span>
            <span className="font-semibold text-green-400">
              ${scenario.expectedOutcome.payout.toLocaleString()}
            </span>
          </div>
        )}
      </div>

      <div className={`mt-4 p-3 rounded-lg border ${getStatusColor()}`}>
        <p className="text-xs">{scenario.expectedOutcome.reason}</p>
      </div>

      {isSelected && (
        <div className="absolute bottom-0 left-0 right-0 h-1 bg-gradient-to-r from-primary-500 to-primary-600" />
      )}
    </div>
  );
};

export default InsuranceScenarioCard;