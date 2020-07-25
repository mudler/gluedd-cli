// Copyright © 2019 Ettore Di Giacinto <mudler@gentoo.org>
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see <http://www.gnu.org/licenses/>.

package types

import (
	"bytes"
	"encoding/base64"
	"fmt"
	mjpeg "github.com/marpie/go-mjpeg"
	"github.com/mudler/gluedd/pkg/resource"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	jpeg "github.com/pixiv/go-libjpeg/jpeg"
	live "github.com/saljam/mjpeg"
	"image"
	"io"
	"net/http"
)

type MJpegStreamerOptions struct {
	LiveStreamingURL                                                 string
	ListeningURL, CropMode                                           string
	Stream                                                           *live.Stream
	Buffer, Timeout, CropWidth, CropHeight, CropAnchorX, CropAnchorY int
	Width, Height                                                    uint
	Resize, Approx, Crop, CropAnchor                                 bool
	LivePreview                                                      bool
}

func NewMJpegStreamer(opts MJpegStreamerOptions) resource.Resource {
	return &MJpegStreamer{Options: opts}
}

type MJpegStreamer struct {
	Options MJpegStreamerOptions
}

// processHttp receives the HTTP data and tries to decodes images. The images
// are sent through a chan for further processing.
func processHttp(response *http.Response, nextImg chan *image.Image, quit chan bool) {
	defer response.Body.Close()
	for {
		select {
		case <-quit:
			close(nextImg)
			return
		default:
			img, err := mjpeg.Decode(response.Body)
			if err == io.EOF {
				close(nextImg)
				return
			}
			if err != nil {
				fmt.Println(err)
			}
			if img != nil {
				nextImg <- img
			}
		}
	}
}

// processImage receives images through a chan and prints the dimensions.
func (l *MJpegStreamer) processImage(nextImg chan *image.Image, quit chan bool, files chan string) {
	for i := range nextImg {

		img := *i
		if *i == nil {
			continue
		}

		var resizedImage image.Image
		var err error
		if l.Options.Resize {
			if l.Options.Approx {
				resizedImage = resize.Thumbnail(l.Options.Width, l.Options.Height, img, resize.Lanczos3)
			} else {
				resizedImage = resize.Resize(l.Options.Width, l.Options.Height, img, resize.Lanczos3)
			}
		}

		if l.Options.Crop {
			cfg := cutter.Config{
				Width:  l.Options.CropWidth,
				Height: l.Options.CropHeight,
			}
			switch l.Options.CropMode {
			case "centered":
				cfg.Mode = cutter.Centered
			case "top_left":
				cfg.Mode = cutter.TopLeft
			default:
				cfg.Mode = cutter.Centered
			}
			if l.Options.CropAnchor {
				cfg.Anchor = image.Point{l.Options.CropAnchorX, l.Options.CropAnchorY}
			}
			if l.Options.Resize {
				resizedImage, err = cutter.Crop(resizedImage, cfg)
			} else {
				resizedImage, err = cutter.Crop(img, cfg)
			}
			if err != nil {
				fmt.Println("Error cropping jpeg image: "+err.Error(), "[ERROR]")
				continue
			}
		}

		// Encode as base64
		buffer64 := new(bytes.Buffer)
		if l.Options.Resize || l.Options.Crop {
			err = jpeg.Encode(buffer64, resizedImage, &jpeg.EncoderOptions{Quality: 100})
		} else {
			err = jpeg.Encode(buffer64, img, &jpeg.EncoderOptions{Quality: 100})
		}
		if err != nil {
			fmt.Println("Error encoding image to base64: " + err.Error())
			continue
		}
		imageBase64 := base64.StdEncoding.EncodeToString(buffer64.Bytes())
		if l.Options.Buffer != 0 {

			select {
			case files <- imageBase64:
			default:
				// If API is slow drop frames
			}
		} else {
			files <- imageBase64
		}

	}
	quit <- true
}

func (l *MJpegStreamer) Listen() chan string {

	files := make(chan string, l.Options.Buffer)
	if l.Options.LivePreview {
		go http.Handle("/", l.Options.Stream)
		go http.ListenAndServe(l.Options.ListeningURL, nil)
	}
	response, err := http.Get(l.Options.LiveStreamingURL)
	if err != nil {
		return files
	}
	nextImg := make(chan *image.Image, 30)
	quit := make(chan bool)
	fmt.Println("Waiting for images to process...")
	go l.processImage(nextImg, quit, files)
	go processHttp(response, nextImg, quit)

	return files

}
