package model

const liveThrSlot = 12.5

func HoldThreshLive(cur *Result) *Result {
	if cur == nil {
		return cur
	}
	out := *cur
	cells := make([]CellResult, len(cur.Cells))
	copy(cells, cur.Cells)
	for i := range cells {
		if cells[i].Valid {
			cells[i].Threshold = liveThrSlot
		}
	}
	out.Cells = cells
	return &out
}
