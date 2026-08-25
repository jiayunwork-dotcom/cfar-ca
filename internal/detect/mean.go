package detect

import "cfar-ca/internal/model"

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

func ReferenceMeanExcludingGuards(cfg *model.DetectorConfig, i int) (mean float64, ok bool) {
	return ReferenceMean(cfg, i)
}

func ReferenceMeanSeq(seq *model.Sequence, guards, refs, i int) (mean float64, ok bool) {
	if seq == nil {
		return 0, false
	}
	cfg := model.NewDetectorConfig(seq.Samples(), guards, refs, 0.001)
	return ReferenceMean(cfg, i)
}
