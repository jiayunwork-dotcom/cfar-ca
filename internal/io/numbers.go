package io

import (
	"fmt"
	"math"
	"strconv"
)

func ParsePfa(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("Pfa 不是有效数字 %q: %w", s, err)
	}
	if math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 || v >= 1 {
		return 0, fmt.Errorf("Pfa 必须在 (0,1) 开区间，得到 %v", v)
	}
	return v, nil
}

func ParsePfaList(s string) ([]float64, error) {
	fields := splitFields(s)
	if len(fields) == 0 {
		return nil, fmt.Errorf("Pfa 列表为空")
	}
	out := make([]float64, 0, len(fields))
	for _, f := range fields {
		v, err := ParsePfa(f)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func ParsePositiveInt(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("不是有效整数 %q: %w", s, err)
	}
	if v <= 0 {
		return 0, fmt.Errorf("必须为正整数，得到 %d", v)
	}
	return v, nil
}

func ParseNonNegativeInt(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("不是有效整数 %q: %w", s, err)
	}
	if v < 0 {
		return 0, fmt.Errorf("必须为非负整数，得到 %d", v)
	}
	return v, nil
}

func ParseFloat(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("不是有效数字 %q: %w", s, err)
	}
	return v, nil
}

func splitFields(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == ',' || r == ' ' || r == '\t' || r == '\n' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

func FormatFloatList(vals []float64) string {
	out := ""
	for i, v := range vals {
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%g", v)
	}
	return out
}
