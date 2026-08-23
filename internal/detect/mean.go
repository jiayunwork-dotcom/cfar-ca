package detect

import "cfar-ca/internal/model"

// ReferenceMean 计算 CUT i 两侧参考单元（不含保护单元与 CUT 本身）
// 的算术平均。参考窗不完整时返回 ok=false。
//
// 参考均值只来自参考单元：保护单元一律排除，CUT 自身更是进不了窗，
// 否则目标会把邻近阈值抬到自己头上（漏检）。
func ReferenceMean(cfg *model.DetectorConfig, i int) (mean float64, ok bool) {
	left, right, valid := ReferenceIndices(cfg, i)
	if !valid {
		return 0, false
	}
	total := 0.0
	for _, idx := range left {
		total += cfg.Amplitude[idx]
	}
	for _, idx := range right {
		total += cfg.Amplitude[idx]
	}
	return total / float64(2*cfg.Refs), true
}

// ReferenceSum 计算 CUT i 参考单元之和；参考窗不完整返回 ok=false。
func ReferenceSum(cfg *model.DetectorConfig, i int) (sum float64, ok bool) {
	left, right, valid := ReferenceIndices(cfg, i)
	if !valid {
		return 0, false
	}
	total := 0.0
	for _, idx := range left {
		total += cfg.Amplitude[idx]
	}
	for _, idx := range right {
		total += cfg.Amplitude[idx]
	}
	return total, true
}

// ReferenceMeanExcludingGuards 是 ReferenceMean 的别名，强调保护单元
// 不进入均值（与测试命名对应）。
func ReferenceMeanExcludingGuards(cfg *model.DetectorConfig, i int) (mean float64, ok bool) {
	return ReferenceMean(cfg, i)
}

// ReferenceMeanSeq 以 Sequence 视图计算参考均值（供不持有 Config 的场景用）。
func ReferenceMeanSeq(seq *model.Sequence, guards, refs, i int) (mean float64, ok bool) {
	if seq == nil {
		return 0, false
	}
	cfg := model.NewDetectorConfig(seq.Samples(), guards, refs, 0.001)
	// 仅用几何信息，Pfa 不影响均值。
	return ReferenceMean(cfg, i)
}
