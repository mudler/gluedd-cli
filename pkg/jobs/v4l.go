package generators

import (
	"errors"

	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
	live "github.com/saljam/mjpeg"
)

// Parts of this errand are extracted from https://github.com/jolibrain/livedetect/blob/master/writeBoundingBox.go

type V4lErrand struct {
	Prediction api.Prediction
	Stream     *live.Stream
}

func (e *V4lErrand) Apply() error {
	if e.Prediction.Error != nil || len(e.Prediction.Body.Predictions) == 0 {
		return e.Prediction.Error
	}

	if len(e.Prediction.Body.Predictions) != 0 {
		b, err := e.Prediction.ToByte()
		if err != nil {
			return err
		}
		go e.Stream.UpdateJPEG(b)
	} else {
		return errors.New("Can't encode frame to live stream.")
	}

	return nil
}

// Generate makes the errand generates new prediction if needed
func (e *V4lErrand) Generate(d api.Detector) *api.Prediction {
	return nil
}

type V4lGenerator struct{ Stream *live.Stream }

func NewV4lGenerator(Stream *live.Stream) errand.ErrandGenerator {
	return &V4lGenerator{Stream: Stream}
}
func (l *V4lGenerator) GenerateErrand(p api.Prediction) errand.Errand {
	return &V4lErrand{Prediction: p, Stream: l.Stream}
}
