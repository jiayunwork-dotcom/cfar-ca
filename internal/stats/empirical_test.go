package stats

import (
	"testing"

	"cfar-ca/internal/detect"
	"cfar-ca/internal/model"
)

func TestUniformNoiseFalseAlarmWithinBand(t *testing.T) {
	const n = 8000
	guards, refs := 2, 16
	pfa := 1e-3

	noise, err := Exponential(NewSource(7), n, 1.0)
	if err != nil {
		t.Fatalf("生成噪声报错: %v", err)
	}
	cfg := model.NewDetectorConfig(noise, guards, refs, pfa)
	res, err := detect.Detect(cfg)
	if err != nil {
		t.Fatalf("Detect 报错: %v", err)
	}
	st := res.Stats()
	if st.ValidCells == 0 {
		t.Fatal("无有效 CUT")
	}
	band, err := BandFor(st.ValidCells, pfa, 3)
	if err != nil {
		t.Fatalf("BandFor 报错: %v", err)
	}
	if !InBand(st.EmpiricalPfa, band) {
		t.Errorf("经验虚警 %v 不在波动带 [%v, %v] 内（Pfa=%v, 有效 CUT=%d, 检出=%d）",
			st.EmpiricalPfa, band.Lo, band.Hi, pfa, st.ValidCells, st.DetectedCells)
	}
	alpha, _ := detect.AlphaFor(pfa, refs)
	if res.Alpha != alpha {
		t.Errorf("结果 α=%v 与教科书 α=%v 不一致", res.Alpha, alpha)
	}
}

func TestBandForValidation(t *testing.T) {
	if _, err := BandFor(0, 1e-3, 3); err == nil {
		t.Error("n=0 应报错")
	}
	if _, err := BandFor(100, 0, 3); err == nil {
		t.Error("pfa=0 应报错")
	}
	if _, err := BandFor(100, 1, 3); err == nil {
		t.Error("pfa=1 应报错")
	}
	if _, err := BandFor(-5, 1e-3, 3); err == nil {
		t.Error("n<0 应报错")
	}
}

func TestEmpiricalFalseAlarmMath(t *testing.T) {
	if got := EmpiricalFalseAlarm(2, 100); got != 0.02 {
		t.Errorf("EmpiricalFalseAlarm(2,100) = %v, 期望 0.02", got)
	}
	if got := EmpiricalFalseAlarm(0, 100); got != 0 {
		t.Errorf("EmpiricalFalseAlarm(0,100) = %v, 期望 0", got)
	}
	if !isNaN(EmpiricalFalseAlarm(0, 0)) {
		t.Error("有效单元为 0 时应返回 NaN")
	}
	if got := CountDetected([]bool{true, false, true}); got != 2 {
		t.Errorf("CountDetected = %d, 期望 2", got)
	}
}

func isNaN(v float64) bool {
	return v != v
}
