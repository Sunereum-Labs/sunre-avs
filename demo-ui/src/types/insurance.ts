// TypeScript interfaces for SunRe AVS Insurance Demo

export type InsuranceType = 'crop' | 'event' | 'travel' | 'property' | 'energy';
export type WeatherPeril = 'heat_wave' | 'cold_snap' | 'drought' | 'excess_rain' | 'frost' | 'high_wind' | 'hail' | 'low_wind';
export type ClaimStatus = 'approved' | 'rejected' | 'partial' | 'pending' | 'investigate';

export interface Location {
  latitude: number;
  longitude: number;
  city: string;
  country: string;
}

export interface TriggerConditions {
  temperatureMin?: number;
  temperatureMax?: number;
  consecutiveDays?: number;
  humidityMin?: number;
  humidityMax?: number;
  windSpeedMin?: number;
  windSpeedMax?: number;
  precipitationMin?: number;
  precipitationMax?: number;
  timeWindow?: {
    startMonth: number;
    endMonth: number;
    startHour?: number;
    endHour?: number;
  };
}

export interface InsuranceTrigger {
  triggerId: string;
  peril: WeatherPeril;
  conditions: TriggerConditions;
  payoutRatio: number;
  description: string;
}

export interface InsurancePolicy {
  policyId: string;
  policyHolder: string;
  insuranceType: InsuranceType;
  location: Location;
  coverageAmount: number;
  premium: number;
  startDate: Date;
  endDate: Date;
  triggers: InsuranceTrigger[];
  metadata?: Record<string, any>;
}

export interface WeatherDataPoint {
  source: string;
  temperature: number;
  humidity?: number;
  windSpeed?: number;
  precipitation?: number;
  timestamp: Date;
  confidence: number;
}

export interface ClaimResult {
  claimId: string;
  policyId: string;
  claimStatus: ClaimStatus;
  triggeredPerils: {
    peril: WeatherPeril;
    conditionsMet: boolean;
    payoutRatio: number;
  }[];
  payoutAmount: number;
  weatherData: WeatherDataPoint[];
  verificationHash: string;
  timestamp: Date;
  processingTime: number;
}

export interface DemoScenario {
  id: string;
  name: string;
  description: string;
  policy: InsurancePolicy;
  weatherPattern: 'normal' | 'heat_wave' | 'cold_snap' | 'storm' | 'drought';
  expectedOutcome: {
    status: ClaimStatus;
    payout: number;
    reason: string;
  };
}

export interface AVSTaskRequest {
  type: 'insurance_claim';
  claimRequest: {
    policyId: string;
    policy: InsurancePolicy;
    claimDate: Date;
    automatedCheck: boolean;
  };
  demoMode: boolean;
  demoScenario?: string;
}

export interface AVSConnection {
  endpoint: string;
  status: 'connected' | 'disconnected' | 'connecting';
  networkType: 'devnet' | 'testnet' | 'mainnet';
}