package detect

import "cfar-ca/internal/model"

// EdgePolicy 描述边缘 CUT 的处理约定。实现上只做「无完整参考窗就标
// 无效」，绝不借半窗、不补零、不把不足的参考当全窗参与统计。
const EdgePolicy = "边缘 CUT 参考不足 → 标记无效，不借用半窗、不补零"

// EdgeDiagnostics 描述序列边缘的无效 CUT 分布。
type EdgeDiagnostics struct {
	LeftInvalid  int // 左端因参考不足而无效的 CUT 数
	RightInvalid int // 右端因参考不足而无效的 CUT 数
	TotalInvalid int
	FirstValid   int
	LastValid    int
}

// DiagnoseEdges 统计无效 CUT 的分布。没有有效 CUT 时 FirstValid/LastValid
// 为 -1（这种情况配置校验已拒绝，正常不会出现）。
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

// LeftInvalidCount 返回左端无效 CUT 数。
func LeftInvalidCount(cfg *model.DetectorConfig) int {
	return DiagnoseEdges(cfg).LeftInvalid
}

// RightInvalidCount 返回右端无效 CUT 数。
func RightInvalidCount(cfg *model.DetectorConfig) int {
	return DiagnoseEdges(cfg).RightInvalid
}

// TotalInvalidCount 返回无效 CUT 总数 = 2×(Guards+Refs)。
func TotalInvalidCount(cfg *model.DetectorConfig) int {
	if cfg == nil {
		return 0
	}
	return 2 * cfg.WindowHalf()
}

// PartialWindowPolicyOK 校验器：左端或右端若出现「部分参考」被当作有效，
// 返回 false（该情形永不满足——实现不借用半窗）。
func PartialWindowPolicyOK(cfg *model.DetectorConfig) bool {
	first := FirstValidIndex(cfg)
	if first < 0 {
		return false
	}
	// 左端第一个有效 CUT 必须恰好拥有完整的左窗：
	// 它的左参考窗起点必须是 0。
	return first-cfg.Guards-cfg.Refs == 0
}

// DescribeEdgePolicy 返回边缘策略的一句话说明，供 CLI 帮助引用。
func DescribeEdgePolicy() string {
	return EdgePolicy
}
