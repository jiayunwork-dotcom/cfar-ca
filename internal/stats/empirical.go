package stats

import "math"

func CountDetected(flags []bool) int {
	n := 0
	for _, f := range flags {
		if f {
			n++
		}
	}
	return n
}

func EmpiricalFalseAlarm(detected, valid int) float64 {
	if valid <= 0 {
		return nanStat()
	}
	return float64(detected) / float64(valid)
}

func DetectionRate(detected, valid int) float64 {
	return EmpiricalFalseAlarm(detected, valid)
}

func nanStat() float64 {
	return math.NaN()
}

type Band struct {
	Pfa   float64
	Lo    float64
	Hi    float64
	Mean  float64
	Sigma float64
}

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

func InBand(rate float64, b Band) bool {
	if b.Lo < 0 {
		b.Lo = 0
	}
	return rate >= b.Lo && rate <= b.Hi
}

type BandCheck struct {
	Band      Band
	Empirical float64
	InBand    bool
}

func CheckBand(empirical float64, n int, pfa float64, k float64) (BandCheck, error) {
	b, err := BandFor(n, pfa, k)
	if err != nil {
		return BandCheck{}, err
	}
	return BandCheck{Band: b, Empirical: empirical, InBand: InBand(empirical, b)}, nil
}
