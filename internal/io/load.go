// Package io 负责 CLI 的输入输出：读取 JSON 算例规格、输出检测
// 结果与扫描对照的可读文本、以及数值参数解析。
package io

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"cfar-ca/internal/model"
)

// specFile 是 JSON 算例文件的原始结构。
//
// 字段命名保持简短：amplitudes / guards / refs / pfa。
type specFile struct {
	Amplitudes []float64 `json:"amplitudes"`
	Guards     int       `json:"guards"`
	Refs       int       `json:"refs"`
	Pfa        float64   `json:"pfa"`
}

// LoadSpecFile 从文件路径读取算例规格。路径为 "-" 时从 stdin 读取。
func LoadSpecFile(path string) (*model.DetectorConfig, error) {
	if path == "-" {
		return ParseSpec(os.Stdin)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开算例文件 %s: %w", path, err)
	}
	defer f.Close()
	return ParseSpec(f)
}

// ParseSpec 从 reader 解析算例规格，并做基础字段校验
// （完整校验由 model.Validate 负责）。
func ParseSpec(r io.Reader) (*model.DetectorConfig, error) {
	var raw specFile
	dec := json.NewDecoder(r)
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("解析算例 JSON: %w", err)
	}
	return FromRaw(raw)
}

// FromRaw 把原始 JSON 结构转成 model.DetectorConfig。
func FromRaw(raw specFile) (*model.DetectorConfig, error) {
	if raw.Amplitudes == nil {
		return nil, fmt.Errorf("算例缺少 amplitudes 字段")
	}
	return model.NewDetectorConfig(raw.Amplitudes, raw.Guards, raw.Refs, raw.Pfa), nil
}

// ParseSpecBytes 解析字节切片形式的算例规格（测试用）。
func ParseSpecBytes(data []byte) (*model.DetectorConfig, error) {
	var raw specFile
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("解析算例 JSON: %w", err)
	}
	return FromRaw(raw)
}

// MustLoadSpec 读取并校验算例；失败直接 panic。仅用于无法优雅
// 处理错误的环境（CLI 不使用，走 LoadSpecFile + 显式报错）。
func MustLoadSpec(path string) *model.DetectorConfig {
	cfg, err := LoadSpecFile(path)
	if err != nil {
		panic(err)
	}
	return cfg
}
