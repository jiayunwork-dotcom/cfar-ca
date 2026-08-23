package stats

import "math"

// Mean 计算样本均值；空序列返回 NaN。
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

// Variance 计算样本方差（除以 n−1）；样本数少于 2 返回 NaN。
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

// StdDev 返回样本标准差。
func StdDev(samples []float64) float64 {
	v := Variance(samples)
	if math.IsNaN(v) {
		return math.NaN()
	}
	return math.Sqrt(v)
}

// Min 返回最小值；空序列返回 NaN。
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

// Max 返回最大值；空序列返回 NaN。
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

// EstimateMean 估计序列均值并给出与理论值的偏差百分比。
func EstimateMean(samples []float64, theory float64) (estimate, relError float64) {
	m := Mean(samples)
	if theory == 0 {
		return m, math.NaN()
	}
	return m, math.Abs(m-theory) / theory
}

// EstimateVariance 估计序列方差并与理论值对照。
func EstimateVariance(samples []float64, theory float64) (estimate, relError float64) {
	v := Variance(samples)
	if theory == 0 {
		return v, math.NaN()
	}
	return v, math.Abs(v-theory) / theory
}
