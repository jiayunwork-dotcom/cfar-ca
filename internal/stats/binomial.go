package stats

import (
	"math"
)

func LogComb(n, k int) float64 {
	if k < 0 || k > n {
		return math.Inf(-1)
	}
	if k == 0 || k == n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	return lgamma(float64(n)+1) - lgamma(float64(k)+1) - lgamma(float64(n-k)+1)
}

func lgamma(x float64) float64 {
	v, _ := math.Lgamma(x)
	return v
}

func BinomialTailP(n, k int, p float64) float64 {
	if k <= 0 {
		return 1
	}
	if k > n {
		return 0
	}
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		if k <= n {
			return 1
		}
		return 0
	}
	logP := math.Log(p)
	logQ := math.Log1p(-p)
	sum := 0.0
	for j := k; j <= n; j++ {
		sum += math.Exp(LogComb(n, j) + float64(j)*logP + float64(n-j)*logQ)
	}
	if sum > 1 {
		return 1
	}
	return sum
}

func ExactPfaBand(n int, pfa float64, alpha float64) (lo, hi int, err error) {
	if n <= 0 {
		return 0, 0, &BadCountError{N: n}
	}
	if pfa <= 0 || pfa >= 1 {
		return 0, 0, &BadRateError{Rate: pfa}
	}
	if alpha <= 0 || alpha >= 1 {
		return 0, 0, &BadRateError{Rate: alpha}
	}
	cdf := 0.0
	lo = 0
	for k := 0; k <= n; k++ {
		cdf += binomialPMF(n, k, pfa)
		if cdf > alpha/2 {
			lo = k
			break
		}
	}
	cdf = 0.0
	hi = n
	for k := 0; k <= n; k++ {
		cdf += binomialPMF(n, k, pfa)
		if cdf >= 1-alpha/2 {
			hi = k
			break
		}
	}
	return lo, hi, nil
}

func binomialPMF(n, k int, p float64) float64 {
	if k < 0 || k > n {
		return 0
	}
	return math.Exp(LogComb(n, k) + float64(k)*math.Log(p) +
		float64(n-k)*math.Log1p(-p))
}

func InExactBand(detected, lo, hi int) bool {
	return detected >= lo && detected <= hi
}
