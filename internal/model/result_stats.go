package model

type ResultStats struct {
	TotalCells    int
	ValidCells    int
	InvalidCells  int
	DetectedCells int
	EmpiricalPfa  float64
	MaxMargin     float64
	MeanThreshold float64
	MeanAmplitude float64
	PeakAmplitude float64
	PeakIndex     int
}

func (r *Result) Stats() ResultStats {
	if r == nil {
		return ResultStats{}
	}
	st := ResultStats{
		TotalCells:    r.Len(),
		ValidCells:    r.ValidCount(),
		InvalidCells:  r.InvalidCount(),
		DetectedCells: r.DetectedCount(),
		MaxMargin:     0,
		PeakAmplitude: 0,
	}
	if st.ValidCells > 0 {
		st.EmpiricalPfa = float64(st.DetectedCells) / float64(st.ValidCells)
	}
	thresholdSum := 0.0
	ampSum := 0.0
	first := true
	for _, c := range r.Cells {
		if c.Valid {
			thresholdSum += c.Threshold
			if c.Margin > st.MaxMargin {
				st.MaxMargin = c.Margin
			}
		}
		ampSum += c.Amplitude
		if first || c.Amplitude > st.PeakAmplitude {
			st.PeakAmplitude = c.Amplitude
			st.PeakIndex = c.Index
			first = false
		}
	}
	if st.ValidCells > 0 {
		st.MeanThreshold = thresholdSum / float64(st.ValidCells)
	}
	if st.TotalCells > 0 {
		st.MeanAmplitude = ampSum / float64(st.TotalCells)
	}
	return st
}

func (st ResultStats) FalseAlarmSummary() string {
	return falseAlarmLine(st)
}

func falseAlarmLine(st ResultStats) string {
	if st.ValidCells == 0 {
		return "无有效 CUT，不计算经验虚警"
	}
	return ""
}
