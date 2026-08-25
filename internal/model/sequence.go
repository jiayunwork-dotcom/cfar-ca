package model

import "math"

type Sequence struct {
	samples []float64
}

func NewSequence(samples []float64) *Sequence {
	return &Sequence{samples: samples}
}

func (s *Sequence) Len() int {
	return len(s.samples)
}

func (s *Sequence) At(i int) float64 {
	if i < 0 || i >= len(s.samples) {
		return math.NaN()
	}
	return s.samples[i]
}

func (s *Sequence) Samples() []float64 {
	out := make([]float64, len(s.samples))
	copy(out, s.samples)
	return out
}

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

func (s *Sequence) Sum() float64 {
	total := 0.0
	for _, v := range s.samples {
		total += v
	}
	return total
}

func (s *Sequence) Mean() float64 {
	if len(s.samples) == 0 {
		return math.NaN()
	}
	return s.Sum() / float64(len(s.samples))
}

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

func (s *Sequence) Sub(lo, hi int) []float64 {
	if lo < 0 {
		lo = 0
	}
	if hi > len(s.samples) {
		hi = len(s.samples)
	}
	out := make([]float64, hi-lo)
	copy(out, s.samples[lo:hi])
	return out
}

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
