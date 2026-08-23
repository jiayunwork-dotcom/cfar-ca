package detect

import "cfar-ca/internal/model"

// ReferenceIndices 返回 CUT 下标 i 的参考单元下标列表。
//
// 布局（每侧 Guards 个保护单元 + Refs 个参考单元）：
//
//	参考        保护   CUT   保护        参考
//	[……r……] [g] [i] [g] [r……]
//
// 返回值 ok=false 表示 CUT 位于序列边缘、两侧参考窗不完整。
// 此时不借用半窗、不补零，直接把该 CUT 标为无效。
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

// HasFullSupport 判断 CUT 下标 i 是否有完整的双侧参考窗。
func HasFullSupport(cfg *model.DetectorConfig, i int) bool {
	_, _, ok := ReferenceIndices(cfg, i)
	return ok
}

// FirstValidIndex 返回第一个参考窗完整的 CUT 下标；没有则返回 -1。
func FirstValidIndex(cfg *model.DetectorConfig) int {
	for i := 0; i < cfg.SequenceLength(); i++ {
		if HasFullSupport(cfg, i) {
			return i
		}
	}
	return -1
}

// LastValidIndex 返回最后一个参考窗完整的 CUT 下标；没有则返回 -1。
func LastValidIndex(cfg *model.DetectorConfig) int {
	for i := cfg.SequenceLength() - 1; i >= 0; i-- {
		if HasFullSupport(cfg, i) {
			return i
		}
	}
	return -1
}

// EdgeCuts 返回无效（参考窗不完整）的 CUT 下标。
func EdgeCuts(cfg *model.DetectorConfig) []int {
	out := make([]int, 0, 2*cfg.WindowHalf())
	for i := 0; i < cfg.SequenceLength(); i++ {
		if !HasFullSupport(cfg, i) {
			out = append(out, i)
		}
	}
	return out
}

// GuardIndices 返回 CUT i 两侧保护单元的下标；参考窗不完整时 ok=false。
// 保护单元进不了参考均值——这是目标不被自身抬阈的关键保证。
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
