package io

import (
	"strings"
	"testing"

	"cfar-ca/internal/detect"
	"cfar-ca/internal/model"
)

func TestFormatDetectOutput(t *testing.T) {
	amps := make([]float64, 30)
	for i := range amps {
		amps[i] = 1.0
	}
	amps[15] = 30.0
	cfg := model.NewDetectorConfig(amps, 2, 4, 1e-3)
	res, err := detect.Detect(cfg)
	if err != nil {
		t.Fatalf("Detect 报错: %v", err)
	}
	out := FormatDetect(res)
	for _, want := range []string{
		"CA-CFAR 检测结果",
		"α=",
		"经验虚警",
		"检出索引: [15]",
		"无效(窗不足)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("输出缺少 %q", want)
		}
	}
}

func TestFormatAlphaOutput(t *testing.T) {
	out := FormatAlpha(1e-3, 8, 8.638824416951874)
	for _, want := range []string{"N=16", "8.63882", "Pfa=0.001"} {
		if !strings.Contains(out, want) {
			t.Errorf("输出缺少 %q", want)
		}
	}
}

func TestFormatSweepOutput(t *testing.T) {
	cfg := model.NewDetectorConfig(
		[]float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 2, 3, 1e-3)
	sweep, err := detect.SweepPfa(cfg, []float64{1e-3, 1e-4})
	if err != nil {
		t.Fatalf("SweepPfa 报错: %v", err)
	}
	out := FormatSweep(sweep)
	for _, want := range []string{"Pfa", "α", "检出"} {
		if !strings.Contains(out, want) {
			t.Errorf("输出缺少 %q", want)
		}
	}
}

func TestFormatCrossRuleOutput(t *testing.T) {
	cfg := model.NewDetectorConfig(
		[]float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 2, 3, 1e-3)
	rule, err := detect.ComparePfa(cfg, 1e-3, 1e-4)
	if err != nil {
		t.Fatalf("ComparePfa 报错: %v", err)
	}
	out := FormatCrossRule(rule)
	for _, want := range []string{"α 升高: true", "检出减少: true"} {
		if !strings.Contains(out, want) {
			t.Errorf("输出缺少 %q", want)
		}
	}
}
