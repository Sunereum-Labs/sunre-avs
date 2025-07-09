package types

import (
	"time"
)

// Insurance-specific types for weather-triggered claims

type InsuranceType string

const (
	InsuranceCrop      InsuranceType = "crop"
	InsuranceProperty  InsuranceType = "property"
	InsuranceTravel    InsuranceType = "travel"
	InsuranceEvent     InsuranceType = "event"
	InsuranceTransport InsuranceType = "transport"
)

type WeatherPeril string

const (
	PerilHeatWave    WeatherPeril = "heat_wave"
	PerilColdSnap    WeatherPeril = "cold_snap"
	PerilDrought     WeatherPeril = "drought"
	PerilExcessRain  WeatherPeril = "excess_rain"
	PerilFrost       WeatherPeril = "frost"
	PerilHighWind    WeatherPeril = "high_wind"
	PerilHail        WeatherPeril = "hail"
)

type InsurancePolicy struct {
	PolicyID       string                  `json:"policy_id"`
	PolicyHolder   string                  `json:"policy_holder"`
	InsuranceType  InsuranceType           `json:"insurance_type"`
	Location       Location                `json:"location"`
	CoverageAmount float64                 `json:"coverage_amount"`
	Premium        float64                 `json:"premium"`
	StartDate      time.Time               `json:"start_date"`
	EndDate        time.Time               `json:"end_date"`
	Triggers       []InsuranceTrigger      `json:"triggers"`
	Metadata       map[string]interface{}  `json:"metadata"`
}

type InsuranceTrigger struct {
	TriggerID    string              `json:"trigger_id"`
	Peril        WeatherPeril        `json:"peril"`
	Conditions   TriggerConditions   `json:"conditions"`
	PayoutRatio  float64             `json:"payout_ratio"` // 0.0 to 1.0
	Description  string              `json:"description"`
}

type TriggerConditions struct {
	TemperatureMin      *float64      `json:"temperature_min,omitempty"`
	TemperatureMax      *float64      `json:"temperature_max,omitempty"`
	ConsecutiveDays     int           `json:"consecutive_days,omitempty"`
	HumidityMin         *float64      `json:"humidity_min,omitempty"`
	HumidityMax         *float64      `json:"humidity_max,omitempty"`
	WindSpeedMin        *float64      `json:"wind_speed_min,omitempty"`
	PrecipitationMin    *float64      `json:"precipitation_min,omitempty"`
	PrecipitationMax    *float64      `json:"precipitation_max,omitempty"`
	TimeWindow          TimeWindow    `json:"time_window,omitempty"`
}

type TimeWindow struct {
	StartMonth int `json:"start_month"` // 1-12
	EndMonth   int `json:"end_month"`   // 1-12
	StartHour  int `json:"start_hour"`  // 0-23
	EndHour    int `json:"end_hour"`    // 0-23
}

type InsuranceClaimRequest struct {
	PolicyID      string        `json:"policy_id"`
	Policy        InsurancePolicy `json:"policy"`
	ClaimDate     time.Time     `json:"claim_date"`
	AutomatedCheck bool          `json:"automated_check"`
}

type InsuranceClaimResponse struct {
	ClaimID          string                `json:"claim_id"`
	PolicyID         string                `json:"policy_id"`
	ClaimStatus      ClaimStatus           `json:"claim_status"`
	TriggeredPerils  []TriggeredPeril      `json:"triggered_perils"`
	PayoutAmount     float64               `json:"payout_amount"`
	WeatherData      WeatherAssessment     `json:"weather_data"`
	VerificationHash string                `json:"verification_hash"`
	Timestamp        time.Time             `json:"timestamp"`
}

type ClaimStatus string

const (
	ClaimApproved    ClaimStatus = "approved"
	ClaimRejected    ClaimStatus = "rejected"
	ClaimPartial     ClaimStatus = "partial"
	ClaimPending     ClaimStatus = "pending"
	ClaimInvestigate ClaimStatus = "investigate"
)

type TriggeredPeril struct {
	Peril            WeatherPeril    `json:"peril"`
	TriggerID        string          `json:"trigger_id"`
	ConditionsMet    bool            `json:"conditions_met"`
	PayoutRatio      float64         `json:"payout_ratio"`
	Evidence         WeatherEvidence `json:"evidence"`
}

type WeatherEvidence struct {
	DataPoints       []DataPoint     `json:"data_points"`
	AverageTemp      float64         `json:"average_temp"`
	MinTemp          float64         `json:"min_temp"`
	MaxTemp          float64         `json:"max_temp"`
	ConsecutiveDays  int             `json:"consecutive_days"`
	Confidence       float64         `json:"confidence"`
}

type WeatherAssessment struct {
	AssessmentPeriod DateRange       `json:"assessment_period"`
	LocationVerified bool            `json:"location_verified"`
	DataSources      int             `json:"data_sources"`
	ConsensusMethod  string          `json:"consensus_method"`
	WeatherSummary   WeatherSummary  `json:"weather_summary"`
}

type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type WeatherSummary struct {
	AverageTemperature float64         `json:"average_temperature"`
	MaxTemperature     float64         `json:"max_temperature"`
	MinTemperature     float64         `json:"min_temperature"`
	TotalPrecipitation float64         `json:"total_precipitation"`
	MaxWindSpeed       float64         `json:"max_wind_speed"`
	ExtremeEvents      []ExtremeEvent  `json:"extreme_events"`
}

type ExtremeEvent struct {
	Date        time.Time    `json:"date"`
	EventType   WeatherPeril `json:"event_type"`
	Severity    string       `json:"severity"`
	Description string       `json:"description"`
	Values      map[string]float64 `json:"values"`
}

// Predefined insurance product templates

var InsuranceTemplates = map[string]InsurancePolicy{
	"crop_heat_protection": {
		InsuranceType: InsuranceCrop,
		Triggers: []InsuranceTrigger{
			{
				Peril: PerilHeatWave,
				Conditions: TriggerConditions{
					TemperatureMax:  floatPtr(35.0),
					ConsecutiveDays: 3,
					TimeWindow: TimeWindow{
						StartMonth: 6,
						EndMonth:   8,
					},
				},
				PayoutRatio: 0.5,
				Description: "Heat stress protection for crops",
			},
			{
				Peril: PerilHeatWave,
				Conditions: TriggerConditions{
					TemperatureMax:  floatPtr(40.0),
					ConsecutiveDays: 2,
				},
				PayoutRatio: 1.0,
				Description: "Extreme heat protection",
			},
		},
	},
	"event_weather_insurance": {
		InsuranceType: InsuranceEvent,
		Triggers: []InsuranceTrigger{
			{
				Peril: PerilExcessRain,
				Conditions: TriggerConditions{
					PrecipitationMin: floatPtr(50.0), // mm in 24h
					TimeWindow: TimeWindow{
						StartHour: 8,
						EndHour:   20,
					},
				},
				PayoutRatio: 1.0,
				Description: "Event cancellation due to rain",
			},
			{
				Peril: PerilHighWind,
				Conditions: TriggerConditions{
					WindSpeedMin: floatPtr(60.0), // km/h
				},
				PayoutRatio: 1.0,
				Description: "Event cancellation due to high winds",
			},
		},
	},
	"travel_delay_insurance": {
		InsuranceType: InsuranceTravel,
		Triggers: []InsuranceTrigger{
			{
				Peril: PerilColdSnap,
				Conditions: TriggerConditions{
					TemperatureMax: floatPtr(-10.0),
				},
				PayoutRatio: 0.2, // Daily compensation
				Description: "Flight delays due to extreme cold",
			},
		},
	},
}

func floatPtr(f float64) *float64 {
	return &f
}