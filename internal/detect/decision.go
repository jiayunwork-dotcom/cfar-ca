package detect

import (
	"math"

	"cfar-ca/internal/model"
)

// Decide 判定一个有效 CUT 是否越过阈值：幅度严格大于阈值才算检出。
// 等于阈值不算检出（与教科书 CFAR 保持一致）。
func Decide(amplitude, threshold float64) bool {
	if math.IsNaN(amplitude) || math.IsNaN(threshold) || threshold < 0 {
		return false
	}
	return amplitude > threshold
}

// MarginRatio 返回检测裕量 = 幅度/阈值。阈值非正或输入非法返回 NaN。
// 裕量 > 1 表示越过阈值，越大越稳。
func MarginRatio(amplitude, threshold float64) float64 {
	if threshold <= 0 || math.IsNaN(amplitude) || math.IsNaN(threshold) {
		return math.NaN()
	}
	return amplitude / threshold
}

// Excess 返回超出阈值的绝对量 = 幅度 − 阈值；未越过为负值。
func Excess(amplitude, threshold float64) float64 {
	if math.IsNaN(amplitude) || math.IsNaN(threshold) {
		return math.NaN()
	}
	return amplitude - threshold
}

// ThresholdFromMean 计算阈值 = α × 参考均值。参考均值非正（例如全零
// 序列）返回 0，此时任何幅度 > 0 都会检出。
func ThresholdFromMean(alpha, mean float64) float64 {
	if math.IsNaN(alpha) || math.IsNaN(mean) || mean <= 0 {
		return 0
	}
	return alpha * mean
}

// DetectCell 组合求值一个 CUT：返回 (阈值, 是否检出, 裕量, 是否有效)。
// 无效 CUT 阈值与裕量为 NaN，检出恒为 false。
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
