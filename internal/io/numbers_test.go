package io

import (
	"strings"
	"testing"
)

func TestParsePfa(t *testing.T) {
	if v, err := ParsePfa("1e-3"); err != nil || v != 1e-3 {
		t.Errorf("ParsePfa(1e-3) = %v, %v", v, err)
	}
	for _, bad := range []string{"0", "1", "-0.5", "1.5", "abc", "nan", "inf"} {
		if _, err := ParsePfa(bad); err == nil {
			t.Errorf("ParsePfa(%q) 应报错", bad)
		}
	}
}

func TestParsePfaList(t *testing.T) {
	vals, err := ParsePfaList("1e-2 1e-3,1e-4")
	if err != nil {
		t.Fatalf("ParsePfaList 报错: %v", err)
	}
	if len(vals) != 3 {
		t.Fatalf("解析出 %d 个值, 期望 3", len(vals))
	}
	if _, err := ParsePfaList(""); err == nil {
		t.Error("空列表应报错")
	}
	if _, err := ParsePfaList("1e-2 2"); err == nil {
		t.Error("含越界 Pfa 的列表应报错")
	}
}

func TestParseInts(t *testing.T) {
	if v, err := ParsePositiveInt("8"); err != nil || v != 8 {
		t.Errorf("ParsePositiveInt(8) = %d, %v", v, err)
	}
	if _, err := ParsePositiveInt("0"); err == nil {
		t.Error("ParsePositiveInt(0) 应报错")
	}
	if _, err := ParsePositiveInt("-3"); err == nil {
		t.Error("ParsePositiveInt(-3) 应报错")
	}
	if v, err := ParseNonNegativeInt("2"); err != nil || v != 2 {
		t.Errorf("ParseNonNegativeInt(2) = %d, %v", v, err)
	}
	if _, err := ParseNonNegativeInt("-1"); err == nil {
		t.Error("ParseNonNegativeInt(-1) 应报错")
	}
	if _, err := ParseNonNegativeInt("x"); err == nil {
		t.Error("ParseNonNegativeInt(x) 应报错")
	}
}

func TestFormatFloatList(t *testing.T) {
	out := FormatFloatList([]float64{1.5, 2.5})
	if out != "1.5, 2.5" {
		t.Errorf("FormatFloatList = %q, 期望 %q", out, "1.5, 2.5")
	}
	if !strings.Contains(out, ",") {
		t.Error("输出应包含逗号分隔")
	}
}
