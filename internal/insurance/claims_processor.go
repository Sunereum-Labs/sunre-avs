package insurance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/Layr-Labs/hourglass-avs-template/internal/consensus"
	"github.com/Layr-Labs/hourglass-avs-template/internal/types"
	log "github.com/sirupsen/logrus"
)

type ClaimsProcessor struct {
	consensusEngine *consensus.ConsensusEngine
}

func NewClaimsProcessor(consensusEngine *consensus.ConsensusEngine) *ClaimsProcessor {
	return &ClaimsProcessor{
		consensusEngine: consensusEngine,
	}
}

func (cp *ClaimsProcessor) ProcessClaim(
	policy types.InsurancePolicy,
	weatherData []types.DataPoint,
	claimDate time.Time,
) (*types.InsuranceClaimResponse, error) {
	
	log.Infof("Processing insurance claim for policy %s, type: %s", 
		policy.PolicyID, policy.InsuranceType)
	
	// Verify claim date is within policy period
	if claimDate.Before(policy.StartDate) || claimDate.After(policy.EndDate) {
		return &types.InsuranceClaimResponse{
			ClaimID:     cp.generateClaimID(policy.PolicyID, claimDate),
			PolicyID:    policy.PolicyID,
			ClaimStatus: types.ClaimRejected,
			Timestamp:   time.Now(),
		}, fmt.Errorf("claim date outside policy period")
	}
	
	// Analyze weather data against policy triggers
	triggeredPerils := cp.evaluateTriggers(policy, weatherData, claimDate)
	
	// Calculate payout
	payoutAmount := cp.calculatePayout(policy, triggeredPerils)
	
	// Determine claim status
	claimStatus := cp.determineClaimStatus(triggeredPerils, payoutAmount)
	
	// Create weather assessment
	assessment := cp.createWeatherAssessment(weatherData, claimDate)
	
	// Generate verification hash
	verificationHash := cp.generateVerificationHash(policy, weatherData, triggeredPerils)
	
	response := &types.InsuranceClaimResponse{
		ClaimID:          cp.generateClaimID(policy.PolicyID, claimDate),
		PolicyID:         policy.PolicyID,
		ClaimStatus:      claimStatus,
		TriggeredPerils:  triggeredPerils,
		PayoutAmount:     payoutAmount,
		WeatherData:      assessment,
		VerificationHash: verificationHash,
		Timestamp:        time.Now(),
	}
	
	log.Infof("Claim processed: Status=%s, Payout=%.2f, Perils=%d", 
		claimStatus, payoutAmount, len(triggeredPerils))
	
	return response, nil
}

func (cp *ClaimsProcessor) evaluateTriggers(
	policy types.InsurancePolicy,
	weatherData []types.DataPoint,
	claimDate time.Time,
) []types.TriggeredPeril {
	
	var triggeredPerils []types.TriggeredPeril
	
	for _, trigger := range policy.Triggers {
		log.Debugf("Evaluating trigger %s for peril %s", trigger.TriggerID, trigger.Peril)
		
		conditionsMet, evidence := cp.checkTriggerConditions(
			trigger, 
			weatherData, 
			claimDate,
		)
		
		if conditionsMet {
			triggeredPerils = append(triggeredPerils, types.TriggeredPeril{
				Peril:         trigger.Peril,
				TriggerID:     trigger.TriggerID,
				ConditionsMet: true,
				PayoutRatio:   trigger.PayoutRatio,
				Evidence:      evidence,
			})
			
			log.Infof("Trigger activated: %s - %s", trigger.Peril, trigger.Description)
		}
	}
	
	return triggeredPerils
}

func (cp *ClaimsProcessor) checkTriggerConditions(
	trigger types.InsuranceTrigger,
	weatherData []types.DataPoint,
	claimDate time.Time,
) (bool, types.WeatherEvidence) {
	
	conditions := trigger.Conditions
	evidence := types.WeatherEvidence{
		DataPoints: weatherData,
	}
	
	// Calculate statistics from weather data
	temps := make([]float64, len(weatherData))
	for i, dp := range weatherData {
		temps[i] = dp.Temperature
	}
	
	if len(temps) == 0 {
		return false, evidence
	}
	
	// Calculate temperature statistics
	avgTemp := average(temps)
	minTemp := min(temps)
	maxTemp := max(temps)
	
	evidence.AverageTemp = avgTemp
	evidence.MinTemp = minTemp
	evidence.MaxTemp = maxTemp
	
	// Check temperature conditions
	if conditions.TemperatureMin != nil && minTemp < *conditions.TemperatureMin {
		return false, evidence
	}
	
	if conditions.TemperatureMax != nil && maxTemp > *conditions.TemperatureMax {
		// For heat-related perils, check consecutive days
		if conditions.ConsecutiveDays > 0 {
			consecutiveDays := cp.checkConsecutiveDays(
				weatherData, 
				*conditions.TemperatureMax,
				conditions.ConsecutiveDays,
			)
			evidence.ConsecutiveDays = consecutiveDays
			
			if consecutiveDays < conditions.ConsecutiveDays {
				return false, evidence
			}
		}
		
		// Check time window if specified
		if !cp.isInTimeWindow(claimDate, conditions.TimeWindow) {
			return false, evidence
		}
		
		return true, evidence
	}
	
	// Add more condition checks here (wind, precipitation, etc.)
	// For demo, we're focusing on temperature-based triggers
	
	return false, evidence
}

func (cp *ClaimsProcessor) checkConsecutiveDays(
	weatherData []types.DataPoint,
	threshold float64,
	requiredDays int,
) int {
	// Simplified: assume each data point represents daily data
	consecutiveCount := 0
	maxConsecutive := 0
	
	for _, dp := range weatherData {
		if dp.Temperature > threshold {
			consecutiveCount++
			if consecutiveCount > maxConsecutive {
				maxConsecutive = consecutiveCount
			}
		} else {
			consecutiveCount = 0
		}
	}
	
	return maxConsecutive
}

func (cp *ClaimsProcessor) isInTimeWindow(date time.Time, window types.TimeWindow) bool {
	if window.StartMonth == 0 && window.EndMonth == 0 {
		return true // No time window restriction
	}
	
	month := int(date.Month())
	
	if window.StartMonth <= window.EndMonth {
		return month >= window.StartMonth && month <= window.EndMonth
	}
	
	// Handle wrap-around (e.g., Nov-Feb)
	return month >= window.StartMonth || month <= window.EndMonth
}

func (cp *ClaimsProcessor) calculatePayout(
	policy types.InsurancePolicy,
	triggeredPerils []types.TriggeredPeril,
) float64 {
	
	if len(triggeredPerils) == 0 {
		return 0
	}
	
	// Use the highest payout ratio among triggered perils
	maxPayoutRatio := 0.0
	for _, peril := range triggeredPerils {
		if peril.PayoutRatio > maxPayoutRatio {
			maxPayoutRatio = peril.PayoutRatio
		}
	}
	
	return policy.CoverageAmount * maxPayoutRatio
}

func (cp *ClaimsProcessor) determineClaimStatus(
	triggeredPerils []types.TriggeredPeril,
	payoutAmount float64,
) types.ClaimStatus {
	
	if len(triggeredPerils) == 0 {
		return types.ClaimRejected
	}
	
	if payoutAmount == 0 {
		return types.ClaimRejected
	}
	
	// Check confidence levels
	for _, peril := range triggeredPerils {
		if peril.Evidence.Confidence < 0.7 {
			return types.ClaimInvestigate
		}
	}
	
	// Check if partial payout
	hasPartialPayout := false
	for _, peril := range triggeredPerils {
		if peril.PayoutRatio < 1.0 {
			hasPartialPayout = true
			break
		}
	}
	
	if hasPartialPayout {
		return types.ClaimPartial
	}
	
	return types.ClaimApproved
}

func (cp *ClaimsProcessor) createWeatherAssessment(
	weatherData []types.DataPoint,
	claimDate time.Time,
) types.WeatherAssessment {
	
	// Calculate assessment period
	startDate := claimDate.AddDate(0, 0, -7) // 7 days before claim
	endDate := claimDate
	
	// Calculate weather summary
	temps := make([]float64, len(weatherData))
	for i, dp := range weatherData {
		temps[i] = dp.Temperature
	}
	
	avgTemp := average(temps)
	maxTemp := max(temps)
	minTemp := min(temps)
	
	// Detect extreme events
	extremeEvents := cp.detectExtremeEvents(weatherData)
	
	return types.WeatherAssessment{
		AssessmentPeriod: types.DateRange{
			Start: startDate,
			End:   endDate,
		},
		LocationVerified: true,
		DataSources:      countUniqueSources(weatherData),
		ConsensusMethod:  "MAD (Median Absolute Deviation)",
		WeatherSummary: types.WeatherSummary{
			AverageTemperature: avgTemp,
			MaxTemperature:     maxTemp,
			MinTemperature:     minTemp,
			ExtremeEvents:      extremeEvents,
		},
	}
}

func (cp *ClaimsProcessor) detectExtremeEvents(
	weatherData []types.DataPoint,
) []types.ExtremeEvent {
	
	var events []types.ExtremeEvent
	
	// Simple extreme heat detection
	for _, dp := range weatherData {
		if dp.Temperature > 35.0 {
			events = append(events, types.ExtremeEvent{
				Date:        dp.Timestamp,
				EventType:   types.PerilHeatWave,
				Severity:    cp.getHeatSeverity(dp.Temperature),
				Description: fmt.Sprintf("High temperature: %.1f°C", dp.Temperature),
				Values: map[string]float64{
					"temperature": dp.Temperature,
				},
			})
		}
		
		if dp.Temperature < -10.0 {
			events = append(events, types.ExtremeEvent{
				Date:        dp.Timestamp,
				EventType:   types.PerilColdSnap,
				Severity:    cp.getColdSeverity(dp.Temperature),
				Description: fmt.Sprintf("Low temperature: %.1f°C", dp.Temperature),
				Values: map[string]float64{
					"temperature": dp.Temperature,
				},
			})
		}
	}
	
	return events
}

func (cp *ClaimsProcessor) getHeatSeverity(temp float64) string {
	switch {
	case temp > 45:
		return "extreme"
	case temp > 40:
		return "severe"
	case temp > 35:
		return "moderate"
	default:
		return "mild"
	}
}

func (cp *ClaimsProcessor) getColdSeverity(temp float64) string {
	switch {
	case temp < -20:
		return "extreme"
	case temp < -15:
		return "severe"
	case temp < -10:
		return "moderate"
	default:
		return "mild"
	}
}

func (cp *ClaimsProcessor) generateClaimID(policyID string, claimDate time.Time) string {
	data := fmt.Sprintf("%s-%d", policyID, claimDate.Unix())
	hash := sha256.Sum256([]byte(data))
	return "CLM-" + hex.EncodeToString(hash[:8])
}

func (cp *ClaimsProcessor) generateVerificationHash(
	policy types.InsurancePolicy,
	weatherData []types.DataPoint,
	triggeredPerils []types.TriggeredPeril,
) string {
	
	h := sha256.New()
	h.Write([]byte(policy.PolicyID))
	
	for _, dp := range weatherData {
		h.Write([]byte(fmt.Sprintf("%.2f", dp.Temperature)))
		h.Write([]byte(dp.Source))
	}
	
	for _, peril := range triggeredPerils {
		h.Write([]byte(string(peril.Peril)))
		h.Write([]byte(fmt.Sprintf("%.2f", peril.PayoutRatio)))
	}
	
	return hex.EncodeToString(h.Sum(nil))
}

// Helper functions

func average(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	sum := 0.0
	for _, n := range nums {
		sum += n
	}
	return sum / float64(len(nums))
}

func min(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	minVal := nums[0]
	for _, n := range nums[1:] {
		if n < minVal {
			minVal = n
		}
	}
	return minVal
}

func max(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	maxVal := nums[0]
	for _, n := range nums[1:] {
		if n > maxVal {
			maxVal = n
		}
	}
	return maxVal
}

func countUniqueSources(dataPoints []types.DataPoint) int {
	sources := make(map[string]bool)
	for _, dp := range dataPoints {
		sources[dp.Source] = true
	}
	return len(sources)
}

// Demo helper to simulate historical weather data
func GenerateDemoWeatherData(location types.Location, days int, scenario string) []types.DataPoint {
	var dataPoints []types.DataPoint
	baseTime := time.Now().AddDate(0, 0, -days)
	
	for i := 0; i < days; i++ {
		timestamp := baseTime.AddDate(0, 0, i)
		
		var temp float64
		switch scenario {
		case "heat_wave":
			// Simulate heat wave with 5 consecutive days above 35°C
			if i >= 2 && i <= 6 {
				temp = 36.0 + float64(i-2)*1.5 + math.Sin(float64(i))*2
			} else {
				temp = 28.0 + math.Sin(float64(i))*3
			}
			
		case "cold_snap":
			// Simulate cold snap
			if i >= 3 && i <= 5 {
				temp = -12.0 - float64(i-3)*2
			} else {
				temp = 5.0 + math.Sin(float64(i))*3
			}
			
		case "normal":
			// Normal weather pattern
			temp = 20.0 + math.Sin(float64(i)*0.5)*5
			
		default:
			temp = 22.0
		}
		
		// Add data from multiple sources for consensus
		sources := []string{"OpenMeteo", "WeatherAPI", "VisualCrossing"}
		for _, source := range sources {
			dp := types.DataPoint{
				Source:      source,
				Temperature: temp + (math.Sin(float64(i)+float64(len(source)))*0.5), // Small variation
				Timestamp:   timestamp,
				Confidence:  0.9 + math.Sin(float64(i))*0.05,
			}
			dp.Signature = []byte(fmt.Sprintf("sig-%s-%d", source, i))
			dataPoints = append(dataPoints, dp)
		}
	}
	
	return dataPoints
}