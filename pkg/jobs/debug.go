package generators

import (
	"fmt"
	"strings"

	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
	live "github.com/saljam/mjpeg"
)

type DebugErrand struct {
	Live       bool
	Prediction api.Prediction
	Stream     *live.Stream
}

func (e *DebugErrand) Apply() error {
	if e.Prediction.Error != nil || len(e.Prediction.Body.Predictions) == 0 {
		return e.Prediction.Error
	}
	e.Prediction.Explain()
	for _, c := range e.Prediction.Body.Predictions[0].Classes {
		if strings.Contains(c.Cat, "Van") || strings.Contains(c.Cat, "Truck") || strings.Contains(c.Cat, "Car") {
			fmt.Println("Vehicle detected")
		}
		if strings.Contains(c.Cat, "Man") {
			fmt.Println("Man detected")
		}

		if strings.Contains(c.Cat, "Woman") {
			fmt.Println("Woman detected")
		}

		if strings.Contains(c.Cat, "Person") || strings.Contains(c.Cat, "Face") {
			fmt.Println("Person detected")
		}
	}
	if e.Live {
		go func() {
			if len(e.Prediction.Url) == 0 {
				return
			}
			b, err := e.Prediction.ToByte()
			if err != nil {
				return
			}
			go e.Stream.UpdateJPEG(b)
		}()
	}
	return nil
}

// Generate makes the errand generates new prediction if needed
func (e *DebugErrand) Generate(d api.Detector) *api.Prediction {
	// for _, c := range e.Prediction.Body.Predictions[0].Classes {
	// 	if strings.Contains(c.Cat, "Van") || strings.Contains(c.Cat, "Truck") || strings.Contains(c.Cat, "Car")  || strings.Contains(c.Cat, "Tree") {

	// 		fmt.Println("Extra generation!")
	// 		pre := d.DetectService(e.Prediction.Url,"ilsvrc_googlenet")
	// 		return &pre
	// 	}
	// }

	return nil
}

type DebugGenerator struct {
	Live   bool
	Stream *live.Stream
}

func NewDebugGenerator(stream *live.Stream, live bool) errand.ErrandGenerator {
	return &DebugGenerator{Live: live, Stream: stream}
}
func (l *DebugGenerator) GenerateErrand(p api.Prediction) errand.Errand {
	return &DebugErrand{Prediction: p, Stream: l.Stream}
}
