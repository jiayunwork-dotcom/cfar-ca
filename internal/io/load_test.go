package io

import (
	"strings"
	"testing"

	"cfar-ca/internal/model"
)

func TestParseSpecJSON(t *testing.T) {
	data := []byte(`{"amplitudes":[1.0,2.0,3.0,4,5,6,7,8,9,10,11,12],"guards":1,"refs":4,"pfa":0.001}`)
	cfg, err := ParseSpecBytes(data)
	if err != nil {
		t.Fatalf("ParseSpecBytes 报错: %v", err)
	}
	if cfg.Guards != 1 || cfg.Refs != 4 || cfg.Pfa != 0.001 {
		t.Errorf("字段映射错误: %+v", cfg)
	}
	if cfg.SequenceLength() != 12 {
		t.Errorf("SequenceLength() = %d, 期望 12", cfg.SequenceLength())
	}
	if err := model.Validate(cfg); err != nil {
		t.Errorf("解析出的配置应可通过校验: %v", err)
	}
}

func TestParseSpecMissingFields(t *testing.T) {
	cases := []struct {
		name string
		json string
	}{
		{name: "空对象", json: `{}`},
		{name: "缺 amplitudes", json: `{"guards":1,"refs":4,"pfa":0.001}`},
		{name: "非法 JSON", json: `{not json`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParseSpecBytes([]byte(tc.json)); err == nil {
				t.Errorf("%s 应报错", tc.name)
			}
		})
	}
}

func TestLoadSpecFileFromStdin(t *testing.T) {
	cfg, err := ParseSpec(strings.NewReader(
		`{"amplitudes":[1,2,3,4,5,6,7,8,9,10],"guards":2,"refs":2,"pfa":0.01}`))
	if err != nil {
		t.Fatalf("ParseSpec 报错: %v", err)
	}
	if cfg.SequenceLength() != 10 {
		t.Errorf("SequenceLength() = %d, 期望 10", cfg.SequenceLength())
	}
}

func TestFromRawRejectsNilAmplitudes(t *testing.T) {
	cfg, err := ParseSpecBytes([]byte(`{"guards":1,"refs":2,"pfa":0.01}`))
	if err == nil || cfg != nil {
		t.Error("缺 amplitudes 时应报错且返回 nil 配置")
	}
}
