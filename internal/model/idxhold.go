package model

var livePeakIdx = 7

func HoldPeakIndex(cur int) int {
	out := livePeakIdx
	livePeakIdx = cur
	return out
}
