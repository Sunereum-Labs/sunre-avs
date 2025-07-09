package types

import (
	"time"
)

type TemperatureTask struct {
	TaskID      string
	Location    Location
	Threshold   float64
	Timestamp   time.Time
	ChainID     uint64
}

type Location struct {
	Latitude  float64
	Longitude float64
	City      string
	Country   string
}

type DataPoint struct {
	Source      string
	Temperature float64
	Timestamp   time.Time
	Confidence  float64
	Signature   []byte
}

type ConsensusResult struct {
	TaskID           string
	Temperature      float64
	MeetsThreshold   bool
	Confidence       float64
	DataPoints       []DataPoint
	AggregatedSig    []byte
	Timestamp        time.Time
}

type WeatherAPIConfig struct {
	Name      string
	BaseURL   string
	APIKey    string
	RateLimit int
}

type Config struct {
	AVS struct {
		Aggregator struct {
			MinOperators      int           `yaml:"min_operators"`
			ResponseTimeout   time.Duration `yaml:"response_timeout"`
			ConsensusThreshold float64      `yaml:"consensus_threshold"`
		} `yaml:"aggregator"`
	} `yaml:"avs"`
	
	WeatherAPIs map[string]struct {
		BaseURL   string `yaml:"base_url"`
		RateLimit int    `yaml:"rate_limit"`
		APIKey    string `yaml:"api_key,omitempty"`
	} `yaml:"weather_apis"`
	
	Consensus struct {
		MinSources    int     `yaml:"min_sources"`
		MADThreshold  float64 `yaml:"mad_threshold"`
		CacheTTL      int     `yaml:"cache_ttl"`
	} `yaml:"consensus"`
}

type OperatorResponse struct {
	OperatorID string
	TaskID     string
	DataPoints []DataPoint
	Signature  []byte
	Timestamp  time.Time
}

type TaskDistribution struct {
	TaskID       string
	Task         TemperatureTask
	AssignedAPIs []string
	Deadline     time.Time
}

type WeatherResponse struct {
	Temperature float64
	Humidity    float64
	Pressure    float64
	WindSpeed   float64
	Timestamp   time.Time
	Source      string
}

type CacheEntry struct {
	Data      WeatherResponse
	ExpiresAt time.Time
}

type TaskStatus int

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusDistributed
	TaskStatusExecuting
	TaskStatusAggregating
	TaskStatusCompleted
	TaskStatusFailed
)

type TaskState struct {
	Task           TemperatureTask
	Status         TaskStatus
	Operators      []string
	Responses      []OperatorResponse
	ConsensusResult *ConsensusResult
	CreatedAt      time.Time
	UpdatedAt      time.Time
}