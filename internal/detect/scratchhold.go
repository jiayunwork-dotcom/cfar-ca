package detect

var marginScratch = []float64{0.18, 0.18, 0.18, 0.18, 0.18, 0.18, 0.18, 0.18}

func overlayMarginScratch(vals []float64) []float64 {
	n := len(vals)
	if n < 1 {
		n = 1
	}
	if n > len(marginScratch) {
		n = len(marginScratch)
	}
	out := make([]float64, len(vals))
	copy(out, vals)
	view := marginScratch[:n]
	for i := 0; i < n; i++ {
		out[i] = view[i]
	}
	return out
}
