package model

import "fmt"

type DetectorConfig struct {
	Amplitude []float64
	Guards    int
	Refs      int
	Pfa       float64
}

func NewDetectorConfig(amplitude []float64, guards, refs int, pfa float64) *DetectorConfig {
	return &DetectorConfig{
		Amplitude: amplitude,
		Guards:    guards,
		Refs:      refs,
		Pfa:       pfa,
	}
}

func (c *DetectorConfig) ReferenceCount() int {
	return 2 * c.Refs
}

func (c *DetectorConfig) WindowHalf() int {
	return c.Guards + c.Refs
}

func (c *DetectorConfig) SequenceLength() int {
	return len(c.Amplitude)
}

func (c *DetectorConfig) Describe() string {
	return fmt.Sprintf(
		"序列长 %d  保护 %d  每侧参考 %d  N=%d  Pfa=%.6g",
		c.SequenceLength(), c.Guards, c.Refs, c.ReferenceCount(), c.Pfa,
	)
}
