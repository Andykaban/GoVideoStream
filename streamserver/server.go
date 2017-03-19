package streamserver

import (
	"log"
	"net/http"
	"fmt"
	"os"
	"strconv"
)

type Server struct {
	host string
	port int
	camera Camera
}

func New() (s *Server, err error) {
	host := getEnv("HOST", "127.0.0.1")
	portStr := getEnv("PORT", "1488")
	port, err := strconv.Atoi(portStr)
	if (err != nil) {
		return nil, fmt.Errorf("Error parse port number %s", portStr)
	}
	cameraNumStr := getEnv("CAMERA_NUMBER", "0")
	cameraNum, err := strconv.Atoi(cameraNumStr)
	if (err != nil) {
		return nil, fmt.Errorf("Error parse camera number %s", cameraNumStr)
	}
	camera, err := CameraInit(cameraNum)
	if (err != nil) {
		log.Println(err)
		return nil, fmt.Errorf("Camera %d not initialized", cameraNum)
	}

	return &Server{
		host: host,
		port: port,
		camera: camera,
	}, nil
}

func getEnv(envName string, defVal string) (val string) {
	val = os.Getenv(envName)
	if (val == "") {
		val = defVal
	}

	return val
}

func (s *Server) Run() (err error) {
	log.Printf("Start server on %s:%d\n", s.host, s.port)
	defer s.camera.Close()
	http.Handle("/", s)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), nil)
}
