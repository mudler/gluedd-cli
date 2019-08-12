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
	"os"
	"strconv"

	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/mjpeg"
	"github.com/mudler/gluedd/pkg/resource"
	jpeg "github.com/pixiv/go-libjpeg/jpeg"
	live "github.com/saljam/mjpeg"

	"net/http"
)

type V4lStreamerOptions struct {
	Device        int
	StreamURL     string
	Width, Height int
	Stream        *live.Stream
	Buffer        int
}

func NewV4lStreamer(opts *V4lStreamerOptions) resource.Resource {
	// Declare stream for web preview

	devicePath := "/dev/video" + strconv.Itoa(opts.Device)

	cam, err := v4l.Open(devicePath)
	if err != nil {
		fmt.Println("Error opening video capture device", "[FATAL]")
		os.Exit(1)
	}

	// Set camera properties
	// Fetch config
	cfg, err := cam.GetConfig()
	if err != nil {
		fmt.Println(err.Error(), "[ERROR]")
		os.Exit(1)
	}

	// Set parameters
	cfg.Format = mjpeg.FourCC
	cfg.Width = opts.Width
	cfg.Height = opts.Height

	// FIXME: This is a workaround - we have goroutine buffers but
	// we need to drop some fps or we will quickly start to detect every single frame.. and we will start to detect the past.

	cfg.FPS = v4l.Frac{N: uint32(1), D: uint32(opts.Buffer)} // Dummy - Make 1/Buff FPS to ensure that cam drops some frames to avoid congestion

	// Apply config
	err = cam.SetConfig(cfg)
	if err != nil {
		fmt.Println(err.Error(), "[ERROR]")
		os.Exit(1)
	}

	// Turn on cam
	err = cam.TurnOn()
	if err != nil {
		fmt.Println(err.Error(), "[ERROR]")
		os.Exit(1)
	}

	// Verify config
	cfg, err = cam.GetConfig()
	if err != nil {
		fmt.Println(err.Error(), "[ERROR]")
		os.Exit(1)
	}
	if cfg.Format != mjpeg.FourCC {
		fmt.Println("Failed to set MJPEG format.", "[ERROR]")
		os.Exit(1)
	}

	return &V4lStreamer{Options: opts, Cam: cam}
}

type V4lStreamer struct {
	Options *V4lStreamerOptions
	Cam     *v4l.Device
}

func (l *V4lStreamer) Listen() chan string {
	files := make(chan string, l.Options.Buffer)

	go http.Handle("/", l.Options.Stream)
	go http.ListenAndServe(l.Options.StreamURL, nil)
	go func() {
		for {
			// Read frame from camera
			buf, err := l.Cam.Capture()
			if err != nil {
				fmt.Println("Capture:", err)
				proc, _ := os.FindProcess(os.Getpid())
				proc.Signal(os.Interrupt)
				break
			}

			// Decode frame to jpeg
			//buf.Seek(0, 0)
			img, err := jpeg.DecodeIntoRGBA(buf, &jpeg.DecoderOptions{})
			if err != nil {
				fmt.Println("Error decoding jpeg image: "+err.Error(), "[ERROR]")
				os.Exit(1)
			}

			// Encode as base64
			buffer64 := new(bytes.Buffer)
			err = jpeg.Encode(buffer64, img, &jpeg.EncoderOptions{Quality: 100})
			if err != nil {
				fmt.Println("Error encoding image to base64: " + err.Error())
				os.Exit(1)
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
