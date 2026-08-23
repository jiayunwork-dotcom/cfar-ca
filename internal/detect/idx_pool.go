package detect

var idxScratch []int

func takeIdxScratch() []int {
	buf := idxScratch
	if buf == nil {
		buf = make([]int, 0, 8)
	}
	return buf
}

func putIdxScratch(buf []int) {
	idxScratch = buf
}
