package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"cfar-ca/internal/detect"
	"cfar-ca/internal/io"
	"cfar-ca/internal/model"
)

type sweepRequest struct {
	Spec json.RawMessage `json:"spec"`
	Pfas []float64       `json:"pfas"`
}

type compareRequest struct {
	Spec json.RawMessage `json:"spec"`
	High float64         `json:"high"`
	Low  float64         `json:"low"`
}

type alphaRequest struct {
	Pfa  float64 `json:"pfa"`
	Refs int     `json:"refs"`
}

type sweepRow struct {
	Pfa        float64 `json:"pfa"`
	Alpha      float64 `json:"alpha"`
	ValidCells int     `json:"valid_cells"`
	Detected   int     `json:"detected"`
	Empirical  float64 `json:"empirical"`
	MaxMargin  float64 `json:"max_margin"`
}

type endpointDoc struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Doc    string `json:"doc"`
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, map[string]string{
		"service": "cfar-ca",
		"detect":  "POST /api/detect",
		"sweep":   "POST /api/sweep",
		"compare": "POST /api/compare",
		"alpha":   "POST /api/alpha",
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}

func handleDetect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cfg, err := specFromBody(r)
	if err != nil {
		writeError(w, err)
		return
	}
	res, err := detect.Detect(cfg)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := io.WriteResultJSON(w, res); err != nil {
		writeError(w, err)
		return
	}
}

func handleSweep(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req sweepRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}
	cfg, err := parseSpec(req.Spec)
	if err != nil {
		writeError(w, err)
		return
	}
	pfas := req.Pfas
	if len(pfas) == 0 {
		pfas = []float64{1e-2, 1e-3, 1e-4}
	}
	sw, err := detect.SweepPfa(cfg, pfas)
	if err != nil {
		writeError(w, err)
		return
	}
	rows := make([]sweepRow, 0, len(sw.Entries))
	for _, e := range sw.Entries {
		rows = append(rows, sweepRow{
			Pfa:        e.Pfa,
			Alpha:      e.Alpha,
			ValidCells: e.ValidCells,
			Detected:   e.Detected,
			Empirical:  e.Empirical,
			MaxMargin:  e.MaxMargin,
		})
	}
	writeJSON(w, rows)
}

func handleCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req compareRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}
	cfg, err := parseSpec(req.Spec)
	if err != nil {
		writeError(w, err)
		return
	}
	rule, err := detect.ComparePfa(cfg, req.High, req.Low)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]any{
		"pfa_high":    rule.PfaHigh,
		"pfa_low":     rule.PfaLow,
		"alpha_high":  rule.AlphaHigh,
		"alpha_low":   rule.AlphaLow,
		"det_high":    rule.DetHigh,
		"det_low":     rule.DetLow,
		"alpha_rises": rule.AlphaRises,
		"det_falls":   rule.DetFalls,
	})
}

func handleAlpha(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req alphaRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}
	a, err := detect.AlphaFor(req.Pfa, req.Refs)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]any{"alpha": a, "pfa": req.Pfa, "refs": req.Refs})
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, map[string]any{
		"name":    "cfar-ca",
		"version": "1.0.0",
		"routes":  []string{"/api/detect", "/api/sweep", "/api/compare", "/api/alpha", "/health"},
	})
}

func handleEndpoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, []endpointDoc{
		{Method: "POST", Path: "/api/detect", Doc: "CA-CFAR 检测"},
		{Method: "POST", Path: "/api/sweep", Doc: "Pfa 扫描"},
		{Method: "POST", Path: "/api/compare", Doc: "交叉规则对照"},
		{Method: "POST", Path: "/api/alpha", Doc: "放大系数"},
		{Method: "GET", Path: "/health", Doc: "健康检查"},
	})
}

func specFromBody(r *http.Request) (*model.DetectorConfig, error) {
	var body map[string]json.RawMessage
	if err := decodeJSON(r, &body); err != nil {
		return nil, err
	}
	data, ok := body["spec"]
	if !ok {
		return nil, errors.New("请求体缺少 spec 字段")
	}
	return parseSpec(data)
}

func parseSpec(data []byte) (*model.DetectorConfig, error) {
	cfg, err := io.ParseSpecBytes(data)
	if err != nil {
		return nil, err
	}
	if err := model.Validate(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func decodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
