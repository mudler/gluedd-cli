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
	"encoding/json"
	"errors"

	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
	live "github.com/saljam/mjpeg"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/mqtt"
)

// Parts of this errand are extracted from https://github.com/jolibrain/livedetect/blob/master/writeBoundingBox.go

type MQTTErrand struct {
	Prediction api.Prediction
	Stream     *live.Stream
	MQTTBroker *mqtt.Adaptor
}

func (e *MQTTErrand) Apply() error {
	if e.Prediction.Error != nil || len(e.Prediction.Body.Predictions) == 0 {
		return e.Prediction.Error
	}

	work := func() {
		data, _ := json.Marshal(e.Prediction)
		e.MQTTBroker.Publish("detection", data)
	}

	robot := gobot.NewRobot("mqttBot",
		[]gobot.Connection{e.MQTTBroker},
		work,
	)

	robot.Start(false)
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
func (e *MQTTErrand) Generate(d api.Detector) *api.Prediction {
	return nil
}

type MQTTGenerator struct {
	Stream     *live.Stream
	MQTTBroker *mqtt.Adaptor
}

func NewMQTTGenerator(Stream *live.Stream, MQTTBroker, MQTTChannel string) errand.ErrandGenerator {
	mqttAdaptor := mqtt.NewAdaptor(MQTTBroker, MQTTChannel)
	return &MQTTGenerator{Stream: Stream, MQTTBroker: mqttAdaptor}
}
func (l *MQTTGenerator) GenerateErrand(p api.Prediction) errand.Errand {
	return &MQTTErrand{Prediction: p, Stream: l.Stream, MQTTBroker: l.MQTTBroker}
}
