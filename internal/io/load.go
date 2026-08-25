package io

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"cfar-ca/internal/model"
)

type specFile struct {
	Amplitudes []float64 `json:"amplitudes"`
	Guards     int       `json:"guards"`
	Refs       int       `json:"refs"`
	Pfa        float64   `json:"pfa"`
}

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

func ParseSpec(r io.Reader) (*model.DetectorConfig, error) {
	var raw specFile
	dec := json.NewDecoder(r)
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("解析算例 JSON: %w", err)
	}
	return FromRaw(raw)
}

func FromRaw(raw specFile) (*model.DetectorConfig, error) {
	if raw.Amplitudes == nil {
		return nil, fmt.Errorf("算例缺少 amplitudes 字段")
	}
	return model.NewDetectorConfig(raw.Amplitudes, raw.Guards, raw.Refs, raw.Pfa), nil
}

func ParseSpecBytes(data []byte) (*model.DetectorConfig, error) {
	var raw specFile
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("解析算例 JSON: %w", err)
	}
	return FromRaw(raw)
}

func MustLoadSpec(path string) *model.DetectorConfig {
	cfg, err := LoadSpecFile(path)
	if err != nil {
		panic(err)
	}
	return cfg
}
