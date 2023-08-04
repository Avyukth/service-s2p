package handlers

import (
	"encoding/json"
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/Avyukth/service3-clone/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/dimfeld/httptreemux/v5"
	"go.uber.org/zap"
)

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())
	return mux
}

func DebugMux(build string, log *zap.SugaredLogger) http.Handler {

	mux := DebugStandardLibraryMux()

	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}

	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)
	return mux
}

func APIMux(cfg APIMuxConfig) http.Handler {
	mux := httptreemux.NewContextMux()
	h := func(w http.ResponseWriter, r *http.Request) {
		status := struct {
			Status string `json:"status"`
		}{
			Status: "OK",
		}
		json.NewEncoder(w).Encode(status)
	}
	mux.Handle(http.MethodGet, "/test", h)
	return mux
}
