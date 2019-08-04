package types

import (
	"github.com/mudler/gluedd/pkg/resource"
)

func NewDummy() resource.Resource {
	return &DummyResource{}
}

// fs listener
type DummyResource struct {
}

func (l *DummyResource) Listen() chan string {

	files := make(chan string)
	go func() {
		for {

			files <- "http://192.168.1.182:4000/file/3acc.jpeg"
			files <- "http://192.168.1.182:4000/file/photo.jpg"

		}
	}()

	return files

}
