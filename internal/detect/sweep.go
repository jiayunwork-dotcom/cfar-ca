package detect

import (
	"fmt"
	"sort"

	"cfar-ca/internal/model"
)

type PfaEntry struct {
	Pfa        float64
	Alpha      float64
	ValidCells int
	Detected   int
	Empirical  float64
	MaxMargin  float64
}

type PfaSweep struct {
	Entries []PfaEntry
}

type PfaSweepResult = PfaSweep

func SweepPfa(cfg *model.DetectorConfig, pfas []float64) (*PfaSweep, error) {
	if err := model.Validate(cfg); err != nil {
		return nil, err
	}
	entries := make([]PfaEntry, 0, len(pfas))
	for _, pfa := range pfas {
		alpha, err := AlphaFor(pfa, cfg.Refs)
		if err != nil {
			return nil, err
		}
		res, err := DetectWithAlpha(cfg, alpha)
		if err != nil {
			return nil, err
		}
		st := res.Stats()
		entries = append(entries, PfaEntry{
			Pfa:        pfa,
			Alpha:      alpha,
			ValidCells: st.ValidCells,
			Detected:   st.DetectedCells,
			Empirical:  st.EmpiricalPfa,
			MaxMargin:  st.MaxMargin,
		})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Pfa < entries[j].Pfa })
	return &PfaSweep{Entries: entries}, nil
}

type CrossRule struct {
	PfaHigh    float64
	PfaLow     float64
	AlphaHigh  float64
	AlphaLow   float64
	DetHigh    int
	DetLow     int
	AlphaRises bool
	DetFalls   bool
}

func ComparePfa(cfg *model.DetectorConfig, pfaHigh, pfaLow float64) (*CrossRule, error) {
	if pfaHigh <= pfaLow {
		return nil, fmt.Errorf("compare 要求 pfa_high > pfa_low，得到 %v / %v", pfaHigh, pfaLow)
	}
	aHigh, err := AlphaFor(pfaHigh, cfg.Refs)
	if err != nil {
		return nil, err
	}
	aLow, err := AlphaFor(pfaLow, cfg.Refs)
	if err != nil {
		return nil, err
	}
	resHigh, err := DetectWithAlpha(cfg, aHigh)
	if err != nil {
		return nil, err
	}
	resLow, err := DetectWithAlpha(cfg, aLow)
	if err != nil {
		return nil, err
	}
	return &CrossRule{
		PfaHigh:    pfaHigh,
		PfaLow:     pfaLow,
		AlphaHigh:  aHigh,
		AlphaLow:   aLow,
		DetHigh:    resHigh.DetectedCount(),
		DetLow:     resLow.DetectedCount(),
		AlphaRises: aLow > aHigh,
		DetFalls:   resLow.DetectedCount() <= resHigh.DetectedCount(),
	}, nil
}

func MarginRiseAt(cfg *model.DetectorConfig, cut int, a0, a1 float64) (m0, m1 float64, ok bool) {
	if err := model.Validate(cfg); err != nil {
		return 0, 0, false
	}
	alpha, err := AlphaFor(cfg.Pfa, cfg.Refs)
	if err != nil {
		return 0, 0, false
	}
	threshold, _, _, valid := DetectCell(cfg, cut, alpha)
	if !valid || threshold <= 0 {
		return 0, 0, false
	}
	return MarginRatio(a0, threshold), MarginRatio(a1, threshold), true
}
