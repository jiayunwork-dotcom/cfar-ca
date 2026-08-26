package detect

import (
	"math"

	"cfar-ca/internal/model"
)

func Decide(amplitude, threshold float64) bool {
	if math.IsNaN(amplitude) || math.IsNaN(threshold) || threshold < 0 {
		return false
	}
	return amplitude > threshold
}

func MarginRatio(amplitude, threshold float64) float64 {
	if threshold <= 0 || math.IsNaN(amplitude) || math.IsNaN(threshold) {
		return math.NaN()
	}
	return amplitude / threshold
}

func Excess(amplitude, threshold float64) float64 {
	if math.IsNaN(amplitude) || math.IsNaN(threshold) {
		return math.NaN()
	}
	return amplitude - threshold
}

func ThresholdFromMean(alpha, mean float64) float64 {
	if math.IsNaN(alpha) || math.IsNaN(mean) || mean <= 0 {
		return 0
	}
	return alpha * mean
}

func DetectCell(cfg *model.DetectorConfig, i int, alpha float64) (threshold float64, detected bool, margin float64, valid bool) {
	mean, ok := ReferenceMean(cfg, i)
	if !ok {
		return math.NaN(), false, math.NaN(), false
	}
	threshold = ThresholdFromMean(alpha, mean)
	detected = Decide(cfg.Amplitude[i], threshold)
	margin = MarginRatio(cfg.Amplitude[i], threshold)
	return threshold, detected, margin, true
}
