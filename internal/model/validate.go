package model

import (
	"fmt"
	"math"
	"strings"
)

// Validate 校验检测配置。任一不变量被破坏就返回描述清晰的错误，
// 由调用方决定如何上报（CLI 打印到 stderr 并非零退出）。
//
// 校验规则：
//   - Pfa 必须严格落在 (0,1) 开区间，0、1、负数、大于 1 全部拒绝；
//   - Refs（每侧参考单元数）必须为正，窗长 0 拒绝；
//   - Guards 必须非负，且严格小于 Refs（保护≥参考拒绝）；
//   - 幅度序列必须非空且不含负值；
//   - 参考窗不得伸出序列两端：序列长度必须大于 2×(Guards+Refs)，
//     否则连一个充分支撑的 CUT 都没有。
func Validate(c *DetectorConfig) error {
	if c == nil {
		return &ConfigError{Reason: "配置为空"}
	}
	if err := validatePfa(c.Pfa); err != nil {
		return err
	}
	if err := validateWindow(c.Guards, c.Refs); err != nil {
		return err
	}
	if err := validateAmplitudes(c.Amplitude); err != nil {
		return err
	}
	if len(c.Amplitude) <= 2*c.WindowHalf() {
		return &ConfigError{
			Reason: fmt.Sprintf(
				"参考窗伸出序列：序列长 %d 需要 > 2×(保护%d+参考%d)=%d",
				len(c.Amplitude), c.Guards, c.Refs, 2*c.WindowHalf(),
			),
		}
	}
	return nil
}

// validatePfa 检查虚警率落在 (0,1) 开区间。
func validatePfa(pfa float64) error {
	if math.IsNaN(pfa) || math.IsInf(pfa, 0) || pfa <= 0 || pfa >= 1 {
		return &ConfigError{
			Reason: fmt.Sprintf("虚警率 Pfa 必须在 (0,1) 开区间，得到 %v", pfa),
		}
	}
	return nil
}

// validateWindow 检查参考窗参数：窗长非零、保护单元小于参考单元。
func validateWindow(guards, refs int) error {
	if refs <= 0 {
		return &ConfigError{
			Reason: fmt.Sprintf("参考窗长必须为正（每侧参考单元数 > 0），得到 %d", refs),
		}
	}
	if guards < 0 {
		return &ConfigError{
			Reason: fmt.Sprintf("保护单元数必须非负，得到 %d", guards),
		}
	}
	if guards >= refs {
		return &ConfigError{
			Reason: fmt.Sprintf(
				"保护单元数必须小于每侧参考单元数：guards=%d, refs=%d",
				guards, refs,
			),
		}
	}
	return nil
}

// validateAmplitudes 检查幅度序列非空且非负。
func validateAmplitudes(samples []float64) error {
	if len(samples) == 0 {
		return &ConfigError{Reason: "幅度序列为空"}
	}
	bad := firstBadAmplitude(samples)
	if bad >= 0 {
		return &ConfigError{
			Reason: fmt.Sprintf(
				"幅度序列第 %d 个单元为负：%.6g（幅度不允许为负）", bad, samples[bad],
			),
		}
	}
	return nil
}

// firstBadAmplitude 返回第一个负值下标；全为非负返回 -1。
func firstBadAmplitude(samples []float64) int {
	for i, v := range samples {
		if math.IsNaN(v) || v < 0 {
			return i
		}
	}
	return -1
}

// ConfigError 描述检测配置不合规的原因。
type ConfigError struct {
	Reason string
}

func (e *ConfigError) Error() string {
	return "配置不合法: " + e.Reason
}

// IsConfigError 判断 err 是否为配置校验错误（供调用方区分参数错误
// 与运行时错误的场景；CLI 统一非零退出）。
func IsConfigError(err error) bool {
	if err == nil {
		return false
	}
	var ce *ConfigError
	return asConfigError(err, &ce)
}

// asConfigError 支持错误链中的 *ConfigError。
func asConfigError(err error, target **ConfigError) bool {
	for err != nil {
		if ce, ok := err.(*ConfigError); ok {
			*target = ce
			return true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}

// 收集全部校验问题（调试/测试用）。
func collectIssues(c *DetectorConfig) []string {
	var issues []string
	if c == nil {
		return []string{"配置为空"}
	}
	if err := validatePfa(c.Pfa); err != nil {
		issues = append(issues, err.Error())
	}
	if err := validateWindow(c.Guards, c.Refs); err != nil {
		issues = append(issues, err.Error())
	}
	if err := validateAmplitudes(c.Amplitude); err != nil {
		issues = append(issues, err.Error())
	}
	if len(issues) == 0 && len(c.Amplitude) <= 2*c.WindowHalf() {
		issues = append(issues, "参考窗伸出序列")
	}
	return issues
}

// Issues 返回配置的全部校验问题；合规时为空切片。
func (c *DetectorConfig) Issues() []string {
	if c == nil {
		return []string{"配置为空"}
	}
	return collectIssues(c)
}

// DescribeIssues 把所有问题拼成一行，用于错误汇总输出。
func (c *DetectorConfig) DescribeIssues() string {
	issues := c.Issues()
	if len(issues) == 0 {
		return ""
	}
	return strings.Join(issues, "; ")
}
