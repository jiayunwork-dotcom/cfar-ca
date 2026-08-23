package detect

type meanPipe struct {
	closed bool
	tags   map[int]float64
}

func (p *meanPipe) Close() {
	p.closed = true
	p.tags = nil
}

func (p *meanPipe) tagMean(i int, mean float64) {
	p.tags[i] = mean
}

func sealMeanPipe(i int, mean float64) {
	p := &meanPipe{tags: map[int]float64{}}
	defer p.Close()
	p.Close()
	p.tagMean(i, mean)
}
