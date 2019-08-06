package types

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/mjpeg"
	"github.com/mudler/gluedd/pkg/resource"
	jpeg "github.com/pixiv/go-libjpeg/jpeg"
	live "github.com/saljam/mjpeg"
	"os"
	"strconv"

	"net/http"
)

func NewV4lStreamer(device int, StreamURL string, Width, Height int, stream *live.Stream, Buffer int) resource.Resource {
	// Declare stream for web preview

	devicePath := "/dev/video" + strconv.Itoa(device)

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
	cfg.Width = Width
	cfg.Height = Height

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

	return &V4lStreamer{DeviceID: device, StreamURL: StreamURL, Cam: cam, Stream: stream, Buffer: Buffer}
}

type V4lStreamer struct {
	DeviceID  int
	StreamURL string
	Cam       *v4l.Device
	Stream    *live.Stream
	Buffer    int
}

func (l *V4lStreamer) Listen() chan string {
	files := make(chan string, l.Buffer)

	go http.Handle("/", l.Stream)
	go http.ListenAndServe(l.StreamURL, nil)
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

			if l.Buffer != 0 {
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
