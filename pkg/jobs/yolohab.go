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

package generators

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jolibrain/godd"
	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
	live "github.com/saljam/mjpeg"
)

type YoloHabConnector struct {
	Prediction api.Prediction
	Stream     *live.Stream
	Server     string
}

func (e *YoloHabConnector) Apply() error {
	if e.Prediction.Error != nil || len(e.Prediction.Body.Predictions) == 0 {
		return e.Prediction.Error
	}

	e.Prediction.Explain()
	PredictionToCategory(e.Prediction) // Use only to print

	if len(e.Prediction.Url) != 0 {
		go func() {
			b, err := e.Prediction.ToByte()
			if err != nil {
				return
			}
			e.Stream.UpdateJPEG(b)
		}()
	} else {
		return errors.New("Can't encode frame to live stream.")
	}

	return nil
}

// Generate makes the errand generates new prediction if needed
func (e *YoloHabConnector) Generate(d api.Detector) *api.Prediction {
	return nil
}

type YoloGenerator struct{ Stream *live.Stream }

func NewYoloGenerator(Stream *live.Stream) errand.ErrandGenerator {
	return &YoloGenerator{Stream: Stream}
}
func (l *YoloGenerator) GenerateErrand(p api.Prediction) errand.Errand {
	return &YoloHabConnector{Prediction: p, Stream: l.Stream}
}

func (d *YoloHabConnector) CreateService(URL string, arguments api.ServiceOptions) error {
	return nil
}

func (d *YoloHabConnector) DetectService(photo, service string) api.Prediction {
	return api.Prediction{}
}

type Request struct {
	Image string `json:"image"`
}

type BoundingBox struct {
	StartPoint image.Point
	EndPoint   image.Point
}
type Detection struct {
	BoundingBox

	ClassIDs      []int
	ClassNames    []string
	Probabilities []float32
}

// DetectionResult represents the inference results from the network.
type DetectionResult struct {
	Detections           []*Detection
	NetworkOnlyTimeTaken time.Duration
	OverallTimeTaken     time.Duration
}

type Bbox struct {
	Ymax float64 "json:\"ymax,omitempty\""
	Xmax float64 "json:\"xmax,omitempty\""
	Ymin float64 "json:\"ymin,omitempty\""
	Xmin float64 "json:\"xmin,omitempty\""
}
type Class struct {
	Prob float64 "json:\"prob,omitempty\""
	Last bool    "json:\"last,omitempty\""
	Bbox Bbox    "json:\"bbox,omitempty\""
	Cat  string  "json:\"cat,omitempty\""
}

type Prediction struct {
	Classes []Class
}
type Body struct {
	Predictions []Prediction "json:\"predictions,omitempty\""
}
type PredictResult struct {
	Body Body "json:\"body,omitempty\""
}

func (d *YoloHabConnector) Detect(photo string) api.Prediction {

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&Request{Image: photo})
	req, _ := http.NewRequest("POST", d.Server, buf)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, e := client.Do(req)
	if e != nil {
		log.Fatal(e)
	}

	defer res.Body.Close()

	fmt.Println("response Status:", res.Status)

	result := bytes.NewBuffer(nil)
	// Print the body to the stdout
	io.Copy(result, res.Body)

	var detectResult DetectionResult
	json.Unmarshal(result.Bytes(), &detectResult)

	returnedResult := godd.PredictResult{}
	fakeResult := PredictResult{}
	fakeResult.Body.Predictions = make([]Prediction, 0)
	fakeResult.Body.Predictions = append(fakeResult.Body.Predictions, Prediction{
		Classes: make([]Class, 0),
	})

	fmt.Println(detectResult.Detections)
	for _, dets := range detectResult.Detections {
		//returnedResult.Body.Predictions
		fakeResult.Body.Predictions[0].Classes = append(fakeResult.Body.Predictions[0].Classes,
			Class{
				//	Prob: float64(dets.Probabilities[0]),
				Cat: strings.Join(dets.ClassNames, ", "),
				Bbox: Bbox{
					Xmax: float64(dets.BoundingBox.StartPoint.X),
					Ymax: float64(dets.BoundingBox.StartPoint.Y),
					Xmin: float64(dets.BoundingBox.EndPoint.X),
					Ymin: float64(dets.BoundingBox.EndPoint.Y),
				},
			})

	}

	data, _ := json.Marshal(fakeResult)
	_ = json.Unmarshal(data, &returnedResult)

	fmt.Println("Returning", returnedResult)
	return api.Prediction{
		PredictResult: returnedResult,
		Url:           photo,
	}
}

func (d *YoloHabConnector) WithService(s string) api.Detector {
	return d
}
