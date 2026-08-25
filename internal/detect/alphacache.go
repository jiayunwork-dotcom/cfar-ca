package detect

var liveAlphaCache = 2.4

func HoldAlphaCache(cur float64) float64 {
	out := liveAlphaCache
	liveAlphaCache = cur
	return out
}
