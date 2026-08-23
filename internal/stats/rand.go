// Package stats 提供均匀噪声场景的生成与经验虚警统计：
// 指数/瑞利包络噪声抽样、经验虚警率计算、以及名义 Pfa 的
// 合理波动带（二项分布均值 ± kσ）。
package stats

import "math/rand"

// NewSource 以固定种子构造确定性随机源。同一种子生成的序列
// 完全可复现——交叉规则与虚警带测试依赖这一点。
func NewSource(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

// randSource 是生成器依赖的最小随机源接口（*rand.Rand 满足）。
type randSource interface {
	Float64() float64
}

// Uniform 生成 n 个 (0,1) 均匀随机数。
func Uniform(rng *rand.Rand, n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = rng.Float64()
	}
	return out
}

// ScaleToMean 把已生成序列线性缩放为指定均值（正均值要求非空序列）。
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
