// +build with_opencv

package image_grabber

import (
	"gocv.io/x/gocv"
	"sync"
	"log"
	"fmt"
)

type OpenCVCamera struct {
	webcam *gocv.VideoCapture
	mutex *sync.Mutex
}

func NewOpenCVCamera(camNum int) (ImageGrabber, error) {
	log.Println("Try to init web camera with OpenCV...")
	webcam, err := gocv.VideoCaptureDevice(int(camNum))
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Error init web camera with %d number", camNum)
	}

	return &OpenCVCamera{
		webcam: webcam,
		mutex: &sync.Mutex{},
	}, nil
}

func (c *OpenCVCamera) GrabImage() ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	img := gocv.NewMat()
	defer img.Close()
	if ok := c.webcam.Read(img); !ok {
		return nil, fmt.Errorf("Grabbed image not retreaved from web camera")
	}
	buf, err := gocv.IMEncode(".jpg", img)
	if err !=nil {
		return nil, fmt.Errorf("Grabbed image not encoded to jpeg format")
	}
	return buf, nil
}

func (c *OpenCVCamera) Close() (error) {
	if c.webcam != nil {
		log.Println("Release OpenCV web camera...")
		c.webcam.Close()
	}
	return nil
}
