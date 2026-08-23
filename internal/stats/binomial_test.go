package stats

import "testing"

func TestBinomialTail(t *testing.T) {
	if p := BinomialTailP(100, 0, 0.5); p != 1 {
		t.Errorf("BinomialTailP(100,0,0.5) = %v, 期望 1", p)
	}
	if p := BinomialTailP(100, 101, 0.5); p != 0 {
		t.Errorf("BinomialTailP(100,101,0.5) = %v, 期望 0", p)
	}
	p1 := BinomialTailP(10, 1, 0.5)
	p5 := BinomialTailP(10, 5, 0.5)
	if !(p5 < p1) {
		t.Errorf("尾概率未随 k 增大而下降: %v vs %v", p1, p5)
	}
	if p1 > 1 {
		t.Error("概率不应超过 1")
	}
}

func TestExactPfaBand(t *testing.T) {
	pfa := 1e-3
	lo1, hi1, err := ExactPfaBand(1000, pfa, 0.01)
	if err != nil {
		t.Fatalf("ExactPfaBand(1000) 报错: %v", err)
	}
	lo2, hi2, err := ExactPfaBand(8000, pfa, 0.01)
	if err != nil {
		t.Fatalf("ExactPfaBand(8000) 报错: %v", err)
	}
	relWidth1 := float64(hi1-lo1) / 1000
	relWidth2 := float64(hi2-lo2) / 8000
	if relWidth1 <= relWidth2 {
		t.Errorf("相对区间宽度未随 n 增大收窄: n=1000→%.6g, n=8000→%.6g",
			relWidth1, relWidth2)
	}
	exp := float64(1000) * pfa
	if exp < float64(lo1) || exp > float64(hi1) {
		t.Errorf("名义期望检出 %v 不在区间 [%d, %d]", exp, lo1, hi1)
	}
}

func TestExactPfaBandValidation(t *testing.T) {
	if _, _, err := ExactPfaBand(0, 1e-3, 0.01); err == nil {
		t.Error("n=0 应报错")
	}
	if _, _, err := ExactPfaBand(100, 1, 0.01); err == nil {
		t.Error("pfa=1 应报错")
	}
	if _, _, err := ExactPfaBand(100, 1e-3, 0); err == nil {
		t.Error("alpha=0 应报错")
	}
}

func TestLogComb(t *testing.T) {
	a := LogComb(10, 3)
	b := LogComb(10, 7)
	if a != b {
		t.Errorf("C(10,3) 与 C(10,7) 对数应相等: %v vs %v", a, b)
	}
	if LogComb(5, 0) != 0 || LogComb(5, 5) != 0 {
		t.Error("C(n,0) 与 C(n,n) 对数应为 0")
	}
}
