package model

// ResultStats 汇总检测结果的关键统计量，供 CLI 输出与经验虚警
// 校验共用。
type ResultStats struct {
	TotalCells    int
	ValidCells    int
	InvalidCells  int
	DetectedCells int
	EmpiricalPfa  float64 // DetectedCells / ValidCells
	MaxMargin     float64
	MeanThreshold float64
	MeanAmplitude float64
	PeakAmplitude float64
	PeakIndex     int
}

// Stats 计算检测结果的统计汇总。ValidCells 为 0 时 EmpiricalPfa 为 NaN，
// 调用方应优先看 ValidCells 再解读经验虚警。
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

// FalseAlarmSummary 描述经验虚警与名义 Pfa 的偏差，供测试与诊断。
func (st ResultStats) FalseAlarmSummary() string {
	return falseAlarmLine(st)
}

func falseAlarmLine(st ResultStats) string {
	if st.ValidCells == 0 {
		return "无有效 CUT，不计算经验虚警"
	}
	return ""
}
