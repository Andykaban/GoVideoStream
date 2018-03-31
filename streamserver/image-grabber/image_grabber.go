package image_grabber

type ImageGrabber interface {
	GrabImage() ([]byte, error)
	Close() error
}
