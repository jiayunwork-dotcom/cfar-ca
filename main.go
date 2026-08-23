// Command cfar-ca 是单元平均 CA-CFAR（Cell-Averaging Constant False
// Alarm Rate）检测核算的命令行工具。
//
// 用法：
//
//	cfar-ca detect example/target-in-noise.json  逐单元阈值与检出、经验虚警
//	cfar-ca detect -json <file>                  同前，JSON 导出
//	cfar-ca sweep -pfa "1e-2 1e-3 1e-4" <file>   Pfa 扫描：α 与检出数
//	cfar-ca compare -high 1e-3 -low 1e-4 <file>  交叉规则对照
//	cfar-ca alpha -pfa 1e-3 -refs 8              打印 α = N(Pfa^{−1/N}−1)
//	cfar-ca synth -n 4000 -guards 2 -refs 16 -pfa 1e-3  均匀噪声经验虚警带
//	cfar-ca version                               版本
//
// 非法输入（Pfa 越界、负幅度、参考窗伸出序列、保护≥参考、窗长 0、
// 文件不存在等）一律打印到 stderr 并返回非零退出码。
// 求解主逻辑在 internal/，本文件只做命令行接线。
package main

import (
	"flag"
	"fmt"
	"os"

	"cfar-ca/internal/detect"
	"cfar-ca/internal/io"
	"cfar-ca/internal/model"
	"cfar-ca/internal/stats"
)

// version 是 CLI 版本号。
const version = "1.0.0"

func main() {
	os.Exit(run(os.Args[1:]))
}

// run 派发子命令并返回进程退出码。
func run(args []string) int {
	if len(args) == 0 {
		usage()
		return 2
	}
	switch args[0] {
	case "detect":
		return cmdDetect(args[1:])
	case "sweep":
		return cmdSweep(args[1:])
	case "compare":
		return cmdCompare(args[1:])
	case "alpha":
		return cmdAlpha(args[1:])
	case "synth":
		return cmdSynth(args[1:])
	case "version":
		fmt.Printf("cfar-ca %s\n", version)
		return 0
	case "help", "-h", "--help":
		usage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "cfar-ca: 未知子命令 %q\n", args[0])
		usage()
		return 2
	}
}

// usage 打印到 stderr 的帮助信息。
func usage() {
	fmt.Fprintln(os.Stderr, `cfar-ca — 单元平均 CA-CFAR 检测核算（雷达检测内核，非安防产品）

用法:
  cfar-ca detect [-json] <spec.json>      逐单元阈值/检出/经验虚警（- 读 stdin）
  cfar-ca sweep -pfa "p1 p2 …" <spec.json>  Pfa 扫描：α 与检出数随 Pfa 变化
  cfar-ca compare -high 1e-3 -low 1e-4 <spec.json>  两个 Pfa 的交叉规则对照
  cfar-ca alpha -pfa 1e-3 -refs 8         打印 α = N(Pfa^{−1/N}−1)
  cfar-ca synth -n 4000 -guards 2 -refs 16 -pfa 1e-3 [-seed 7]  均匀噪声经验虚警带
  cfar-ca version                         打印版本

算例 JSON（见 example/）:
  {"amplitudes": [...], "guards": 2, "refs": 8, "pfa": 0.001}

约定:
  阈值 = α × 参考均值；α = N(Pfa^{−1/N}−1)，N=2×refs。
  保护单元与 CUT 不进参考均值；边缘 CUT 参考不足标为无效，不补零。
  非法输入打印到 stderr 并以非零码退出。`)
}

// parseFileArg 解析子命令的位置参数：恰好一个算例文件路径。
func parseFileArg(name string, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("子命令 %s 需要一个算例文件路径（或 - 表示 stdin）", name)
	}
	if len(args) > 1 {
		return "", fmt.Errorf("子命令 %s 只接受一个文件路径，多余参数: %v", name, args[1:])
	}
	return args[0], nil
}

// makeSynthConfig 构造 synth 子命令用的检测配置。
func makeSynthConfig(noise []float64, guards, refs int, pfa float64) *model.DetectorConfig {
	return model.NewDetectorConfig(noise, guards, refs, pfa)
}

// cmdDetect 执行 detect 子命令。
func cmdDetect(args []string) int {
	fs := flag.NewFlagSet("detect", flag.ContinueOnError)
	jsonOut := fs.Bool("json", false, "以 JSON 导出结果")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	path, err := parseFileArg("detect", fs.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 2
	}
	cfg, err := io.LoadSpecFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 1
	}
	res, err := detect.Detect(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 1
	}
	if *jsonOut {
		if err := io.WriteResultJSON(os.Stdout, res); err != nil {
			fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
			return 1
		}
		return 0
	}
	fmt.Print(io.FormatDetect(res))
	return 0
}

// cmdSweep 执行 sweep 子命令：对多个 Pfa 分别检测并对照。
func cmdSweep(args []string) int {
	fs := flag.NewFlagSet("sweep", flag.ContinueOnError)
	pfaList := fs.String("pfa", "1e-2 1e-3 1e-4", "以空白/逗号分隔的 Pfa 列表")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	path, err := parseFileArg("sweep", fs.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 2
	}
	pfas, err := io.ParsePfaList(*pfaList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 2
	}
	cfg, err := io.LoadSpecFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 1
	}
	sweep, err := detect.SweepPfa(cfg, pfas)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 1
	}
	fmt.Print(io.FormatSweep(sweep))
	return 0
}

// cmdCompare 执行 compare 子命令：两个 Pfa 的交叉规则对照。
func cmdCompare(args []string) int {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	high := fs.Float64("high", 1e-3, "较大的 Pfa")
	low := fs.Float64("low", 1e-4, "较小的 Pfa")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	path, err := parseFileArg("compare", fs.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 2
	}
	cfg, err := io.LoadSpecFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 1
	}
	rule, err := detect.ComparePfa(cfg, *high, *low)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 1
	}
	fmt.Print(io.FormatCrossRule(rule))
	return 0
}

// cmdAlpha 执行 alpha 子命令：打印给定 Pfa/参考单元数的放大系数。
func cmdAlpha(args []string) int {
	fs := flag.NewFlagSet("alpha", flag.ContinueOnError)
	pfa := fs.Float64("pfa", 1e-3, "名义虚警率")
	refs := fs.Int("refs", 8, "每侧参考单元数")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "cfar-ca: alpha 不接受位置参数: %v\n", fs.Args())
		return 2
	}
	a, err := detect.AlphaFor(*pfa, *refs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 1
	}
	fmt.Print(io.FormatAlpha(*pfa, *refs, a))
	return 0
}

// cmdSynth 执行 synth 子命令：生成均匀指数噪声，检测并报告经验虚警带。
func cmdSynth(args []string) int {
	fs := flag.NewFlagSet("synth", flag.ContinueOnError)
	n := fs.Int("n", 4000, "序列长度")
	guards := fs.Int("guards", 2, "保护单元数（每侧）")
	refs := fs.Int("refs", 16, "参考单元数（每侧）")
	pfa := fs.Float64("pfa", 1e-3, "名义虚警率")
	seed := fs.Int64("seed", 7, "随机种子（确定性复现）")
	mean := fs.Float64("mean", 1.0, "指数噪声均值")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "cfar-ca: synth 不接受位置参数: %v\n", fs.Args())
		return 2
	}
	noise, err := stats.Exponential(stats.NewSource(*seed), *n, *mean)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 1
	}
	cfg := makeSynthConfig(noise, *guards, *refs, *pfa)
	res, err := detect.Detect(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 1
	}
	st := res.Stats()
	band, err := stats.BandFor(st.ValidCells, *pfa, 3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfar-ca: %v\n", err)
		return 1
	}
	inBand := stats.InBand(st.EmpiricalPfa, band)
	fmt.Printf("均匀指数噪声（均值 %g, 种子 %d）经验虚警带校验\n", *mean, *seed)
	fmt.Printf("序列长 %d  保护 %d  每侧参考 %d  N=%d  α=%.6g  Pfa=%.6g\n",
		*n, *guards, *refs, 2**refs, res.Alpha, *pfa)
	fmt.Printf("有效 CUT %d   检出 %d   经验虚警 %.6g\n", st.ValidCells, st.DetectedCells, st.EmpiricalPfa)
	fmt.Printf("波动带（±3σ）[%.6g, %.6g]  带内: %v\n", band.Lo, band.Hi, inBand)
	if !inBand {
		fmt.Fprintf(os.Stderr, "cfar-ca: 经验虚警超出名义 Pfa 波动带\n")
		return 1
	}
	return 0
}
