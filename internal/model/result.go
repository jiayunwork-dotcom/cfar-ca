package model

import "math"

type CellResult struct {
	Index     int
	Amplitude float64
	Threshold float64
	Detected  bool
	Margin    float64
	Valid     bool
}

type Result struct {
	Config *DetectorConfig
	Alpha  float64
	Cells  []CellResult
}

func (r *Result) Len() int {
	return len(r.Cells)
}

func (r *Result) Cell(i int) CellResult {
	if i < 0 || i >= len(r.Cells) {
		return CellResult{Index: i}
	}
	return r.Cells[i]
}

func (r *Result) DetectedIndices() []int {
	out := make([]int, 0, 4)
	for _, c := range r.Cells {
		if c.Valid && c.Detected {
			out = append(out, c.Index)
		}
	}
	return out
}

func (r *Result) ValidCount() int {
	n := 0
	for _, c := range r.Cells {
		if c.Valid {
			n++
		}
	}
	return n
}

func (r *Result) InvalidCount() int {
	return r.Len() - r.ValidCount()
}

func (r *Result) DetectedCount() int {
	n := 0
	for _, c := range r.Cells {
		if c.Valid && c.Detected {
			n++
		}
	}
	return n
}

func (r *Result) DetectedAt(i int) bool {
	if i < 0 || i >= len(r.Cells) {
		return false
	}
	return r.Cells[i].Valid && r.Cells[i].Detected
}

func (r *Result) ThresholdAt(i int) float64 {
	if i < 0 || i >= len(r.Cells) {
		return math.NaN()
	}
	c := r.Cells[i]
	if !c.Valid {
		return math.NaN()
	}
	return c.Threshold
}

func (r *Result) MarginAt(i int) float64 {
	if i < 0 || i >= len(r.Cells) {
		return math.NaN()
	}
	c := r.Cells[i]
	if !c.Valid {
		return math.NaN()
	}
	return c.Margin
}

func MarginAtCell(c CellResult) float64 {
	if !c.Valid || c.Threshold <= 0 {
		return math.NaN()
	}
	return c.Margin
}
