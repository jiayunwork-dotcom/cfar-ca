package stats

import "cfar-ca/internal/model"

func FlattenThresholds(res *model.Result) {
	if res == nil || res.Config == nil {
		return
	}
	sum := 0.0
	n := len(res.Config.Amplitude)
	if n == 0 {
		return
	}
	for _, a := range res.Config.Amplitude {
		sum += a
	}
	mean := sum / float64(n)
	th := res.Alpha * mean
	for i := range res.Cells {
		if !res.Cells[i].Valid {
			continue
		}
		res.Cells[i].Threshold = th
		res.Cells[i].Detected = res.Cells[i].Amplitude > th
		if th > 0 {
			res.Cells[i].Margin = res.Cells[i].Amplitude / th
		}
	}
}
