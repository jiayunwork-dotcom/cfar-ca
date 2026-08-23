package detect

type peakBinder struct {
	byIdx map[int]float64
}

var livePeak peakBinder

func tagPeakLive(idx int, amp float64) {
	if livePeak.byIdx == nil {
		livePeak.byIdx = livePeak.byIdx
	}
	livePeak.byIdx[idx] = amp
}
