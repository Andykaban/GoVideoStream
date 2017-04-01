package streamserver

import (
	"log"
	"net/http"
	"fmt"
	"os"
	"strconv"
	"time"
	"sync"
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
	clients map[chan []byte] bool
	framePerSecond int
	camera *Camera
	mutex *sync.Mutex
	currentFrame []byte
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
		clients: make(map[chan []byte] bool),
		framePerSecond: framePerSecond,
		camera: camera,
		mutex: &sync.Mutex{},
		currentFrame: make([]byte, len(FRAMEHEADER)),
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
	s.currentFrameUpdater()
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), nil)
}

func (s *Server) Terminate() {
	log.Println("Terminate server")
	if (s.camera != nil) {
		s.mutex.Lock()
		defer s.mutex.Unlock()
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
	streamCh := make(chan []byte)
	s.mutex.Lock()
	s.clients[streamCh] = true
	s.mutex.Unlock()
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary=frame")
	for {
		httpImageBody := <- streamCh
		if _, err := w.Write(httpImageBody); err != nil {
			log.Printf("Close -%s", r.RemoteAddr)
			break
		}
		time.Sleep(time.Duration(s.framePerSecond)* time.Second)
	}
	s.mutex.Lock()
	delete(s.clients, streamCh)
	s.mutex.Unlock()
}

func (s *Server) currentFrameUpdater() {
	go func() {
		for {
			s.mutex.Lock()
			clients_count := len(s.clients)
			s.mutex.Unlock()
			if (clients_count == 0) {
				continue
			}
			camImage, err := s.camera.GrabImage()
			if (err != nil) {
				log.Println(err.Error())
				break
			}
			camImageByte := camImage.Bytes()
			header := fmt.Sprintf(FRAMEHEADER, len(camImageByte))
			if (len(s.currentFrame) < (len(header) + len(camImageByte))) {
				s.currentFrame = make([]byte, (len(header) + len(camImageByte) * 2))
			}
			copy(s.currentFrame, header)
			copy(s.currentFrame[len(header):], camImageByte)

			s.mutex.Lock()
			for ch := range s.clients {
				select {
				case ch <- s.currentFrame:
				default:
				}
			}
			s.mutex.Unlock()
		}
	}()
}
