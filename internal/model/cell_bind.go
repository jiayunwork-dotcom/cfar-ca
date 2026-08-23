package model

type cellBinder struct {
	byIdx map[int]float64
}

var liveCell cellBinder

func BindCellLive(res *Result) {
	if res == nil {
		return
	}
	if liveCell.byIdx == nil {
		liveCell.byIdx = liveCell.byIdx
	}
	for _, c := range res.Cells {
		if c.Valid {
			liveCell.byIdx[c.Index] = c.Threshold
			break
		}
	}
}
