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

func NewV4LCamera(camNum int) (ImageGrabber, error) {
	camPath := fmt.Sprintf("/dev/video%d", camNum)
	cam, err := webcam.Open(camPath)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error init web camera with %s path", camPath)
	}
	format, _ := getWebcamFormatByString("MJPG")
	cam.SetImageFormat(format, 640, 480)
	err = cam.StartStreaming()
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error switch web camera to stream mode")
	}
	return &V4LGrabber{
		cam: cam,
		mutex: &sync.Mutex{},
	}, nil
}

func getWebcamFormatByString(formatStr string) (webcam.PixelFormat, error) {
	chars := []rune(formatStr)
	if len(chars) != 4 {
		return 0, fmt.Errorf("Format string its to long")
	}
	p := chars[0:4]
	formatCode := uint32(p[0]) | (uint32(p[1])<<8) | (uint32(p[2])<<16) | (uint32(p[3])<<24)
	ret := webcam.PixelFormat(formatCode)
	return ret, nil
}

func (c *V4LGrabber) GrabImage() ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.cam.WaitForFrame(1)

	switch err.(type) {
	case nil:
	case *webcam.Timeout:
		return nil, fmt.Errorf("web Camera timeout")
	default:
		log.Println(err)
		return nil, err
	}

	imageJpeg, err := c.cam.ReadFrame()
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("grabbed image not retreaved from web camera")
	}
	if len(imageJpeg) == 0 {
		return nil, fmt.Errorf("grabbed image not encoded to jpeg format")
	}
	return imageJpeg, nil
}

func (c *V4LGrabber) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
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
