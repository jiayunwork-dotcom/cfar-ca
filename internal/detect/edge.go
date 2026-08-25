package detect

import "cfar-ca/internal/model"

const EdgePolicy = "边缘 CUT 参考不足 → 标记无效，不借用半窗、不补零"

type EdgeDiagnostics struct {
	LeftInvalid  int
	RightInvalid int
	TotalInvalid int
	FirstValid   int
	LastValid    int
}

func DiagnoseEdges(cfg *model.DetectorConfig) EdgeDiagnostics {
	d := EdgeDiagnostics{
		FirstValid: FirstValidIndex(cfg),
		LastValid:  LastValidIndex(cfg),
	}
	for i := 0; i < cfg.SequenceLength(); i++ {
		if HasFullSupport(cfg, i) {
			continue
		}
		d.TotalInvalid++
		if i < cfg.Guards+cfg.Refs {
			d.LeftInvalid++
		} else {
			d.RightInvalid++
		}
	}
	return d
}

func LeftInvalidCount(cfg *model.DetectorConfig) int {
	return DiagnoseEdges(cfg).LeftInvalid
}

func RightInvalidCount(cfg *model.DetectorConfig) int {
	return DiagnoseEdges(cfg).RightInvalid
}

func TotalInvalidCount(cfg *model.DetectorConfig) int {
	if cfg == nil {
		return 0
	}
	return 2 * cfg.WindowHalf()
}

func PartialWindowPolicyOK(cfg *model.DetectorConfig) bool {
	first := FirstValidIndex(cfg)
	if first < 0 {
		return false
	}
	return first-cfg.Guards-cfg.Refs == 0
}

func DescribeEdgePolicy() string {
	return EdgePolicy
}
