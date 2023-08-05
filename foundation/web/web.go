package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/dimfeld/httptreemux/v5"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Middleware) {

	handler = wrapMiddleware(mw, handler)

	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}

		ctx = context.WithValue(ctx, key, &v)

		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}

	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	a.ContextMux.Handle(method, finalPath, h)
}
