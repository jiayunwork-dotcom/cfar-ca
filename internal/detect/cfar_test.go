package detect

import (
	"math"
	"testing"

	"cfar-ca/internal/model"
)

func syntheticConfig(n, g, r, cut int, targetAmp float64, pfa float64) *model.DetectorConfig {
	amps := make([]float64, n)
	for i := range amps {
		amps[i] = 1.0
	}
	amps[cut] = targetAmp
	return model.NewDetectorConfig(amps, g, r, pfa)
}

func TestThresholdEqualsAlphaTimesReferenceMean(t *testing.T) {
	amps := []float64{
		0.5, 1.2, 0.9, 2.0, 1.7, 3.0, 0.8, 1.1, 2.3, 1.4,
		0.6, 1.8, 0.7, 2.5, 1.0, 0.3, 1.9, 2.2, 0.4, 1.6,
		0.2, 1.3, 0.5, 2.8, 1.1, 0.9, 1.7, 2.0, 0.6, 1.2,
	}
	g, r := 2, 4
	pfa := 1e-3
	cfg := model.NewDetectorConfig(amps, g, r, pfa)
	alpha, err := AlphaFor(pfa, r)
	if err != nil {
		t.Fatalf("AlphaFor 报错: %v", err)
	}
	res, err := Detect(cfg)
	if err != nil {
		t.Fatalf("Detect 报错: %v", err)
	}
	for i, c := range res.Cells {
		if !c.Valid {
			continue
		}
		left := amps[i-g-r : i-g]
		right := amps[i+g+1 : i+g+1+r]
		sum := 0.0
		for _, v := range left {
			sum += v
		}
		for _, v := range right {
			sum += v
		}
		mean := sum / float64(2*r)
		want := alpha * mean
		if math.Abs(c.Threshold-want) > 1e-9 {
			t.Errorf("CUT %d: 阈值 %v ≠ α×参考均值 %v", i, c.Threshold, want)
		}
	}
}

func TestGuardExclusionPreventsSelfRaising(t *testing.T) {
	cfg := syntheticConfig(101, 2, 8, 50, 40.0, 1e-3)
	res, err := Detect(cfg)
	if err != nil {
		t.Fatalf("Detect 报错: %v", err)
	}
	cut := res.Cell(50)
	if !cut.Valid {
		t.Fatal("CUT 50 应在有效区")
	}
	if !cut.Detected {
		t.Errorf("强目标未检出: 幅度 %v, 阈值 %v", cut.Amplitude, cut.Threshold)
	}
	alpha, _ := AlphaFor(1e-3, 8)
	want := alpha * 1.0
	if math.Abs(cut.Threshold-want) > 1e-9 {
		t.Errorf("目标处阈值 %v 被抬高，期望噪声水平 %v", cut.Threshold, want)
	}
	neighbor := res.Cell(52)
	if math.Abs(neighbor.Threshold-want) > 1e-9 {
		t.Errorf("CUT 52 阈值 %v 受目标影响，期望 %v", neighbor.Threshold, want)
	}
}

func TestEdgeCUTsMarkedInvalid(t *testing.T) {
	g, r := 2, 4
	cfg := syntheticConfig(30, g, r, 15, 40.0, 1e-3)
	res, err := Detect(cfg)
	if err != nil {
		t.Fatalf("Detect 报错: %v", err)
	}
	if got := res.InvalidCount(); got != 2*(g+r) {
		t.Errorf("无效 CUT 数 = %d, 期望 %d", got, 2*(g+r))
	}
	for i := 0; i < g+r; i++ {
		c := res.Cell(i)
		if c.Valid {
			t.Errorf("左端 CUT %d 不应有效", i)
		}
		if !math.IsNaN(c.Threshold) {
			t.Errorf("左端 CUT %d 阈值应为 NaN, 得到 %v", i, c.Threshold)
		}
		if c.Detected {
			t.Errorf("左端 CUT %d 不应检出", i)
		}
	}
	for i := 30 - (g + r); i < 30; i++ {
		if res.Cell(i).Valid {
			t.Errorf("右端 CUT %d 不应有效", i)
		}
	}
	if !PartialWindowPolicyOK(cfg) {
		t.Error("左端第一个有效 CUT 的左参考窗起点应为 0（无借半窗）")
	}
}

func TestCrossRulePfaReducesDetections(t *testing.T) {
	cfg := syntheticConfig(201, 2, 8, 100, 25.0, 1e-3)
	sweep, err := SweepPfa(cfg, []float64{1e-2, 1e-3, 1e-4})
	if err != nil {
		t.Fatalf("SweepPfa 报错: %v", err)
	}
	if len(sweep.Entries) != 3 {
		t.Fatalf("Entries 数 = %d, 期望 3", len(sweep.Entries))
	}
	for i := 1; i < len(sweep.Entries); i++ {
		prev, cur := sweep.Entries[i-1], sweep.Entries[i]
		if !(cur.Pfa > prev.Pfa) {
			t.Fatalf("扫描未按 Pfa 升序: %v >= %v", cur.Pfa, prev.Pfa)
		}
		if !(cur.Alpha < prev.Alpha) {
			t.Errorf("Pfa %v→%v 时 α 未下降: %v → %v",
				prev.Pfa, cur.Pfa, prev.Alpha, cur.Alpha)
		}
		if cur.Detected < prev.Detected {
			t.Errorf("Pfa 升高后检出数减少: %d → %d",
				prev.Detected, cur.Detected)
		}
	}
	rule, err := ComparePfa(cfg, 1e-3, 1e-4)
	if err != nil {
		t.Fatalf("ComparePfa 报错: %v", err)
	}
	if !rule.AlphaRises {
		t.Errorf("Pfa 降低后 α 未升高: %v → %v", rule.AlphaHigh, rule.AlphaLow)
	}
	if !rule.DetFalls {
		t.Errorf("Pfa 降低后检出数不减反增: %d → %d", rule.DetHigh, rule.DetLow)
	}
	last := sweep.Entries[len(sweep.Entries)-1]
	res, err := DetectWithAlpha(cfg, last.Alpha)
	if err != nil {
		t.Fatalf("DetectWithAlpha 报错: %v", err)
	}
	if !res.DetectedAt(100) {
		t.Error("低 Pfa 下强目标漏检：目标被自身抬阈或交叉规则反了")
	}
}

func TestCrossRuleAmplitudeRaisesMargin(t *testing.T) {
	cfg := syntheticConfig(201, 2, 8, 100, 25.0, 1e-3)
	m0, m1, ok := MarginRiseAt(cfg, 100, 25.0, 60.0)
	if !ok {
		t.Fatal("MarginRiseAt 无法计算（CUT 无效或阈值非正）")
	}
	if !(m1 > m0) {
		t.Errorf("目标幅度加大后裕量未升高: %v → %v", m0, m1)
	}
	if m0 <= 1 || m1 <= 1 {
		t.Errorf("裕量应 >1（越过阈值）: %v, %v", m0, m1)
	}
}
