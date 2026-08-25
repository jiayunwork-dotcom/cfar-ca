package stats

import (
	"math"
)

func Rayleigh(rng randSource, n int, sigma float64) ([]float64, error) {
	if n < 0 {
		return nil, &BadCountError{N: n}
	}
	if sigma <= 0 {
		return nil, &BadMeanError{Mean: sigma}
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		u := rng.Float64()
		out[i] = sigma * math.Sqrt(-2.0*math.Log1p(-u))
	}
	return out, nil
}

func RayleighMean(sigma float64) float64 {
	return sigma * math.Sqrt(math.Pi/2.0)
}

func RayleighVariance(sigma float64) float64 {
	return (2.0 - math.Pi/2.0) * sigma * sigma
}

func RayleighCdf(t, sigma float64) float64 {
	if sigma <= 0 {
		return math.NaN()
	}
	if t <= 0 {
		return 0
	}
	return 1.0 - math.Exp(-t*t/(2.0*sigma*sigma))
}

func RayleighTailPfa(threshold, sigma float64) float64 {
	if threshold <= 0 || sigma <= 0 {
		return math.NaN()
	}
	return math.Exp(-threshold * threshold / (2.0 * sigma * sigma))
}

func SigmaFromMean(mean float64) float64 {
	if mean <= 0 {
		return math.NaN()
	}
	return mean / math.Sqrt(math.Pi/2.0)
}
