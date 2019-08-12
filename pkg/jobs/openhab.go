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
	"fmt"
	"net/http"
	"time"

	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
	live "github.com/saljam/mjpeg"
)

type OpenHabErrand struct {
	Prediction api.Prediction
	Options    *OpenHabGeneratorOptions
}

func (e *OpenHabErrand) UpdateItem(item, data string) error {
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	fmt.Println("Updating", item, data)
	req, err := http.NewRequest("POST", e.Options.APIURL+"/rest/items/"+item, bytes.NewBufferString(data))
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

	cat := PredictionToCategory(e.Prediction)

	go func() {
		if cat.Vehicle {
			e.UpdateItem(e.Options.VehicleItem, "ON")
		} else {
			e.UpdateItem(e.Options.VehicleItem, "OFF")
		}
		if cat.Person {
			e.UpdateItem(e.Options.HumanItem, "ON")
		} else {
			e.UpdateItem(e.Options.HumanItem, "OFF")
		}
		if cat.Animal {
			e.UpdateItem(e.Options.AnimalItem, "ON")
		} else {
			e.UpdateItem(e.Options.AnimalItem, "OFF")
		}
	}()
	if e.Options.Live {
		go func() {
			if len(e.Prediction.Url) == 0 {
				return
			}
			b, err := e.Prediction.ToByte()
			if err != nil {
				return
			}
			go e.Options.Stream.UpdateJPEG(b)
		}()
	}
	return nil
}

// Generate makes the errand generates new prediction if needed
func (e *OpenHabErrand) Generate(d api.Detector) *api.Prediction {
	return nil
}

type OpenHabGenerator struct {
	URL     string
	Options *OpenHabGeneratorOptions
}

type OpenHabGeneratorOptions struct {
	APIURL, VehicleItem, HumanItem, AnimalItem string
	Stream                                     *live.Stream
	Live                                       bool
}

func NewOpenHabGenerator(opts *OpenHabGeneratorOptions) errand.ErrandGenerator {
	return &OpenHabGenerator{Options: opts}
}
func (l *OpenHabGenerator) GenerateErrand(p api.Prediction) errand.Errand {
	return &OpenHabErrand{Prediction: p, Options: l.Options}
}
