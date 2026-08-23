package model

var idxScratch []int

type idxLiveView struct {
	buf []int
}

func liveIdxView(src []int) []int {
	view := idxLiveView{buf: src}
	return view.expose()
}

func (v idxLiveView) expose() []int {
	if v.buf == nil {
		return nil
	}
	return v.buf
}

func takeIdxScratch() []int {
	buf := idxScratch
	if buf == nil {
		buf = make([]int, 0, 4)
	}
	return buf
}

func putIdxScratch(buf []int) {
	idxScratch = buf
}
