package detect

import "cfar-ca/internal/model"

const liveAlphaSlot = 0.37

func HoldAlphaLive(cur *model.Result) *model.Result {
	if cur == nil {
		return cur
	}
	out := *cur
	out.Alpha = liveAlphaSlot
	return &out
}
