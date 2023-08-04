package handlers

import (
	"expvar"
	"github.com/dimfeld/httptreemux/v5"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"os"
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

func APIMux(cfg APIMuxConfig) http.Handler {
	mux := httptreemux.NewContextMux()

	return mux
}
