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
	"image"
	"net/http"
	"time"

	"github.com/mudler/gluedd/pkg/resource"
	"github.com/nfnt/resize"
	jpeg "github.com/pixiv/go-libjpeg/jpeg"
	live "github.com/saljam/mjpeg"
)

type JpegStreamerOptions struct {
	LiveStreamingURL string
	ListeningURL     string
	Stream           *live.Stream
	Buffer           int
	Width, Height    uint
	Resize, Approx   bool
	LivePreview      bool
}

func NewJpegStreamer(opts JpegStreamerOptions) resource.Resource {
	return &JpegStreamer{Options: opts}
}

type JpegStreamer struct {
	Options JpegStreamerOptions
}

func (l *JpegStreamer) Listen() chan string {

	files := make(chan string, l.Options.Buffer)
	if l.Options.LivePreview {
		go http.Handle("/", l.Options.Stream)
		go http.ListenAndServe(l.Options.ListeningURL, nil)
	}
	go func() {
		for {

			timeout := time.Duration(2 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}
			// Get the data
			resp, err := client.Get(l.Options.LiveStreamingURL)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
			defer resp.Body.Close()

			img, err := jpeg.DecodeIntoRGBA(resp.Body, &jpeg.DecoderOptions{})
			if err != nil {
				fmt.Println("Error decoding jpeg image: "+err.Error(), "[ERROR]")
				continue
			}

			var resizedImage image.Image
			if l.Options.Resize {
				if l.Options.Approx {
					resizedImage = resize.Thumbnail(l.Options.Width, l.Options.Height, img, resize.Lanczos3)
				} else {
					resizedImage = resize.Resize(l.Options.Width, l.Options.Height, img, resize.Lanczos3)
				}
			}

			// Encode as base64
			buffer64 := new(bytes.Buffer)
			if l.Options.Resize {
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
	}()

	return files

}
