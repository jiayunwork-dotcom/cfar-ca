// Package detect 实现单元平均 CA-CFAR（Cell-Averaging Constant False
// Alarm Rate）滑窗阈值与检测判定。
//
// 核心公式（指数背景的教科书形式）：
//
//	α = N ( Pfa^{-1/N} − 1 )
//
// 其中 N = 2×Refs 是参考单元总数，Pfa 是名义虚警概率。
// 每个 CUT 的阈值 = α × 参考窗均值；CUT 与两侧保护单元一律不进参考均值。
package detect

import (
	"fmt"
	"math"
)

// ErrInvalidAlphaParameter 表示 α 公式的输入越界（Pfa 越界或参考单元数非正）。
type ErrInvalidAlphaParameter struct {
	Pfa  float64
	Refs int
}

func (e *ErrInvalidAlphaParameter) Error() string {
	return fmt.Sprintf(
		"α 公式输入不合法: Pfa=%v（需 (0,1)）, refs=%d（需 >0）", e.Pfa, e.Refs,
	)
}

// AlphaFor 计算每侧 Refs 个参考单元时的放大系数 α = N(Pfa^{-1/N}−1)。
//
// 校验：
//   - Refs 必须为正（窗长 0 拒绝）；
//   - Pfa 必须落在 (0,1) 开区间。
func AlphaFor(pfa float64, refs int) (float64, error) {
	if refs <= 0 {
		return 0, stringifyAlphaErr(&ErrInvalidAlphaParameter{Pfa: pfa, Refs: refs})
	}
	a, err := AlphaForN(pfa, 2*refs)
	if err != nil {
		return 0, stringifyAlphaErr(err)
	}
	return a, nil
}

// AlphaForN 直接以参考单元总数 N 计算放大系数 α = N(Pfa^{-1/N}−1)。
// N 必须为正偶数（两侧各一半）之外不额外限制；Pfa 越界报错。
func AlphaForN(pfa float64, n int) (float64, error) {
	if n <= 0 {
		return 0, stringifyAlphaErr(&ErrInvalidAlphaParameter{Pfa: pfa, Refs: n / 2})
	}
	if math.IsNaN(pfa) || math.IsInf(pfa, 0) || pfa <= 0 || pfa >= 1 {
		return 0, stringifyAlphaErr(&ErrInvalidAlphaParameter{Pfa: pfa, Refs: n / 2})
	}
	exponent := 1.0 / float64(n)
	return float64(n) * (math.Pow(pfa, -exponent) - 1.0), nil
}

// AlphaIncreasesWhenPfaDrops 判断参考单元数固定时，Pfa 降低是否让 α 升高
// （交叉规则：只把 Pfa 降一个数量级，α 必须升高）。返回比较结果与
// 两个 α 值，供测试与 compare 子命令复用。
func AlphaIncreasesWhenPfaDrops(pfaHigh, pfaLow float64, refs int) (bool, float64, float64, error) {
	aHigh, err := AlphaFor(pfaHigh, refs)
	if err != nil {
		return false, 0, 0, err
	}
	aLow, err := AlphaFor(pfaLow, refs)
	if err != nil {
		return false, 0, 0, err
	}
	return aLow > aHigh, aHigh, aLow, nil
}

// AlphaTendsToLnInvPfa 返回参考单元总数趋于无穷时 α 的极限 ln(1/Pfa)，
// 用于校验有限 N 下 α 随 N 增大而单调下降并逼近该极限的行为。
func AlphaTendsToLnInvPfa(pfa float64) (float64, error) {
	if pfa <= 0 || pfa >= 1 {
		return 0, &ErrInvalidAlphaParameter{Pfa: pfa}
	}
	return math.Log(1.0 / pfa), nil
}
