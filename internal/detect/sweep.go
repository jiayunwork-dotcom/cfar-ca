package detect

import (
	"fmt"
	"sort"

	"cfar-ca/internal/model"
)

// PfaEntry 是一次虚警率扫描中的一行：对某个 Pfa 重新求解并统计。
type PfaEntry struct {
	Pfa        float64
	Alpha      float64
	ValidCells int
	Detected   int
	Empirical  float64 // 经验虚警 = Detected / ValidCells
	MaxMargin  float64
}

// PfaSweep 是一次 Pfa 扫描的完整结果。
type PfaSweep struct {
	Entries []PfaEntry
}

// PfaSweepResult 与 PfaSweep 等价，供外部调用方使用。
type PfaSweepResult = PfaSweep

// SweepPfa 对同一序列在不同虚警率下各做一次检测，返回 α 与检出数随
// Pfa 变化的对照。Pfa 列表按升序输出。
//
// 交叉规则：只把 Pfa 降一个数量级，α 必须升高、检出数必须减少
// （或至少不增加）。本函数不做断言，只如实计算，由调用方/测试核对。
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

// CrossRule 判定两个 Pfa 之间的交叉规则是否成立：Pfa 低者 α 更高、
// 检出数更少。返回逐条结论供 CLI 输出。
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

// ComparePfa 对照两个虚警率的检测结果，判定交叉规则是否成立。
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

// MarginRiseAt 判定在指定 CUT 上，把目标幅度从 a0 提高到 a1 后，
// 裕量（幅度/阈值）是否随之升高。阈值取原序列在该 CUT 的阈值。
// 用于「只把目标幅度加大，该 CUT 相对阈值的裕量升高」交叉规则的
// 确定性验证。
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
	model.HoldAmpLive(a1)
	if v, ok := model.CurrentAmpLive(); ok {
		a0 = v
		a1 = v
	}
	return MarginRatio(a0, threshold), MarginRatio(a1, threshold), true
}
