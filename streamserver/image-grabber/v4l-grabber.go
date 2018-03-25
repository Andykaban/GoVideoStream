package image_grabber

import (
	"sync"
	"log"
	"fmt"
	"github.com/blackjack/webcam"
)

type V4LGrabber struct {
	cam *webcam.Webcam
	mutex *sync.Mutex
}

func NewV4LCamera(camNum int) (ImageGrabber, error){
	camPath := fmt.Sprintf("/dev/video%d", camNum)
	cam, err := webcam.Open(camPath)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Error init web camera with %s path", camPath)
	}
	err = cam.StartStreaming()
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Error switch web camera to stream mode")
	}
	return &V4LGrabber{
		cam: cam,
		mutex: &sync.Mutex{},
	}, nil
}

func (c *V4LGrabber) GrabImage() ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	imageJpeg, err := c.cam.ReadFrame()
	if err != nil {
		return nil, fmt.Errorf("Grabbed image not retreaved from web camera")
	}
	if len(imageJpeg) == 0 {
		return nil, fmt.Errorf("Grabbed image not encoded to jpeg format")
	}
	return imageJpeg, nil
}

func (c *V4LGrabber) Close() error {
	if c.cam != nil {
		err := c.cam.StopStreaming()
		if err != nil {
			log.Println(err)
		}
		log.Println("Release web camera...")
		c.cam.Close()
	}
	return nil
}
