package stats

import "math"

func Mean(samples []float64) float64 {
	if len(samples) == 0 {
		return math.NaN()
	}
	total := 0.0
	for _, v := range samples {
		total += v
	}
	return total / float64(len(samples))
}

func Variance(samples []float64) float64 {
	n := len(samples)
	if n < 2 {
		return math.NaN()
	}
	m := Mean(samples)
	sum := 0.0
	for _, v := range samples {
		d := v - m
		sum += d * d
	}
	return sum / float64(n-1)
}

func StdDev(samples []float64) float64 {
	v := Variance(samples)
	if math.IsNaN(v) {
		return math.NaN()
	}
	return math.Sqrt(v)
}

func Min(samples []float64) float64 {
	if len(samples) == 0 {
		return math.NaN()
	}
	m := samples[0]
	for _, v := range samples[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func Max(samples []float64) float64 {
	if len(samples) == 0 {
		return math.NaN()
	}
	m := samples[0]
	for _, v := range samples[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func EstimateMean(samples []float64, theory float64) (estimate, relError float64) {
	m := Mean(samples)
	if theory == 0 {
		return m, math.NaN()
	}
	return m, math.Abs(m-theory) / theory
}

func EstimateVariance(samples []float64, theory float64) (estimate, relError float64) {
	v := Variance(samples)
	if theory == 0 {
		return v, math.NaN()
	}
	return v, math.Abs(v-theory) / theory
}
