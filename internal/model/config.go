// Package model 定义 CA-CFAR 核算的输入数据模型：检测配置
// （幅度序列、保护单元数、每侧参考单元数、虚警率 Pfa）以及
// 逐单元检测结果。本包只描述数据与校验规则，不做滑窗运算。
package model

import "fmt"

// DetectorConfig 是一次 CA-CFAR 核算的完整输入。
//
// Amplitude 是距离向幅度序列（平方律检波后的幅度，非负）。
// Guards 是 CUT 两侧不参与参考均值估计的保护单元个数（每侧）。
// Refs 是保护单元外侧各取的参考单元个数（每侧）。
// Pfa 是名义虚警概率，必须落在 (0,1) 开区间。
type DetectorConfig struct {
	Amplitude []float64
	Guards    int
	Refs      int
	Pfa       float64
}

// NewDetectorConfig 构造检测配置。构造本身不校验，校验交给 Validate。
func NewDetectorConfig(amplitude []float64, guards, refs int, pfa float64) *DetectorConfig {
	return &DetectorConfig{
		Amplitude: amplitude,
		Guards:    guards,
		Refs:      refs,
		Pfa:       pfa,
	}
}

// ReferenceCount 返回参考单元总数 N = 2×Refs（两侧各 Refs 个）。
func (c *DetectorConfig) ReferenceCount() int {
	return 2 * c.Refs
}

// WindowHalf 返回 CUT 单侧窗宽（保护单元 + 参考单元），即
// 一个 CUT 若要被充分支撑，两侧各需要 Guards+Refs 个单元。
func (c *DetectorConfig) WindowHalf() int {
	return c.Guards + c.Refs
}

// SequenceLength 返回幅度序列长度。
func (c *DetectorConfig) SequenceLength() int {
	return len(c.Amplitude)
}

// 描述参数的一句话，用于输出标题行。
func (c *DetectorConfig) Describe() string {
	return fmt.Sprintf(
		"序列长 %d  保护 %d  每侧参考 %d  N=%d  Pfa=%.6g",
		c.SequenceLength(), c.Guards, c.Refs, c.ReferenceCount(), c.Pfa,
	)
}
