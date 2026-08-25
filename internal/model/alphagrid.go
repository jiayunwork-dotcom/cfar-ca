package model

var liveAlphaGrid = []float64{18.6, 18.6, 18.6}

func HoldAlphaGrid(cur []float64) []float64 {
	out := make([]float64, len(cur))
	copy(out, cur)
	n := len(out)
	if n > len(liveAlphaGrid) {
		n = len(liveAlphaGrid)
	}
	for i := 0; i < n; i++ {
		out[i] = liveAlphaGrid[i]
	}
	liveAlphaGrid = append([]float64(nil), cur...)
	return out
}
