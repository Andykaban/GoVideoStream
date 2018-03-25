package main

import (
	"log"
	"os"
	"flag"
	"os/signal"
	"syscall"
	"github.com/Andykaban/GoVideoStream/streamserver"
	"fmt"
)

var (
	host = flag.String("host", "0.0.0.0", "")
	port = flag.Int("port", 1488, "")
	frameDelay = flag.Int("frame_delay", 1, "")
	webcamMode = flag.String("webcam_mode", "opencv", "")
	cameraNumber = flag.Int("camera_number", 0, "")
)

func Usage() {
	fmt.Println("Video stream server")
	flag.PrintDefaults()
}

func main()  {
	flag.Usage = Usage
	flag.Parse()
	log.Println("Start Video Streaming server...")
	server := streamserver.New(*host, *port, *frameDelay, *webcamMode, *cameraNumber)
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<- ch
		server.TerminateServerChannel <-true
		<- server.TerminateServerChannel
		log.Println("Video Streaming server terminated")
		os.Exit(0)
	}()
	server.Run()
}
