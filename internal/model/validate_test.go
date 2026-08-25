package model

import "testing"

func TestRejectInvalidParameters(t *testing.T) {
	validAmps := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	cases := []struct {
		name   string
		amps   []float64
		guards int
		refs   int
		pfa    float64
	}{
		{name: "Pfa=0", amps: validAmps, guards: 1, refs: 2, pfa: 0},
		{name: "Pfa=1", amps: validAmps, guards: 1, refs: 2, pfa: 1},
		{name: "Pfa 负", amps: validAmps, guards: 1, refs: 2, pfa: -0.5},
		{name: "Pfa 大于1", amps: validAmps, guards: 1, refs: 2, pfa: 1.5},
		{name: "窗长0", amps: validAmps, guards: 1, refs: 0, pfa: 1e-3},
		{name: "保护等于参考", amps: validAmps, guards: 2, refs: 2, pfa: 1e-3},
		{name: "保护大于参考", amps: validAmps, guards: 3, refs: 2, pfa: 1e-3},
		{name: "负幅度", amps: []float64{1, -2, 3}, guards: 1, refs: 2, pfa: 1e-3},
		{name: "窗伸出序列", amps: []float64{1, 2, 3, 4, 5}, guards: 1, refs: 2, pfa: 1e-3},
		{name: "空序列", amps: nil, guards: 1, refs: 2, pfa: 1e-3},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := NewDetectorConfig(tc.amps, tc.guards, tc.refs, tc.pfa)
			if err := Validate(cfg); err == nil {
				t.Errorf("期望错误，但 Validate 通过: %+v", tc)
			}
		})
	}
}

func TestAcceptValidConfig(t *testing.T) {
	cfg := NewDetectorConfig(
		[]float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 2, 4, 1e-3)
	if err := Validate(cfg); err != nil {
		t.Errorf("合法配置被拒绝: %v", err)
	}
	if cfg.ReferenceCount() != 8 {
		t.Errorf("ReferenceCount() = %d, 期望 8", cfg.ReferenceCount())
	}
	if cfg.WindowHalf() != 6 {
		t.Errorf("WindowHalf() = %d, 期望 6", cfg.WindowHalf())
	}
}

func TestIssuesCollect(t *testing.T) {
	cfg := NewDetectorConfig([]float64{1, 2, 3}, 2, 2, 0)
	issues := cfg.Issues()
	if len(issues) == 0 {
		t.Fatal("期望 Issues 非空")
	}
	if cfg.DescribeIssues() == "" {
		t.Error("DescribeIssues 应为空串以外的内容")
	}
}
