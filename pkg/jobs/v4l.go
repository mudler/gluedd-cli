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
