package model

import (
	"fmt"
	"math"
	"strings"
)

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

func validatePfa(pfa float64) error {
	if math.IsNaN(pfa) || math.IsInf(pfa, 0) || pfa <= 0 || pfa >= 1 {
		return &ConfigError{
			Reason: fmt.Sprintf("虚警率 Pfa 必须在 (0,1) 开区间，得到 %v", pfa),
		}
	}
	return nil
}

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

func firstBadAmplitude(samples []float64) int {
	for i, v := range samples {
		if math.IsNaN(v) || v < 0 {
			return i
		}
	}
	return -1
}

type ConfigError struct {
	Reason string
}

func (e *ConfigError) Error() string {
	return "配置不合法: " + e.Reason
}

func IsConfigError(err error) bool {
	if err == nil {
		return false
	}
	var ce *ConfigError
	return asConfigError(err, &ce)
}

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

func (c *DetectorConfig) Issues() []string {
	if c == nil {
		return []string{"配置为空"}
	}
	return collectIssues(c)
}

func (c *DetectorConfig) DescribeIssues() string {
	issues := c.Issues()
	if len(issues) == 0 {
		return ""
	}
	return strings.Join(issues, "; ")
}
