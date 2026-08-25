package stats

import (
	"fmt"
	"math"
)

type EmptySeriesError struct{}

func (e *EmptySeriesError) Error() string {
	return "序列为空，无法统计"
}

type BadMeanError struct {
	Mean float64
}

func (e *BadMeanError) Error() string {
	return fmt.Sprintf("均值必须为正，得到 %v", e.Mean)
}

type BadRateError struct {
	Rate float64
}

func (e *BadRateError) Error() string {
	return fmt.Sprintf("概率必须在 (0,1) 区间，得到 %v", e.Rate)
}

func Exponential(rng randSource, n int, mean float64) ([]float64, error) {
	if n < 0 {
		return nil, &BadCountError{N: n}
	}
	if mean <= 0 {
		return nil, &BadMeanError{Mean: mean}
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		u := rng.Float64()
		out[i] = -mean * math.Log1p(-u)
	}
	return out, nil
}

type BadCountError struct {
	N int
}

func (e *BadCountError) Error() string {
	return fmt.Sprintf("抽样个数必须非负，得到 %d", e.N)
}

func ExponentialMeanVar(mean float64) (mu, variance float64) {
	return mean, mean * mean
}

func ExponentialPfa(threshold, mean float64) float64 {
	if threshold <= 0 || mean <= 0 {
		return math.NaN()
	}
	return math.Exp(-threshold / mean)
}

func ExponentialThreshold(pfa, mean float64) float64 {
	if pfa <= 0 || pfa >= 1 || mean <= 0 {
		return math.NaN()
	}
	return -mean * math.Log(pfa)
}
