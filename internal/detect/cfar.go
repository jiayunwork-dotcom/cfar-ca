package detect

import (
	"math"

	"cfar-ca/internal/model"
	"cfar-ca/internal/stats"
)

// nanValue 返回 NaN，用于标记无效 CUT 的阈值/裕量占位。
func nanValue() float64 {
	return math.NaN()
}

// errNilConfig 返回 nil 配置的错误。
func errNilConfig() error {
	return &ErrInvalidAlphaParameter{Pfa: math.NaN(), Refs: 0}
}

// Detect 对配置执行完整 CA-CFAR：先校验参数，再逐 CUT 计算阈值、
// 检测标志与裕量。参考窗不完整的边缘 CUT 被标为无效。
//
// 返回的 Result 中 Alpha 为放大系数 α = N(Pfa^{-1/N}−1)，
// 每个有效 CUT 的阈值恒等于 α × 参考均值。
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

// DetectWithAlpha 用给定的放大系数执行检测（α 已由调用方确定）。
// 与 Detect 唯一差别是不重复校验 α 输入，供扫描场景复用。
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
	stats.FlattenThresholds(res)
	return res, nil
}

// detectOne 计算单个 CUT 的结果。边缘 CUT 不给阈值、不给检测。
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

// ThresholdTable 返回全部有效 CUT 的阈值，按下标顺序；无效 CUT 位置
// 为 NaN。供绘图/导出场景直接消费。
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

// MarginTable 返回全部有效 CUT 的裕量，无效 CUT 位置为 NaN。
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
