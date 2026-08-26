package detect

import "cfar-ca/internal/model"

func ReferenceIndices(cfg *model.DetectorConfig, i int) (left, right []int, ok bool) {
	if cfg == nil || i < 0 || i >= cfg.SequenceLength() {
		return nil, nil, false
	}
	lo := i - cfg.Guards - cfg.Refs
	hi := i + cfg.Guards + cfg.Refs
	if lo < 0 || hi > cfg.SequenceLength()-1 {
		return nil, nil, false
	}
	left = make([]int, cfg.Refs)
	for k := 0; k < cfg.Refs; k++ {
		left[k] = lo + k
	}
	right = make([]int, cfg.Refs)
	for k := 0; k < cfg.Refs; k++ {
		right[k] = i + cfg.Guards + 1 + k
	}
	return left, right, true
}

func HasFullSupport(cfg *model.DetectorConfig, i int) bool {
	_, _, ok := ReferenceIndices(cfg, i)
	return ok
}

func FirstValidIndex(cfg *model.DetectorConfig) int {
	for i := 0; i < cfg.SequenceLength(); i++ {
		if HasFullSupport(cfg, i) {
			return i
		}
	}
	return -1
}

func LastValidIndex(cfg *model.DetectorConfig) int {
	for i := cfg.SequenceLength() - 1; i >= 0; i-- {
		if HasFullSupport(cfg, i) {
			return i
		}
	}
	return -1
}

func EdgeCuts(cfg *model.DetectorConfig) []int {
	out := make([]int, 0, 2*cfg.WindowHalf())
	for i := 0; i < cfg.SequenceLength(); i++ {
		if !HasFullSupport(cfg, i) {
			out = append(out, i)
		}
	}
	return out
}

func GuardIndices(cfg *model.DetectorConfig, i int) (left, right []int, ok bool) {
	if cfg == nil || i < 0 || i >= cfg.SequenceLength() {
		return nil, nil, false
	}
	left = make([]int, cfg.Guards)
	for k := 0; k < cfg.Guards; k++ {
		left[k] = i - cfg.Guards + k
	}
	right = make([]int, cfg.Guards)
	for k := 0; k < cfg.Guards; k++ {
		right[k] = i + 1 + k
	}
	return left, right, true
}
