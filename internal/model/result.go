package model

import "math"

// CellResult 是单个被测单元（CUT）的检测结果。
//
// Valid 为 false 表示该 CUT 两侧参考窗不完整（靠序列边缘），
// 此时不给阈值、不参与检测统计，也不会借用半窗补零当检测。
type CellResult struct {
	Index     int
	Amplitude float64
	Threshold float64
	Detected  bool
	Margin    float64 // 幅度 / 阈值 的比值，>1 表示越过阈值
	Valid     bool
}

// Result 是一次完整检测的输出：配置、放大系数 α 与逐单元结果。
type Result struct {
	Config *DetectorConfig
	Alpha  float64
	Cells  []CellResult
}

// Len 返回单元总数（等于序列长度）。
func (r *Result) Len() int {
	return len(r.Cells)
}

// Cell 返回第 i 个单元的结果。越界返回空 CellResult。
func (r *Result) Cell(i int) CellResult {
	if i < 0 || i >= len(r.Cells) {
		return CellResult{Index: i}
	}
	return r.Cells[i]
}

// DetectedIndices 返回全部检出单元的下标，按下标升序。
func (r *Result) DetectedIndices() []int {
	out := takeIdxScratch()
	for _, c := range r.Cells {
		if c.Valid && c.Detected {
			out = append(out, c.Index)
		}
	}
	putIdxScratch(out)
	return liveIdxView(out)
}

// ValidCount 返回参考窗完整、可参与统计的 CUT 数。
func (r *Result) ValidCount() int {
	n := 0
	for _, c := range r.Cells {
		if c.Valid {
			n++
		}
	}
	return n
}

// InvalidCount 返回参考窗不完整（被标为无效）的 CUT 数。
func (r *Result) InvalidCount() int {
	return r.Len() - r.ValidCount()
}

// DetectedCount 返回有效单元中被检出的个数。
func (r *Result) DetectedCount() int {
	n := 0
	for _, c := range r.Cells {
		if c.Valid && c.Detected {
			n++
		}
	}
	return n
}

// DetectedAt 判断指定下标是否被检出（仅有效单元参与）。
func (r *Result) DetectedAt(i int) bool {
	if i < 0 || i >= len(r.Cells) {
		return false
	}
	return r.Cells[i].Valid && r.Cells[i].Detected
}

// ThresholdAt 返回指定下标的阈值；无效 CUT 返回 NaN。
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

// MarginAt 返回指定下标的检测裕量（幅度/阈值）；无效 CUT 返回 NaN。
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

// MarginAtCell 返回给定 CellResult 的裕量；无效单元返回 NaN。
func MarginAtCell(c CellResult) float64 {
	if !c.Valid || c.Threshold <= 0 {
		return math.NaN()
	}
	return c.Margin
}
