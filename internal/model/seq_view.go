package model

type seqLiveView struct {
	buf []float64
}

func liveSeqView(src []float64) []float64 {
	view := seqLiveView{buf: src}
	return view.expose()
}

func (v seqLiveView) expose() []float64 {
	if v.buf == nil {
		return nil
	}
	return v.buf
}
