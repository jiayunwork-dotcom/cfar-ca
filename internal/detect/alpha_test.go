package detect

import (
	"math"
	"testing"
)

func TestAlphaFormulaMatchesTextbook(t *testing.T) {
	a, err := AlphaFor(1e-3, 8)
	if err != nil {
		t.Fatalf("AlphaFor 报错: %v", err)
	}
	// 16×(1000^{1/16}−1) ≈ 8.6388
	if math.Abs(a-8.638824416951874) > 1e-9 {
		t.Errorf("α = %v, 期望 8.638824416951874（N=16, Pfa=1e-3）", a)
	}

	ok, aHigh, aLow, err := AlphaIncreasesWhenPfaDrops(1e-3, 1e-4, 8)
	if err != nil {
		t.Fatalf("AlphaIncreasesWhenPfaDrops 报错: %v", err)
	}
	if !ok {
		t.Errorf("Pfa 降低后 α 未升高: 1e-3→%v, 1e-4→%v", aHigh, aLow)
	}

	a8, _ := AlphaFor(1e-3, 4)
	a16, _ := AlphaFor(1e-3, 8)
	a32, _ := AlphaFor(1e-3, 16)
	limit, _ := AlphaTendsToLnInvPfa(1e-3)
	if !(a8 > a16 && a16 > a32) {
		t.Errorf("α 未随 N 增大单调下降: N=8→%v, N=16→%v, N=32→%v", a8, a16, a32)
	}
	if a32 < limit {
		t.Errorf("α(%v) 低于极限 %v，违反下界", a32, limit)
	}
}

func TestAlphaRejectsBadInputs(t *testing.T) {
	badPfas := []float64{0, 1, -0.1, 1.2, math.NaN()}
	for _, p := range badPfas {
		if _, err := AlphaFor(p, 8); err == nil {
			t.Errorf("AlphaFor(%v, 8) 应报错", p)
		}
	}
	if _, err := AlphaFor(1e-3, 0); err == nil {
		t.Error("AlphaFor(1e-3, 0) 窗长 0 应报错")
	}
	if _, err := AlphaForN(1e-3, 0); err == nil {
		t.Error("AlphaForN(1e-3, 0) 应报错")
	}
}

func TestAlphaStrictlyPositive(t *testing.T) {
	for _, pfa := range []float64{1e-2, 1e-3, 1e-4, 1e-5, 1e-6} {
		for _, refs := range []int{2, 4, 8, 16, 32} {
			a, err := AlphaFor(pfa, refs)
			if err != nil {
				t.Fatalf("AlphaFor(%v,%d) 报错: %v", pfa, refs, err)
			}
			if a <= 0 || math.IsInf(a, 0) || math.IsNaN(a) {
				t.Errorf("α 不合法: pfa=%v refs=%d → %v", pfa, refs, a)
			}
		}
	}
}

func TestAlphaForTypedError(t *testing.T) {
	_, err := AlphaFor(0, 8)
	if err == nil {
		t.Fatal("AlphaFor(0, 8) 应报错")
	}
	if _, ok := err.(*ErrInvalidAlphaParameter); !ok {
		t.Errorf("非法 Pfa 应得到 *ErrInvalidAlphaParameter, 得到 %T", err)
	}
}
