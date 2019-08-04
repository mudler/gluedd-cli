package types

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/mudler/gluedd/pkg/resource"
)

func NewLocalWatcher(dir string) resource.Resource {
	return &LocalResource{dir}
}

// fs listener
type LocalResource struct {
	Directory string
}

func (l *LocalResource) Listen() chan string {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	files := make(chan string)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println(event.Name)
					files <- event.Name
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Fatal(err)
					return
				}
			}
		}
	}()

	err = watcher.Add(l.Directory)
	if err != nil {
		log.Fatal(err)
	}
	return files

}
