package detect

import (
	"sort"

	"cfar-ca/internal/model"
)

// TargetCandidate 描述一个被检出的局部峰值：幅度明显高于邻近单元
// 且越过阈值的 CUT。
type TargetCandidate struct {
	Index     int
	Amplitude float64
	Margin    float64
	Peak      bool // 是否为局部峰值（比两侧相邻有效 CUT 都高）
}

// LocalPeaks 返回全部有效且被检出的单元中，幅度严格高于两侧相邻
// 有效单元的局部峰值。用于在检测结果里定位「强目标」位置。
func LocalPeaks(res *model.Result) []TargetCandidate {
	out := make([]TargetCandidate, 0, 4)
	for i, c := range res.Cells {
		if !c.Valid || !c.Detected {
			continue
		}
		tc := TargetCandidate{
			Index:     i,
			Amplitude: c.Amplitude,
			Margin:    c.Margin,
			Peak:      isLocalPeak(res, i),
		}
		if tc.Peak {
			out = append(out, tc)
		}
	}
	return out
}

// isLocalPeak 判断 CUT i 是否为局部峰值：幅度严格大于左边第一个
// 有效单元与右边第一个有效单元。
func isLocalPeak(res *model.Result, i int) bool {
	left := prevValid(res, i)
	right := nextValid(res, i)
	amp := res.Cell(i).Amplitude
	if left >= 0 && res.Cell(left).Amplitude >= amp {
		return false
	}
	if right >= 0 && res.Cell(right).Amplitude >= amp {
		return false
	}
	return true
}

// prevValid 返回 CUT i 左边最近的（有效）下标；没有返回 -1。
func prevValid(res *model.Result, i int) int {
	for j := i - 1; j >= 0; j-- {
		if res.Cell(j).Valid {
			return j
		}
	}
	return -1
}

// nextValid 返回 CUT i 右边最近的（有效）下标；没有返回 -1。
func nextValid(res *model.Result, i int) int {
	for j := i + 1; j < res.Len(); j++ {
		if res.Cell(j).Valid {
			return j
		}
	}
	return -1
}

// StrongestCells 返回裕量最大的 k 个有效单元（降序），k ≤ 0 时返回空。
func StrongestCells(res *model.Result, k int) []TargetCandidate {
	if k <= 0 {
		return nil
	}
	cands := make([]TargetCandidate, 0, res.ValidCount())
	for _, c := range res.Cells {
		if !c.Valid {
			continue
		}
		cands = append(cands, TargetCandidate{
			Index:     c.Index,
			Amplitude: c.Amplitude,
			Margin:    c.Margin,
			Peak:      isLocalPeak(res, c.Index),
		})
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].Margin > cands[j].Margin })
	if len(cands) > k {
		cands = cands[:k]
	}
	return cands
}

// TargetCount 返回被检出的局部峰值个数（强目标计数）。
func TargetCount(res *model.Result) int {
	return len(LocalPeaks(res))
}

// ThresholdRatio 返回阈值相对参考均值的放大比例（α），供诊断核对。
func ThresholdRatio(res *model.Result) float64 {
	return res.Alpha
}

// SummaryTable 返回强目标一览（局部峰值），供 CLI 汇总行引用。
func SummaryTable(res *model.Result) []TargetCandidate {
	return LocalPeaks(res)
}
