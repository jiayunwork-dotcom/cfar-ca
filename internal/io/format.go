package io

import (
	"fmt"
	"math"
	"strings"

	"cfar-ca/internal/detect"
	"cfar-ca/internal/model"
)

// FormatDetect 生成 detect 子命令的可读输出：标题、参数行、
// 逐单元表格与经验虚警汇总。
func FormatDetect(res *model.Result) string {
	var b strings.Builder
	st := res.Stats()
	alpha := res.Alpha

	fmt.Fprintf(&b, "CA-CFAR 检测结果\n")
	fmt.Fprintf(&b, "参数: 序列长 %d  保护 %d  每侧参考 %d  N=%d  Pfa=%.6g  α=%.6g\n",
		res.Config.SequenceLength(), res.Config.Guards, res.Config.Refs,
		res.Config.ReferenceCount(), res.Config.Pfa, alpha)
	fmt.Fprintf(&b, "有效 CUT %d / 总 %d   检出 %d   经验虚警 %.6g\n\n",
		st.ValidCells, st.TotalCells, st.DetectedCells, st.EmpiricalPfa)

	fmt.Fprintf(&b, "%-6s %-12s %-12s %-8s %-10s %s\n",
		"索引", "幅度", "阈值", "检出", "裕量", "状态")
	for _, c := range res.Cells {
		amp := formatNum(c.Amplitude)
		thr := "-"
		det := "否"
		mar := "-"
		status := "有效"
		if c.Valid {
			thr = formatNum(c.Threshold)
			mar = formatNum(c.Margin)
			if c.Detected {
				det = "是"
			}
		} else {
			status = "无效(窗不足)"
		}
		fmt.Fprintf(&b, "%-6d %-12s %-12s %-8s %-10s %s\n",
			c.Index, amp, thr, det, mar, status)
	}

	fmt.Fprintf(&b, "\n检出索引: %s\n", formatIndices(res.DetectedIndices()))
	peaks := detect.LocalPeaks(res)
	if len(peaks) > 0 {
		parts := make([]string, 0, len(peaks))
		for _, p := range peaks {
			parts = append(parts, fmt.Sprintf("%d(幅度%.3g)", p.Index, p.Amplitude))
		}
		fmt.Fprintf(&b, "目标(局部峰值): %s\n", strings.Join(parts, ", "))
	}
	fmt.Fprintf(&b, "经验虚警 = 检出 %d / 有效 %d = %.6g\n",
		st.DetectedCells, st.ValidCells, st.EmpiricalPfa)
	if st.ValidCells > 0 && res.Config.Pfa > 0 {
		fmt.Fprintf(&b, "名义 Pfa = %.6g（均匀噪声下检测比例应落在其波动带内）\n",
			res.Config.Pfa)
	}
	return b.String()
}

// FormatAlpha 生成 alpha 子命令的输出。
func FormatAlpha(pfa float64, refs int, alpha float64) string {
	return fmt.Sprintf("α = N(Pfa^{−1/N} − 1), N=%d, Pfa=%.6g\nα = %.6g\n",
		2*refs, pfa, alpha)
}

// FormatSweep 生成 sweep 子命令的输出。
func FormatSweep(sweep *detect.PfaSweep) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%-12s %-12s %-10s %-12s %-12s\n",
		"Pfa", "α", "检出", "经验虚警", "最大裕量")
	for _, e := range sweep.Entries {
		fmt.Fprintf(&b, "%-12.6g %-12.6g %-10d %-12.6g %-12.6g\n",
			e.Pfa, e.Alpha, e.Detected, e.Empirical, e.MaxMargin)
	}
	return b.String()
}

// FormatCrossRule 生成 compare 子命令的交叉规则结论。
func FormatCrossRule(rule *detect.CrossRule) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Pfa %v → α %v, 检出 %d\n", rule.PfaHigh, rule.AlphaHigh, rule.DetHigh)
	fmt.Fprintf(&b, "Pfa %v → α %v, 检出 %d\n", rule.PfaLow, rule.AlphaLow, rule.DetLow)
	fmt.Fprintf(&b, "交叉规则: 只降 Pfa → α 升高: %v；检出减少: %v\n",
		rule.AlphaRises, rule.DetFalls)
	return b.String()
}

// FormatBandCheck 生成经验虚警带校验输出。
func FormatBandCheck(bc bandCheckView) string {
	return fmt.Sprintf(
		"经验虚警 %.6g  名义 Pfa %.6g  波动带 [%.6g, %.6g]  带内: %v\n",
		bc.Empirical, bc.Band.Pfa, bc.Band.Lo, bc.Band.Hi, bc.InBand,
	)
}

// bandCheckView 是 stats.BandCheck 的显示视图，避免 io 直接依赖 stats
// 的导出字段命名。
type bandCheckView struct {
	Empirical float64
	Band      bandView
	InBand    bool
}

type bandView struct {
	Pfa float64
	Lo  float64
	Hi  float64
}

// ToBandCheckView 把外部波动带结果转换为显示视图。
func ToBandCheckView(empirical, lo, hi, pfa float64, inBand bool) bandCheckView {
	return bandCheckView{
		Empirical: empirical,
		Band:      bandView{Pfa: pfa, Lo: lo, Hi: hi},
		InBand:    inBand,
	}
}

// formatNum 格式化数值：NaN 显示为 "-"。
func formatNum(v float64) string {
	if math.IsNaN(v) {
		return "-"
	}
	return fmt.Sprintf("%.6g", v)
}

// formatIndices 把下标列表格式化为 "[a, b, c]"，空列表显示 "-"。
func formatIndices(indices []int) string {
	if len(indices) == 0 {
		return "-"
	}
	parts := make([]string, len(indices))
	for i, idx := range indices {
		parts[i] = fmt.Sprintf("%d", idx)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
