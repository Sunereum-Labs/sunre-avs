import React from 'react';
import { DemoScenario } from '../types/insurance';

interface DemoScenarioSelectorProps {
  scenarios: DemoScenario[];
  onSelect: (scenario: DemoScenario) => void;
  selected: DemoScenario | null;
}

const DemoScenarioSelector: React.FC<DemoScenarioSelectorProps> = ({ scenarios, onSelect, selected }) => {
  const getScenarioIcon = (type: string) => {
    const icons: Record<string, string> = {
      crop: '🌾',
      event: '🎪',
      travel: '✈️',
      property: '🏠',
      energy: '⚡'
    };
    return icons[type] || '📋';
  };

  const getWeatherIcon = (pattern: string) => {
    const icons: Record<string, string> = {
      heat_wave: '🌡️',
      cold_snap: '❄️',
      storm: '⛈️',
      normal: '☀️'
    };
    return icons[pattern] || '🌤️';
  };

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      approved: '#4caf50',
      rejected: '#f44336',
      partial: '#ff9800'
    };
    return (
      <span 
        className="status-badge" 
        style={{ 
          background: colors[status] || '#2196f3',
          padding: '2px 8px',
          borderRadius: '12px',
          fontSize: '0.8rem',
          fontWeight: 'bold'
        }}
      >
        {status.toUpperCase()}
      </span>
    );
  };

  return (
    <div className="scenario-selector">
      <div className="scenario-grid">
        {scenarios.map(scenario => (
          <div 
            key={scenario.id}
            className={`scenario-card ${selected?.id === scenario.id ? 'selected' : ''}`}
            onClick={() => onSelect(scenario)}
          >
            <div className="scenario-header">
              <span className="scenario-icon">{getScenarioIcon(scenario.policy.insuranceType)}</span>
              <h3>{scenario.name}</h3>
            </div>
            
            <p className="scenario-description">{scenario.description}</p>
            
            <div className="scenario-details">
              <div className="detail-row">
                <span>Coverage:</span>
                <strong>${scenario.policy.coverageAmount.toLocaleString()}</strong>
              </div>
              <div className="detail-row">
                <span>Location:</span>
                <strong>{scenario.policy.location.city}</strong>
              </div>
              <div className="detail-row">
                <span>Weather:</span>
                <span>{getWeatherIcon(scenario.weatherPattern)} {scenario.weatherPattern.replace('_', ' ')}</span>
              </div>
            </div>
            
            <div className="expected">
              <div className="expected-header">Expected Outcome:</div>
              <div className="expected-result">
                {getStatusBadge(scenario.expectedOutcome.status)}
                {scenario.expectedOutcome.payout > 0 && (
                  <span className="expected-payout">
                    ${scenario.expectedOutcome.payout.toLocaleString()}
                  </span>
                )}
              </div>
              <p className="expected-reason">{scenario.expectedOutcome.reason}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default DemoScenarioSelector;