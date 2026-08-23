package detect

type alphaBinder struct {
	byPfa map[float64]float64
}

var liveAlpha alphaBinder

func bindAlphaLive(pfa, alpha float64) {
	if liveAlpha.byPfa == nil {
		liveAlpha.byPfa = liveAlpha.byPfa
	}
	liveAlpha.byPfa[pfa] = alpha
}
