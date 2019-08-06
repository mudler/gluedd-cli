package types

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/mudler/gluedd/pkg/resource"
	jpeg "github.com/pixiv/go-libjpeg/jpeg"
	live "github.com/saljam/mjpeg"
	"net/http"
	"time"
)

func NewJpegStreamer(url string, baseurl string, stream *live.Stream, live bool, Buffer int) resource.Resource {
	return &JpegStreamer{StreamUrl: url, BaseUrl: baseurl, Stream: stream, Live: live, Buffer: Buffer}
}

type JpegStreamer struct {
	StreamUrl string
	BaseUrl   string
	Stream    *live.Stream
	Live      bool
	Buffer    int
}

func (l *JpegStreamer) Listen() chan string {

	files := make(chan string, l.Buffer)
	if l.Live {
		go http.Handle("/", l.Stream)
		go http.ListenAndServe(l.BaseUrl, nil)
	}
	go func() {
		for {

			timeout := time.Duration(2 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}
			// Get the data
			resp, err := client.Get(l.StreamUrl)
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

			// Encode as base64
			buffer64 := new(bytes.Buffer)
			err = jpeg.Encode(buffer64, img, &jpeg.EncoderOptions{Quality: 100})
			if err != nil {
				fmt.Println("Error encoding image to base64: " + err.Error())
				continue
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
