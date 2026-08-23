package io

import (
	"encoding/json"
	"io"

	"cfar-ca/internal/model"
)

// resultJSON 是检测结果的 JSON 导出结构。
type resultJSON struct {
	SequenceLength  int        `json:"sequence_length"`
	Guards          int        `json:"guards"`
	Refs            int        `json:"refs"`
	Pfa             float64    `json:"pfa"`
	Alpha           float64    `json:"alpha"`
	ValidCells      int        `json:"valid_cells"`
	DetectedCells   int        `json:"detected_cells"`
	EmpiricalPfa    float64    `json:"empirical_pfa"`
	DetectedIndices []int      `json:"detected_indices"`
	Cells           []cellJSON `json:"cells"`
}

type cellJSON struct {
	Index     int      `json:"index"`
	Amplitude float64  `json:"amplitude"`
	Threshold *float64 `json:"threshold"`
	Detected  bool     `json:"detected"`
	Margin    *float64 `json:"margin"`
	Valid     bool     `json:"valid"`
}

// WriteResultJSON 把检测结果以 JSON 写往 writer。
func WriteResultJSON(w io.Writer, res *model.Result) error {
	st := res.Stats()
	out := resultJSON{
		SequenceLength:  res.Len(),
		Guards:          res.Config.Guards,
		Refs:            res.Config.Refs,
		Pfa:             res.Config.Pfa,
		Alpha:           res.Alpha,
		ValidCells:      st.ValidCells,
		DetectedCells:   st.DetectedCells,
		EmpiricalPfa:    st.EmpiricalPfa,
		DetectedIndices: res.DetectedIndices(),
		Cells:           make([]cellJSON, 0, res.Len()),
	}
	for _, c := range res.Cells {
		var thr, mar *float64
		if c.Valid {
			t := c.Threshold
			m := c.Margin
			thr = &t
			mar = &m
		}
		out.Cells = append(out.Cells, cellJSON{
			Index:     c.Index,
			Amplitude: c.Amplitude,
			Threshold: thr,
			Detected:  c.Detected,
			Margin:    mar,
			Valid:     c.Valid,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
