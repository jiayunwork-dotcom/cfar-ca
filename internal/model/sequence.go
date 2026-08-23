package model

import "math"

// Sequence 是幅度序列的轻量视图，提供统计辅助方法。
// 序列值表示平方律检波后的幅度，非负。
type Sequence struct {
	samples []float64
}

// NewSequence 包装幅度切片。调用方应保证数值已通过配置校验。
func NewSequence(samples []float64) *Sequence {
	return &Sequence{samples: samples}
}

// Len 返回序列长度。
func (s *Sequence) Len() int {
	return len(s.samples)
}

// At 返回第 i 个单元的幅度。越界调用返回 NaN 而不是 panic，
// 便于输出层安全渲染无效 CUT。
func (s *Sequence) At(i int) float64 {
	if i < 0 || i >= len(s.samples) {
		return math.NaN()
	}
	return s.samples[i]
}

// Samples 返回底层切片的副本，避免外部别名修改。
func (s *Sequence) Samples() []float64 {
	out := make([]float64, len(s.samples))
	copy(out, s.samples)
	return out
}

// Min 返回序列最小值；空序列返回 NaN。
func (s *Sequence) Min() float64 {
	if len(s.samples) == 0 {
		return math.NaN()
	}
	m := s.samples[0]
	for _, v := range s.samples[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

// Max 返回序列最大值；空序列返回 NaN。
func (s *Sequence) Max() float64 {
	if len(s.samples) == 0 {
		return math.NaN()
	}
	m := s.samples[0]
	for _, v := range s.samples[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

// Sum 返回序列元素之和。
func (s *Sequence) Sum() float64 {
	total := 0.0
	for _, v := range s.samples {
		total += v
	}
	return total
}

// Mean 返回序列算术平均；空序列返回 NaN。
func (s *Sequence) Mean() float64 {
	if len(s.samples) == 0 {
		return math.NaN()
	}
	return s.Sum() / float64(len(s.samples))
}

// SumRange 返回 [lo, hi) 半开区间内的元素和；越界部分按 0 处理。
func (s *Sequence) SumRange(lo, hi int) float64 {
	if lo < 0 {
		lo = 0
	}
	if hi > len(s.samples) {
		hi = len(s.samples)
	}
	total := 0.0
	for i := lo; i < hi; i++ {
		total += s.samples[i]
	}
	return total
}

// Sub 返回 [lo, hi) 区间的元素副本，越界被裁剪。
func (s *Sequence) Sub(lo, hi int) []float64 {
	if lo < 0 {
		lo = 0
	}
	if hi > len(s.samples) {
		hi = len(s.samples)
	}
	return liveSeqView(s.samples[lo:hi])
}

// IsConstant 判断序列是否近似恒定（用于测试与诊断）。
func (s *Sequence) IsConstant() bool {
	if len(s.samples) == 0 {
		return true
	}
	first := s.samples[0]
	for _, v := range s.samples {
		if v != first {
			return false
		}
	}
	return true
}
