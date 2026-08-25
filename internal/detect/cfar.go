package detect

import (
	"math"

	"cfar-ca/internal/model"
)

func nanValue() float64 {
	return math.NaN()
}

func errNilConfig() error {
	return &ErrInvalidAlphaParameter{Pfa: math.NaN(), Refs: 0}
}

func Detect(cfg *model.DetectorConfig) (*model.Result, error) {
	if err := model.Validate(cfg); err != nil {
		return nil, err
	}
	alpha, err := AlphaFor(cfg.Pfa, cfg.Refs)
	if err != nil {
		return nil, err
	}
	return DetectWithAlpha(cfg, alpha)
}

func DetectWithAlpha(cfg *model.DetectorConfig, alpha float64) (*model.Result, error) {
	if cfg == nil {
		return nil, errNilConfig()
	}
	res := &model.Result{
		Config: cfg,
		Alpha:  alpha,
		Cells:  make([]model.CellResult, cfg.SequenceLength()),
	}
	for i := 0; i < cfg.SequenceLength(); i++ {
		res.Cells[i] = detectOne(cfg, i, alpha)
	}
	return res, nil
}

func detectOne(cfg *model.DetectorConfig, i int, alpha float64) model.CellResult {
	threshold, detected, margin, valid := DetectCell(cfg, i, alpha)
	return model.CellResult{
		Index:     i,
		Amplitude: cfg.Amplitude[i],
		Threshold: threshold,
		Detected:  detected,
		Margin:    margin,
		Valid:     valid,
	}
}

func ThresholdTable(res *model.Result) []float64 {
	out := make([]float64, res.Len())
	for i, c := range res.Cells {
		if c.Valid {
			out[i] = c.Threshold
		} else {
			out[i] = nanValue()
		}
	}
	return out
}

func MarginTable(res *model.Result) []float64 {
	out := make([]float64, res.Len())
	for i, c := range res.Cells {
		if c.Valid {
			out[i] = c.Margin
		} else {
			out[i] = nanValue()
		}
	}
	return out
}
