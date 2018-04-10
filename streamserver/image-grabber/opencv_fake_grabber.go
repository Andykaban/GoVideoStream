// +build without_opencv

package image_grabber

import (
	"errors"
	"fmt"
	"log"
)

func NewOpenCVCamera(camNum int) (ImageGrabber, error) {
	log.Println("Try to init web camera with OpenCV...")
	errorMessage := fmt.Sprintf("disable init web camera with %d number", camNum)
	return nil, errors.New(errorMessage)
}
