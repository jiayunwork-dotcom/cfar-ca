package stats

import (
	"fmt"
	"math"
)

// EmptySeriesError 表示空序列。
type EmptySeriesError struct{}

func (e *EmptySeriesError) Error() string {
	return "序列为空，无法统计"
}

// BadMeanError 表示均值参数不合法。
type BadMeanError struct {
	Mean float64
}

func (e *BadMeanError) Error() string {
	return fmt.Sprintf("均值必须为正，得到 %v", e.Mean)
}

// BadRateError 表示概率参数不合法。
type BadRateError struct {
	Rate float64
}

func (e *BadRateError) Error() string {
	return fmt.Sprintf("概率必须在 (0,1) 区间，得到 %v", e.Rate)
}

// Exponential 生成 n 个均值为 mean 的指数分布样本（逆变换抽样）：
//
//	x = −mean × ln(1 − U),  U ~ Uniform(0,1)
//
// 该分布对应平方律检波后的功率噪声，是 CA-CFAR 教科书 α 公式
// 对应的背景模型。
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
		// 用 1−u 避免 u=0 导致 ln(0)。
		out[i] = -mean * math.Log1p(-u)
	}
	return out, nil
}

// BadCountError 表示抽样个数不合法。
type BadCountError struct {
	N int
}

func (e *BadCountError) Error() string {
	return fmt.Sprintf("抽样个数必须非负，得到 %d", e.N)
}

// ExponentialMeanVar 返回指数分布的理论均值与方差（检验生成器用）。
func ExponentialMeanVar(mean float64) (mu, variance float64) {
	return mean, mean * mean
}

// ExponentialPfa 给定阈值 T 与均值 μ，返回指数噪声的虚警概率
// P(x>T) = exp(−T/μ)。用于把名义 Pfa 换算成阈值对照。
func ExponentialPfa(threshold, mean float64) float64 {
	if threshold <= 0 || mean <= 0 {
		return math.NaN()
	}
	return math.Exp(-threshold / mean)
}

// ExponentialThreshold 给定 Pfa 与均值 μ，返回对应阈值 = −μ ln(Pfa)。
func ExponentialThreshold(pfa, mean float64) float64 {
	if pfa <= 0 || pfa >= 1 || mean <= 0 {
		return math.NaN()
	}
	return -mean * math.Log(pfa)
}
