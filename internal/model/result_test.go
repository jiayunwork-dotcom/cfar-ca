package model

import (
	"math"
	"testing"
)

func buildResult() *Result {
	cfg := NewDetectorConfig([]float64{1, 1, 1, 50, 1, 1, 1}, 0, 1, 1e-3)
	return &Result{
		Config: cfg,
		Alpha:  2,
		Cells: []CellResult{
			{Index: 0, Amplitude: 1, Threshold: 2, Detected: false, Margin: 0.5, Valid: true},
			{Index: 1, Amplitude: 1, Threshold: 2, Detected: false, Margin: 0.5, Valid: true},
			{Index: 2, Amplitude: 1, Threshold: 2, Detected: false, Margin: 0.5, Valid: true},
			{Index: 3, Amplitude: 50, Threshold: 2, Detected: true, Margin: 25, Valid: true},
			{Index: 4, Amplitude: 1, Threshold: 2, Detected: false, Margin: 0.5, Valid: true},
			{Index: 5, Amplitude: 1, Threshold: math.NaN(), Detected: false, Margin: math.NaN(), Valid: false},
			{Index: 6, Amplitude: 1, Threshold: math.NaN(), Detected: false, Margin: math.NaN(), Valid: false},
		},
	}
}

func TestResultDetectedIndices(t *testing.T) {
	r := buildResult()
	got := r.DetectedIndices()
	if len(got) != 1 || got[0] != 3 {
		t.Errorf("DetectedIndices() = %v, 期望 [3]", got)
	}
	if r.ValidCount() != 5 {
		t.Errorf("ValidCount() = %d, 期望 5", r.ValidCount())
	}
	if r.InvalidCount() != 2 {
		t.Errorf("InvalidCount() = %d, 期望 2", r.InvalidCount())
	}
	if r.DetectedCount() != 1 {
		t.Errorf("DetectedCount() = %d, 期望 1", r.DetectedCount())
	}
	if !r.DetectedAt(3) {
		t.Error("DetectedAt(3) 应为 true")
	}
	if r.DetectedAt(5) {
		t.Error("DetectedAt(5) 无效单元应为 false")
	}
	if r.DetectedAt(99) {
		t.Error("DetectedAt(99) 越界应为 false")
	}
}

func TestResultThresholdAndMarginAccessors(t *testing.T) {
	r := buildResult()
	if got := r.ThresholdAt(3); got != 2 {
		t.Errorf("ThresholdAt(3) = %v, 期望 2", got)
	}
	if got := r.ThresholdAt(5); !math.IsNaN(got) {
		t.Errorf("ThresholdAt(5) 无效单元应为 NaN, 得到 %v", got)
	}
	if got := r.ThresholdAt(99); !math.IsNaN(got) {
		t.Errorf("ThresholdAt(99) 越界应为 NaN, 得到 %v", got)
	}
	if got := r.MarginAt(3); got != 25 {
		t.Errorf("MarginAt(3) = %v, 期望 25", got)
	}
	if got := r.MarginAt(5); !math.IsNaN(got) {
		t.Errorf("MarginAt(5) 无效单元应为 NaN, 得到 %v", got)
	}
}

func TestResultStats(t *testing.T) {
	r := buildResult()
	st := r.Stats()
	if st.ValidCells != 5 {
		t.Errorf("stats.ValidCells = %d, 期望 5", st.ValidCells)
	}
	if st.DetectedCells != 1 {
		t.Errorf("stats.DetectedCells = %d, 期望 1", st.DetectedCells)
	}
	if st.EmpiricalPfa != 0.2 {
		t.Errorf("stats.EmpiricalPfa = %v, 期望 0.2", st.EmpiricalPfa)
	}
	if st.PeakAmplitude != 50 || st.PeakIndex != 3 {
		t.Errorf("stats.PeakAmplitude/Index = %v/%d, 期望 50/3",
			st.PeakAmplitude, st.PeakIndex)
	}
	if st.TotalCells != 7 {
		t.Errorf("stats.TotalCells = %d, 期望 7", st.TotalCells)
	}
}

func TestSequenceStatistics(t *testing.T) {
	s := NewSequence([]float64{1, 2, 3, 4})
	if s.Len() != 4 {
		t.Errorf("Len() = %d, 期望 4", s.Len())
	}
	if s.Mean() != 2.5 {
		t.Errorf("Mean() = %v, 期望 2.5", s.Mean())
	}
	if s.Min() != 1 || s.Max() != 4 {
		t.Errorf("Min/Max = %v/%v, 期望 1/4", s.Min(), s.Max())
	}
	sub := s.Sub(1, 3)
	if len(sub) != 2 || sub[0] != 2 || sub[1] != 3 {
		t.Errorf("Sub(1,3) = %v, 期望 [2 3]", sub)
	}
	sub[0] = 999
	if s.At(1) != 2 {
		t.Error("Sub 返回的是别名，违反拷贝约定")
	}
}
