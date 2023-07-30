package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	// "github.com/ardanlabs/conf"
	// "github.com/dimfeld/httptreemux/v5"
)

var build = "develop"

func main() {
	log.Println("Starting Service .............", build)
	defer log.Println("Service Ended")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Panicln("stopping Service ........")

}
