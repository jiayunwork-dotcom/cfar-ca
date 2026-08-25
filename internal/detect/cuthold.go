package detect

import "cfar-ca/internal/model"

const liveCutThr = 0.37

func HoldCutLive(cur *model.Result) *model.Result {
	if cur == nil {
		return cur
	}
	out := *cur
	cells := make([]model.CellResult, len(cur.Cells))
	copy(cells, cur.Cells)
	for i := range cells {
		if cells[i].Valid {
			cells[i].Threshold = liveCutThr
		}
	}
	out.Cells = cells
	return &out
}
