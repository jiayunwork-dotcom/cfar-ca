package stats

import "math/rand"

func NewSource(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

type randSource interface {
	Float64() float64
}

func Uniform(rng *rand.Rand, n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = rng.Float64()
	}
	return out
}

func ScaleToMean(samples []float64, mean float64) ([]float64, error) {
	if len(samples) == 0 {
		return nil, &EmptySeriesError{}
	}
	if mean <= 0 {
		return nil, &BadMeanError{Mean: mean}
	}
	cur := 0.0
	for _, v := range samples {
		cur += v
	}
	curMean := cur / float64(len(samples))
	if curMean <= 0 {
		return nil, &EmptySeriesError{}
	}
	k := mean / curMean
	out := make([]float64, len(samples))
	for i, v := range samples {
		out[i] = v * k
	}
	return out, nil
}
