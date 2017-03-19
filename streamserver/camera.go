package streamserver

import (
	"github.com/Andykaban/go-opencv/opencv"
	"sync"
	"log"
	"fmt"
	"bytes"
	"image/jpeg"
)

type Camera struct {
	cap *opencv.Capture
	mutex *sync.Mutex
}

func New(camNum int) (c *Camera, err error) {
	log.Println("Try to init web camera...")
	cap := opencv.NewCameraCapture(camNum)
	if (cap == nil) {
		return nil, fmt.Errorf("Error init web camera with %d number", camNum)
	}

	return &Camera{
		cap: cap,
		mutex: &sync.Mutex{},
	}, nil
}

func (c *Camera) GrabImage() (grabImage *bytes.Buffer, err error){
	log.Println("Try to grab image...")
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.cap.GrabFrame() {
		img := c.cap.RetrieveFrame(0)
		if (img != nil) {
			imageBuffer := new(bytes.Buffer)
			convertImage := img.ToImage()
			err = jpeg.Encode(imageBuffer, convertImage, nil)
			if err != nil {
				return nil, fmt.Errorf("Grabbed image not encoded to jpeg format")
			}
			return imageBuffer, nil
		} else {
			return nil, fmt.Errorf("Grabbed image not retreaved from web camera")
		}
	} else {
		return nil, fmt.Errorf("Image not grabbed from web camera")
	}
}

func (c *Camera) Close() (err error) {
	if (c.cap != nil) {
		log.Println("Release web camera...")
		c.cap.Release()
	}
	return nil
}
