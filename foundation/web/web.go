package web

import (
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	return &App{ContextMux: httptreemux.NewContextMux(), shutdown: shutdown}
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}
