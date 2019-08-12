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
	PredictionToCategory(e.Prediction) // Use only to print
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
