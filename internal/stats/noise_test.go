package stats

import (
	"math"
	"testing"
)

func TestExponentialNoiseMoments(t *testing.T) {
	n := 20000
	mean := 1.5
	s1, err := Exponential(NewSource(42), n, mean)
	if err != nil {
		t.Fatalf("Exponential 报错: %v", err)
	}
	s2, err := Exponential(NewSource(42), n, mean)
	if err != nil {
		t.Fatalf("Exponential(同一种子) 报错: %v", err)
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			t.Fatalf("同一种子生成不一致（i=%d）: %v vs %v", i, s1[i], s2[i])
		}
	}
	mu, variance := ExponentialMeanVar(mean)
	estMean, relErr := EstimateMean(s1, mu)
	if relErr > 0.03 {
		t.Errorf("样本均值 %v 偏离理论 %v 超过 3%%", estMean, mu)
	}
	estVar, relErrVar := EstimateVariance(s1, variance)
	if relErrVar > 0.06 {
		t.Errorf("样本方差 %v 偏离理论 %v 超过 6%%", estVar, variance)
	}
	if Min(s1) < 0 {
		t.Error("指数噪声出现负样本")
	}
}

func TestRayleighMoments(t *testing.T) {
	n := 20000
	sigma := 2.0
	s, err := Rayleigh(NewSource(7), n, sigma)
	if err != nil {
		t.Fatalf("Rayleigh 报错: %v", err)
	}
	mu := RayleighMean(sigma)
	variance := RayleighVariance(sigma)
	estMean, relErr := EstimateMean(s, mu)
	if relErr > 0.03 {
		t.Errorf("瑞利样本均值 %v 偏离理论 %v 超过 3%%", estMean, mu)
	}
	estVar, relErrVar := EstimateVariance(s, variance)
	if relErrVar > 0.06 {
		t.Errorf("瑞利样本方差 %v 偏离理论 %v 超过 6%%", estVar, variance)
	}
	if Min(s) < 0 {
		t.Error("瑞利噪声出现负样本")
	}
}

func TestExponentialTailMath(t *testing.T) {
	mean := 1.0
	pfa := 1e-3
	t0 := ExponentialThreshold(pfa, mean)
	p := ExponentialPfa(t0, mean)
	if math.Abs(p-pfa) > 1e-12 {
		t.Errorf("P(x>%.6g) = %.6g, 期望 %.6g", t0, p, pfa)
	}
}

func TestScaleToMean(t *testing.T) {
	s, err := ScaleToMean([]float64{1, 2, 3}, 6.0)
	if err != nil {
		t.Fatalf("ScaleToMean 报错: %v", err)
	}
	if Mean(s) != 6.0 {
		t.Errorf("缩放后均值 = %v, 期望 6", Mean(s))
	}
	if _, err := ScaleToMean([]float64{1, 2, 3}, -1); err == nil {
		t.Error("负均值应报错")
	}
	if _, err := ScaleToMean(nil, 1); err == nil {
		t.Error("空序列应报错")
	}
}
