package stats

import "math"

// CountDetected 统计布尔标志序列中 true 的个数。
func CountDetected(flags []bool) int {
	n := 0
	for _, f := range flags {
		if f {
			n++
		}
	}
	return n
}

// EmpiricalFalseAlarm 计算经验虚警率 = 检出数 / 有效单元数。
// 有效单元为 0 时返回 NaN。
func EmpiricalFalseAlarm(detected, valid int) float64 {
	if valid <= 0 {
		return nanStat()
	}
	return float64(detected) / float64(valid)
}

// DetectionRate 是 EmpiricalFalseAlarm 的别名（检测比例）。
func DetectionRate(detected, valid int) float64 {
	return EmpiricalFalseAlarm(detected, valid)
}

// nanStat 返回 NaN 占位。
func nanStat() float64 {
	return math.NaN()
}

// Band 描述经验虚警的合理波动带。
type Band struct {
	Pfa   float64
	Lo    float64
	Hi    float64
	Mean  float64
	Sigma float64
}

// BandFor 给定有效单元数 n 与名义 Pfa，用二项分布近似给出经验虚警
// 的合理波动带：均值 = p，σ = sqrt(p(1−p)/n)，带取 [p − kσ, p + kσ]。
//
// 该带用于「均匀噪声长序列，检测比例应落在 Pfa 的合理波动带」的
// 自动校验；k 越大越宽松。
func BandFor(n int, pfa float64, k float64) (Band, error) {
	if n <= 0 {
		return Band{}, &BadCountError{N: n}
	}
	if pfa <= 0 || pfa >= 1 {
		return Band{}, &BadRateError{Rate: pfa}
	}
	mean := pfa
	sigma := math.Sqrt(pfa * (1 - pfa) / float64(n))
	return Band{
		Pfa:   pfa,
		Lo:    mean - k*sigma,
		Hi:    mean + k*sigma,
		Mean:  mean,
		Sigma: sigma,
	}, nil
}

// InBand 判断经验值是否落在波动带内。带下界被钳制到 ≥0。
func InBand(rate float64, b Band) bool {
	if b.Lo < 0 {
		b.Lo = 0
	}
	return rate >= b.Lo && rate <= b.Hi
}

// CheckEmpiricalBand 把经验虚警率与名义 Pfa 的波动带一起返回，
// 便于调用方打印「是否落在带内」。
type BandCheck struct {
	Band      Band
	Empirical float64
	InBand    bool
}

// CheckBand 组合完成波动带校验。
func CheckBand(empirical float64, n int, pfa float64, k float64) (BandCheck, error) {
	b, err := BandFor(n, pfa, k)
	if err != nil {
		return BandCheck{}, err
	}
	return BandCheck{Band: b, Empirical: empirical, InBand: InBand(empirical, b)}, nil
}
