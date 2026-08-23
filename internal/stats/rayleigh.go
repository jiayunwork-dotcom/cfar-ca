package stats

import (
	"math"
)

// Rayleigh 生成 n 个尺度参数为 sigma 的瑞利包络样本：
//
//	x = sigma × sqrt(−2 ln(1 − U)),  U ~ Uniform(0,1)
//
// 用于生成幅度域的瑞利噪声序列（如未做平方律检波的线性检波输出）。
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

// RayleighMean 返回瑞利分布的理论均值 μ = σ√(π/2)。
func RayleighMean(sigma float64) float64 {
	return sigma * math.Sqrt(math.Pi/2.0)
}

// RayleighVariance 返回瑞利分布的理论方差 = (2−π/2)σ²。
func RayleighVariance(sigma float64) float64 {
	return (2.0 - math.Pi/2.0) * sigma * sigma
}

// RayleighCdf 返回瑞利分布 CDF P(x ≤ t) = 1 − exp(−t²/(2σ²))。
func RayleighCdf(t, sigma float64) float64 {
	if sigma <= 0 {
		return math.NaN()
	}
	if t <= 0 {
		return 0
	}
	return 1.0 - math.Exp(-t*t/(2.0*sigma*sigma))
}

// RayleighTailPfa 给定阈值 T 与尺度 σ，返回瑞利噪声虚警概率
// P(x>T) = exp(−T²/(2σ²))。注意：教科书 CA-CFAR 的 α 公式对应
// 指数背景，瑞利幅度下经验虚警会显著低于名义 Pfa——这正是把
// 本仓序列约定为「平方律检波后幅度（指数背景）」的原因。
func RayleighTailPfa(threshold, sigma float64) float64 {
	if threshold <= 0 || sigma <= 0 {
		return math.NaN()
	}
	return math.Exp(-threshold * threshold / (2.0 * sigma * sigma))
}

// SigmaFromMean 由瑞利均值反推尺度参数 σ = μ / √(π/2)。
func SigmaFromMean(mean float64) float64 {
	if mean <= 0 {
		return math.NaN()
	}
	return mean / math.Sqrt(math.Pi/2.0)
}
