package server

import "net/http"

func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/detect", handleDetect)
	mux.HandleFunc("/api/sweep", handleSweep)
	mux.HandleFunc("/api/compare", handleCompare)
	mux.HandleFunc("/api/alpha", handleAlpha)
	mux.HandleFunc("/api/info", handleInfo)
	mux.HandleFunc("/api/endpoints", handleEndpoints)
	return mux
}

func ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, New())
}
