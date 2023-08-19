package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	v1CheckGrp "github.com/Avyukth/service3-clone/app/services/sales-api/handlers/debug/checkgrp"
	v1TestGrp "github.com/Avyukth/service3-clone/app/services/sales-api/handlers/v1/testgrp"
	v1UserGrp "github.com/Avyukth/service3-clone/app/services/sales-api/handlers/v1/usergrp"
	userCore "github.com/Avyukth/service3-clone/business/core/user"
	"github.com/Avyukth/service3-clone/business/sys/auth"
	"github.com/Avyukth/service3-clone/business/web/mid"
	"github.com/Avyukth/service3-clone/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
	DB       *sqlx.DB
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

func DebugMux(build string, log *zap.SugaredLogger, db *sqlx.DB) http.Handler {

	mux := DebugStandardLibraryMux()

	cgh := v1CheckGrp.Handlers{
		Build: build,
		Log:   log,
		DB:    db,
	}

	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)
	return mux
}

func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)

	//
	v1(app, cfg)
	return app
}

func v1(app *web.App, cfg APIMuxConfig) {

	const version = "v1"
	tgh := v1TestGrp.Handlers{
		Log: cfg.Log,
	}

	app.Handle(http.MethodGet, version, "/test", tgh.Test)
	app.Handle(http.MethodGet, version, "/testauth", tgh.Test, mid.Authenticate(cfg.Auth), mid.Authorize(auth.RoleAdmin))

	ugh := v1UserGrp.Handlers{
		User: userCore.NewCore(cfg.Log, cfg.DB),
		Auth: cfg.Auth,
	}
	app.Handle(http.MethodGet, version, "/users/token", ugh.Token)
	app.Handle(http.MethodGet, version, "/users/:page/:rows", ugh.Query, mid.Authenticate(cfg.Auth))
	app.Handle(http.MethodGet, version, "/users/:id", ugh.QueryByID, mid.Authenticate(cfg.Auth), mid.Authorize(auth.RoleAdmin))

	app.Handle(http.MethodPost, version, "/users", ugh.Create, mid.Authenticate(cfg.Auth), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodPut, version, "/users/:id", ugh.Update, mid.Authenticate(cfg.Auth), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodDelete, version, "/users/:id", ugh.Delete, mid.Authenticate(cfg.Auth), mid.Authorize(auth.RoleAdmin))
}
