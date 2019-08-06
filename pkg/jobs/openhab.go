package generators

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
	jpeg "github.com/pixiv/go-libjpeg/jpeg"
	live "github.com/saljam/mjpeg"
	"image"
	"image/draw"
	"net/http"
	"strings"
	"time"
)

type OpenHabErrand struct {
	Prediction             api.Prediction
	URL                    string
	VehicleItem, HumanItem string
	Stream                 *live.Stream
	Live                   bool
}

func (e *OpenHabErrand) UpdateItem(item, data string) error {
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	fmt.Println("Updating", item, data)
	req, err := http.NewRequest("POST", e.URL+"/rest/items/"+item, bytes.NewBufferString(data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (e *OpenHabErrand) Apply() error {
	if e.Prediction.Error != nil || len(e.Prediction.Body.Predictions) == 0 {
		return e.Prediction.Error
	}

	person := false
	vehicle := false
	for _, c := range e.Prediction.Body.Predictions[0].Classes {
		if strings.Contains(c.Cat, "Van") || strings.Contains(c.Cat, "Truck") || strings.Contains(c.Cat, "Car") {
			vehicle = true
		}
		if strings.Contains(c.Cat, "Man") {
			person = true
		}

		if strings.Contains(c.Cat, "Woman") {
			person = true
		}

		if strings.Contains(c.Cat, "Person") || strings.Contains(c.Cat, "Face") {
			person = true
		}
	}
	go func() {
		if vehicle {
			e.UpdateItem(e.VehicleItem, "ON")
		} else {
			e.UpdateItem(e.VehicleItem, "OFF")
		}
		if person {
			e.UpdateItem(e.HumanItem, "ON")
		} else {
			e.UpdateItem(e.HumanItem, "OFF")
		}
	}()
	if e.Live {
		go func() {

			unbased, err := base64.StdEncoding.DecodeString(e.Prediction.Url)
			if err != nil {
				return
			}
			img, err := jpeg.Decode(bytes.NewReader(unbased), &jpeg.DecoderOptions{})
			if err != nil {
				return
			}
			if len(e.Prediction.Body.Predictions) != 0 {
				// Convert to RGBA
				b := img.Bounds()
				imgRGBA := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
				draw.Draw(imgRGBA, imgRGBA.Bounds(), img, b.Min, draw.Src)

				for i, _ := range e.Prediction.Body.Predictions[0].Classes {
					imgRGBA = writeBoundingBox(img, e.Prediction, i)
					img = imgRGBA
				}

				buf := new(bytes.Buffer)
				if imgRGBA != nil {
					err := jpeg.Encode(buf, imgRGBA, &jpeg.EncoderOptions{Quality: 50})
					if err == nil {
						go e.Stream.UpdateJPEG(buf.Bytes())
					}
				}

			}
		}()
	}
	return nil
}

// Generate makes the errand generates new prediction if needed
func (e *OpenHabErrand) Generate(d api.Detector) *api.Prediction {
	return nil
}

type OpenHabGenerator struct {
	URL                    string
	VehicleItem, HumanItem string
	Stream                 *live.Stream
	Live                   bool
}

func NewOpenHabGenerator(OpenHabURL, VehicleItem, HumanItem string, stream *live.Stream, live bool) errand.ErrandGenerator {
	return &OpenHabGenerator{URL: OpenHabURL, VehicleItem: VehicleItem, HumanItem: HumanItem, Stream: stream, Live: live}
}
func (l *OpenHabGenerator) GenerateErrand(p api.Prediction) errand.Errand {
	return &OpenHabErrand{Prediction: p, URL: l.URL, VehicleItem: l.VehicleItem, HumanItem: l.HumanItem, Stream: l.Stream, Live: l.Live}
}
