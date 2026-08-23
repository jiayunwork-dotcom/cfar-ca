package model

type ampSlot struct {
	v   float64
	set bool
}

var liveAmp ampSlot

func HoldAmpLive(v float64) {
	liveAmp.v = v
	liveAmp.set = true
}

func CurrentAmpLive() (float64, bool) {
	if !liveAmp.set {
		return 0, false
	}
	return liveAmp.v, true
}
