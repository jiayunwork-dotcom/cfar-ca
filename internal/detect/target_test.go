package detect

import (
	"testing"

	"cfar-ca/internal/model"
)

func TestTargetCandidateLocalization(t *testing.T) {
	amps := make([]float64, 101)
	for i := range amps {
		amps[i] = 1.0
	}
	amps[50] = 30.0
	cfg := model.NewDetectorConfig(amps, 2, 8, 1e-3)
	res, err := Detect(cfg)
	if err != nil {
		t.Fatalf("Detect 报错: %v", err)
	}
	peaks := LocalPeaks(res)
	found := false
	for _, p := range peaks {
		if p.Index == 50 {
			found = true
			if !p.Peak {
				t.Error("目标 50 应被标记为局部峰值")
			}
		}
	}
	if !found {
		t.Errorf("目标未被定位为局部峰值，peaks=%v", peaks)
	}
	if TargetCount(res) == 0 {
		t.Error("TargetCount 应为至少 1")
	}
}

func TestStrongestCells(t *testing.T) {
	amps := make([]float64, 101)
	for i := range amps {
		amps[i] = 1.0
	}
	amps[50] = 30.0
	cfg := model.NewDetectorConfig(amps, 2, 8, 1e-3)
	res, err := Detect(cfg)
	if err != nil {
		t.Fatalf("Detect 报错: %v", err)
	}
	top := StrongestCells(res, 1)
	if len(top) != 1 {
		t.Fatalf("StrongestCells(1) = %d 个, 期望 1", len(top))
	}
	if top[0].Index != 50 {
		t.Errorf("最强单元 = %d, 期望目标 50", top[0].Index)
	}
	if top[0].Margin <= 1 {
		t.Errorf("目标裕量 %v 应 >1", top[0].Margin)
	}
	if got := StrongestCells(res, 0); got != nil {
		t.Error("k=0 应返回 nil")
	}
}

func TestEdgeValidNeighborSearch(t *testing.T) {
	cfg := syntheticConfig(30, 2, 4, 15, 40.0, 1e-3)
	res, err := Detect(cfg)
	if err != nil {
		t.Fatalf("Detect 报错: %v", err)
	}
	first := FirstValidIndex(cfg)
	if got := prevValid(res, first); got != -1 {
		t.Errorf("第一个有效 CUT 左侧应有有效单元, 得到 %d", got)
	}
	last := LastValidIndex(cfg)
	if got := nextValid(res, last); got != -1 {
		t.Errorf("最后一个有效 CUT 右侧应有有效单元, 得到 %d", got)
	}
}
