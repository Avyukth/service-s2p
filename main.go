package main

import (
	"go.uber.org/automaxprocs"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	// "github.com/ardanlabs/conf"
	// "github.com/dimfeld/httptreemux/v5"
)

var build = "develop"

func main() {
	g := runtime.GOMAXPROCS(0)
	log.Printf("Starting the Service build[%s] CPU[%d].............", build, g)
	defer log.Println("Service Ended")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Panicln("stopping Service ........")
}
