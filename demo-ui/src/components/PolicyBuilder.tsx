import React from 'react';
import { InsurancePolicy, InsuranceTrigger } from '../types/insurance';

interface PolicyBuilderProps {
  policy: InsurancePolicy;
  onChange: (policy: InsurancePolicy) => void;
  readOnly?: boolean;
}

const PolicyBuilder: React.FC<PolicyBuilderProps> = ({ policy, onChange, readOnly = false }) => {
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
  };

  const formatDate = (date: Date) => {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const getTriggerConditionDisplay = (trigger: InsuranceTrigger) => {
    const conditions = [];
    const c = trigger.conditions;

    if (c.temperatureMax !== undefined) {
      conditions.push(`Temp > ${c.temperatureMax}°C`);
    }
    if (c.temperatureMin !== undefined) {
      conditions.push(`Temp < ${c.temperatureMin}°C`);
    }
    if (c.consecutiveDays) {
      conditions.push(`${c.consecutiveDays} consecutive days`);
    }
    if (c.precipitationMin !== undefined) {
      conditions.push(`Rain > ${c.precipitationMin}mm`);
    }
    if (c.windSpeedMin !== undefined) {
      conditions.push(`Wind > ${c.windSpeedMin}km/h`);
    }

    return conditions;
  };

  const getPerilIcon = (peril: string) => {
    const icons: Record<string, string> = {
      heat_wave: '🌡️',
      cold_snap: '❄️',
      excess_rain: '🌧️',
      drought: '☀️',
      high_wind: '💨',
      frost: '🧊',
      hail: '🌨️'
    };
    return icons[peril] || '⚠️';
  };

  return (
    <div className="policy-builder">
      <div className="policy-details">
        <div className="policy-field">
          <label>Policy ID:</label>
          <span>{policy.policyId}</span>
        </div>

        <div className="policy-field">
          <label>Policy Holder:</label>
          <span>{policy.policyHolder}</span>
        </div>

        <div className="policy-field">
          <label>Insurance Type:</label>
          <span className="insurance-type">{policy.insuranceType.toUpperCase()}</span>
        </div>

        <div className="policy-field">
          <label>Location:</label>
          <span>{policy.location.city}, {policy.location.country} ({policy.location.latitude.toFixed(2)}, {policy.location.longitude.toFixed(2)})</span>
        </div>

        <div className="policy-field">
          <label>Coverage Amount:</label>
          <span className="coverage-amount">{formatCurrency(policy.coverageAmount)}</span>
        </div>

        <div className="policy-field">
          <label>Premium:</label>
          <span>{formatCurrency(policy.premium)}</span>
        </div>

        <div className="policy-field">
          <label>Policy Period:</label>
          <span>{formatDate(policy.startDate)} - {formatDate(policy.endDate)}</span>
        </div>

        <div className="triggers-section">
          <h3>Coverage Triggers</h3>
          <div className="triggers-list">
            {policy.triggers.map((trigger, index) => (
              <div key={trigger.triggerId} className="trigger-item">
                <div className="trigger-header">
                  <span className="trigger-icon">{getPerilIcon(trigger.peril)}</span>
                  <span className="trigger-peril">{trigger.peril.replace('_', ' ').toUpperCase()}</span>
                  <span className="trigger-payout">{(trigger.payoutRatio * 100).toFixed(0)}% Payout</span>
                </div>
                <p className="trigger-description">{trigger.description}</p>
                <div className="trigger-conditions">
                  {getTriggerConditionDisplay(trigger).map((condition, i) => (
                    <span key={i} className="condition-badge">{condition}</span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default PolicyBuilder;