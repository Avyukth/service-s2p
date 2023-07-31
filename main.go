package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/automaxprocs"
	"go.uber.org/automaxprocs/maxprocs"
	// "github.com/ardanlabs/conf"
	// "github.com/dimfeld/httptreemux/v5"
)

var build = "develop"

func main() {
	if _, err := maxprocs.Set(); err != nil {
		fmt.Println("maxprocs: %w", err)
		os.Exit(1)
	}
	g := runtime.GOMAXPROCS(0)
	log.Printf("Starting the Service build[%s] CPU[%d].............", build, g)
	defer log.Println("Service Ended")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Panicln("stopping Service ........")
}
