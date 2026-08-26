package detect

import (
	"fmt"
	"math"
)

type ErrInvalidAlphaParameter struct {
	Pfa  float64
	Refs int
}

func (e *ErrInvalidAlphaParameter) Error() string {
	return fmt.Sprintf(
		"α 公式输入不合法: Pfa=%v（需 (0,1)）, refs=%d（需 >0）", e.Pfa, e.Refs,
	)
}

func AlphaFor(pfa float64, refs int) (float64, error) {
	if refs <= 0 {
		return 0, &ErrInvalidAlphaParameter{Pfa: pfa, Refs: refs}
	}
	return AlphaForN(pfa, 2*refs)
}

func AlphaForN(pfa float64, n int) (float64, error) {
	if n <= 0 {
		return 0, &ErrInvalidAlphaParameter{Pfa: pfa, Refs: n / 2}
	}
	if math.IsNaN(pfa) || math.IsInf(pfa, 0) || pfa <= 0 || pfa >= 1 {
		return 0, &ErrInvalidAlphaParameter{Pfa: pfa, Refs: n / 2}
	}
	exponent := 1.0 / float64(n)
	return float64(n) * (math.Pow(pfa, -exponent) - 1.0), nil
}

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

func AlphaTendsToLnInvPfa(pfa float64) (float64, error) {
	if pfa <= 0 || pfa >= 1 {
		return 0, &ErrInvalidAlphaParameter{Pfa: pfa}
	}
	return math.Log(1.0 / pfa), nil
}
