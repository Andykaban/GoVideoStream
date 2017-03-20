package main

import (
	"./streamserver"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	log.Println("Start Video Streaming server...")
	server, err := streamserver.New()
	if (err != nil) {
		panic(err)
	}
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<- ch
		server.Terminate()
		os.Exit(0)
	}()
	server.Run()
}
