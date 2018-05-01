package streamserver

import (
	"log"
	"net/http"
	"fmt"
	"time"
	"sync"
	"github.com/Andykaban/GoVideoStream/streamserver/image-grabber"
)

const FRAMEHEADER = "\r\n" +
	"--frame\r\n" +
	"Content-Type: image/jpeg\r\n" +
	"Content-Length: %d\r\n" +
	"X-Timestamp: 0.000000\r\n\r\n"

type Server struct {
	host string
	port int
	clients map[chan []byte] bool
	framePerSecond int
	mutex *sync.Mutex
	currentFrame []byte
	webcamMode string
	webcamNumber int
	TerminateServerChannel chan bool
}

func New(host string, port int, frameDelay int, camMode string, camNum int) (*Server) {
	return &Server{
		host: host,
		port: port,
		clients: make(map[chan []byte] bool),
		framePerSecond: frameDelay,
		mutex: &sync.Mutex{},
		currentFrame: make([]byte, len(FRAMEHEADER)),
		webcamMode: camMode,
		webcamNumber: camNum,
		TerminateServerChannel: make(chan bool),
	}
}

func (s *Server) GetCamera() (image_grabber.ImageGrabber, error) {
	if s.webcamMode == "opencv" {
		camera, err := image_grabber.NewOpenCVCamera(s.webcamNumber)
		if err != nil {
			log.Println(err.Error())
			return nil, fmt.Errorf("camera %d not initialized by opencv", s.webcamNumber)
		}
		return camera, nil
	} else if s.webcamMode == "v4l" {
		camera, err := image_grabber.NewV4LCamera(s.webcamNumber)
		if err != nil {
			return nil, fmt.Errorf("camera %d not initialized by v4l", s.webcamNumber)
		}
		return camera, nil
	} else {
		return nil, fmt.Errorf("unsupported webcam mode %s", s.webcamMode)
	}
}

func (s *Server) Run() (error) {
	camera, err := s.GetCamera()
	if err != nil {
		panic(err)
	}
	log.Printf("Server is started on %s:%d host:port\n", s.host, s.port)
	http.Handle("/", s)
	s.currentFrameUpdater(camera)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), nil)
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
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Transfer-Encoding", "chunked")
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary=frame")
	w.Header().Add("Expires", "-1")
	w.Header().Add("Pragma", "no-cache")

	for {
		httpImageBody := <- streamCh
		if _, err := w.Write(httpImageBody); err != nil {
			log.Printf("Close - %s", r.RemoteAddr)
			break
		}
		time.Sleep(time.Duration(s.framePerSecond)* time.Second)
	}
	s.mutex.Lock()
	delete(s.clients, streamCh)
	s.mutex.Unlock()
}

func (s *Server) currentFrameUpdater(workCam image_grabber.ImageGrabber) {
	go func() {
		for {
			s.mutex.Lock()
			clientsCount := len(s.clients)
			s.mutex.Unlock()
			if clientsCount == 0 {
				continue
			}
			camImage, err := workCam.GrabImage()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			header := fmt.Sprintf(FRAMEHEADER, len(camImage))
			if len(s.currentFrame) < (len(header) + len(camImage)) {
				s.currentFrame = make([]byte, len(header) + len(camImage) * 2)
			}
			copy(s.currentFrame, header)
			copy(s.currentFrame[len(header):], camImage)

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
	go func() {
		<- s.TerminateServerChannel
		log.Println("Try to stop webcam")
		s.mutex.Lock()
		defer s.mutex.Unlock()
		workCam.Close()
		s.TerminateServerChannel <- true
	}()
}
