package streamserver

import (
	"log"
	"net/http"
	"fmt"
	"os"
	"strconv"
	"time"
)

const FRAMEHEADER = "\r\n" +
	"--frame\r\n" +
	"Content-Type: image/jpeg\r\n" +
	"Content-Length: %d\r\n" +
	"X-Timestamp: 0.000000\r\n" +
	"\r\n"

type Server struct {
	host string
	port int
	framePerSecond int
	camera *Camera
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
		log.Println(err.Error())
		return nil, fmt.Errorf("Camera %d not initialized", cameraNum)
	}
	framePerSecondStr := getEnv("FRAME_PER_SECOND", "2")
	framePerSecond, err := strconv.Atoi(framePerSecondStr)
	if (err != nil) {
		return nil, fmt.Errorf("Error parse frame per second value %s", framePerSecondStr)
	}

	return &Server{
		host: host,
		port: port,
		framePerSecond: framePerSecond,
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
	log.Printf("Server is started on %s:%d host:port\n", s.host, s.port)
	http.Handle("/", s)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), nil)
}

func (s *Server) Terminate() {
	log.Println("Terminate server")
	if (s.camera != nil) {
		s.camera.Close()
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s : %s]", r.Method, r.RequestURI)
	switch r.RequestURI {
	case "/":
		s.handleIndex(w, r)
	case "/stream":
		s.handleStream(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) handleStream(w http.ResponseWriter, r *http.Request) {
	log.Printf("Connected - %s", r.RemoteAddr)
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary=frame")
	for {
		camImage, err := s.camera.GrabImage()
		if (err != nil) {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		camImageByte := camImage.Bytes()
		header := fmt.Sprintf(FRAMEHEADER, len(camImageByte))
		httpImageBody := make([]byte, (len(header) + len(camImageByte)) * 2)
		copy(httpImageBody, header)
		copy(httpImageBody[len(header):], camImageByte)
		if _, err := w.Write(httpImageBody); err != nil {
			log.Printf("Close -%s", r.RemoteAddr)
			break
		}
		time.Sleep(time.Duration(s.framePerSecond)* time.Second)
	}
}
