package consensus

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/Layr-Labs/hourglass-avs-template/internal/types"

	log "github.com/sirupsen/logrus"
)

type ConsensusEngine struct {
	MinSources   int
	madThreshold float64
}

func NewConsensusEngine(minSources int, madThreshold float64) *ConsensusEngine {
	return &ConsensusEngine{
		MinSources:   minSources,
		madThreshold: madThreshold,
	}
}

func (c *ConsensusEngine) ReachConsensus(dataPoints []types.DataPoint) (*types.ConsensusResult, error) {
	if len(dataPoints) < c.MinSources {
		return nil, fmt.Errorf("insufficient data sources: %d < %d", len(dataPoints), c.MinSources)
	}
	
	temperatures := make([]float64, len(dataPoints))
	for i, dp := range dataPoints {
		temperatures[i] = dp.Temperature
	}
	
	median := c.calculateMedian(temperatures)
	mad := c.calculateMAD(temperatures, median)
	
	filteredPoints := c.filterOutliers(dataPoints, median, mad)
	
	if len(filteredPoints) < c.MinSources {
		return nil, fmt.Errorf("too many outliers filtered: %d remaining < %d required", 
			len(filteredPoints), c.MinSources)
	}
	
	consensusTemp, confidence := c.weightedConsensus(filteredPoints)
	
	aggregatedSig := c.aggregateSignatures(filteredPoints)
	
	return &types.ConsensusResult{
		Temperature:    consensusTemp,
		Confidence:     confidence,
		DataPoints:     filteredPoints,
		AggregatedSig:  aggregatedSig,
		Timestamp:      time.Now(),
	}, nil
}

func (c *ConsensusEngine) calculateMedian(values []float64) float64 {
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	
	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

func (c *ConsensusEngine) calculateMAD(values []float64, median float64) float64 {
	deviations := make([]float64, len(values))
	for i, v := range values {
		deviations[i] = math.Abs(v - median)
	}
	
	return c.calculateMedian(deviations)
}

func (c *ConsensusEngine) filterOutliers(dataPoints []types.DataPoint, median, mad float64) []types.DataPoint {
	filtered := make([]types.DataPoint, 0)
	
	madThreshold := c.madThreshold
	if mad == 0 {
		mad = 0.01
	}
	
	for _, dp := range dataPoints {
		deviation := math.Abs(dp.Temperature - median)
		if deviation <= madThreshold*mad {
			filtered = append(filtered, dp)
		} else {
			log.Warnf("Filtered outlier from %s: %.2f (median: %.2f, deviation: %.2f, threshold: %.2f)",
				dp.Source, dp.Temperature, median, deviation, madThreshold*mad)
		}
	}
	
	return filtered
}

func (c *ConsensusEngine) weightedConsensus(dataPoints []types.DataPoint) (float64, float64) {
	if len(dataPoints) == 0 {
		return 0, 0
	}
	
	weights := c.calculateReliabilityWeights(dataPoints)
	
	weightedSum := 0.0
	totalWeight := 0.0
	
	for i, dp := range dataPoints {
		weight := weights[i]
		weightedSum += dp.Temperature * weight
		totalWeight += weight
	}
	
	if totalWeight == 0 {
		temperatures := make([]float64, len(dataPoints))
		for i, dp := range dataPoints {
			temperatures[i] = dp.Temperature
		}
		return c.calculateMedian(temperatures), 0.5
	}
	
	consensusTemp := weightedSum / totalWeight
	
	variance := 0.0
	for i, dp := range dataPoints {
		diff := dp.Temperature - consensusTemp
		variance += weights[i] * diff * diff
	}
	variance /= totalWeight
	
	stdDev := math.Sqrt(variance)
	confidence := 1.0 - math.Min(stdDev/10.0, 1.0)
	
	agreementScore := c.calculateAgreementScore(dataPoints, consensusTemp)
	confidence = (confidence + agreementScore) / 2
	
	return consensusTemp, confidence
}

func (c *ConsensusEngine) calculateReliabilityWeights(dataPoints []types.DataPoint) []float64 {
	weights := make([]float64, len(dataPoints))
	
	for i, dp := range dataPoints {
		weight := 1.0
		
		age := time.Since(dp.Timestamp).Minutes()
		if age > 5 {
			weight *= math.Max(0.5, 1.0-age/60.0)
		}
		
		if dp.Confidence > 0 {
			weight *= dp.Confidence
		}
		
		sourceMultiplier := c.getSourceReliabilityScore(dp.Source)
		weight *= sourceMultiplier
		
		weights[i] = weight
	}
	
	minWeight := 0.1
	for i := range weights {
		if weights[i] < minWeight {
			weights[i] = minWeight
		}
	}
	
	return weights
}

func (c *ConsensusEngine) getSourceReliabilityScore(source string) float64 {
	reliabilityScores := map[string]float64{
		"OpenWeatherMap":  0.95,
		"WeatherAPI":      0.93,
		"TomorrowIO":      0.92,
		"VisualCrossing":  0.90,
		"OpenMeteo":       0.88,
	}
	
	if score, ok := reliabilityScores[source]; ok {
		return score
	}
	return 0.8
}

func (c *ConsensusEngine) calculateAgreementScore(dataPoints []types.DataPoint, consensusTemp float64) float64 {
	if len(dataPoints) == 0 {
		return 0
	}
	
	totalDeviation := 0.0
	for _, dp := range dataPoints {
		totalDeviation += math.Abs(dp.Temperature - consensusTemp)
	}
	
	avgDeviation := totalDeviation / float64(len(dataPoints))
	
	agreementScore := math.Max(0, 1.0-avgDeviation/5.0)
	
	return agreementScore
}

func (c *ConsensusEngine) aggregateSignatures(dataPoints []types.DataPoint) []byte {
	h := sha256.New()
	
	for _, dp := range dataPoints {
		h.Write([]byte(dp.Source))
		h.Write([]byte(fmt.Sprintf("%.2f", dp.Temperature)))
		h.Write(dp.Signature)
	}
	
	return h.Sum(nil)
}

func (c *ConsensusEngine) VerifyThreshold(result *types.ConsensusResult, threshold float64) bool {
	result.MeetsThreshold = result.Temperature >= threshold
	return result.MeetsThreshold
}

func GenerateSignature(operatorID string, taskID string, temperature float64) []byte {
	data := fmt.Sprintf("%s:%s:%.2f", operatorID, taskID, temperature)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

func VerifySignature(signature []byte, operatorID string, taskID string, temperature float64) bool {
	expected := GenerateSignature(operatorID, taskID, temperature)
	return hex.EncodeToString(signature) == hex.EncodeToString(expected)
}

type ConsensusStats struct {
	MedianTemperature float64
	MAD               float64
	OutlierCount      int
	FilteredCount     int
	Confidence        float64
}

func (c *ConsensusEngine) GetConsensusStats(dataPoints []types.DataPoint) ConsensusStats {
	temperatures := make([]float64, len(dataPoints))
	for i, dp := range dataPoints {
		temperatures[i] = dp.Temperature
	}
	
	median := c.calculateMedian(temperatures)
	mad := c.calculateMAD(temperatures, median)
	
	outlierCount := 0
	for _, temp := range temperatures {
		deviation := math.Abs(temp - median)
		if deviation > c.madThreshold*mad {
			outlierCount++
		}
	}
	
	result, _ := c.ReachConsensus(dataPoints)
	confidence := 0.0
	if result != nil {
		confidence = result.Confidence
	}
	
	return ConsensusStats{
		MedianTemperature: median,
		MAD:               mad,
		OutlierCount:      outlierCount,
		FilteredCount:     len(dataPoints) - outlierCount,
		Confidence:        confidence,
	}
}