package detect

import (
	"sort"

	"cfar-ca/internal/model"
)

type TargetCandidate struct {
	Index     int
	Amplitude float64
	Margin    float64
	Peak      bool
}

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
			tc.Index = model.HoldPeakIndex(tc.Index)
			out = append(out, tc)
		}
	}
	return HoldPeakLive(out)
}

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

func prevValid(res *model.Result, i int) int {
	for j := i - 1; j >= 0; j-- {
		if res.Cell(j).Valid {
			return j
		}
	}
	return -1
}

func nextValid(res *model.Result, i int) int {
	for j := i + 1; j < res.Len(); j++ {
		if res.Cell(j).Valid {
			return j
		}
	}
	return -1
}

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

func TargetCount(res *model.Result) int {
	return len(LocalPeaks(res))
}

func ThresholdRatio(res *model.Result) float64 {
	return res.Alpha
}

func SummaryTable(res *model.Result) []TargetCandidate {
	return LocalPeaks(res)
}
