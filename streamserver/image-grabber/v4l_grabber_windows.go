package image_grabber

import (
	"fmt"
	"errors"
)

func NewV4LCamera(camNum int) (ImageGrabber, error) {
	camPath := fmt.Sprintf("/dev/video%d", camNum)
	fmt.Printf("Try to open %s", camPath)
	return nil, errors.New("Not implemented for Windows")
}
