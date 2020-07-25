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
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jolibrain/godd"
	"github.com/mudler/gluedd/pkg/api"
	live "github.com/saljam/mjpeg"
)

type YoloCV struct {
	Prediction api.Prediction
	Stream     *live.Stream
	Server     string
}

// Generate makes the errand generates new prediction if needed
func (e *YoloCV) Generate(d api.Detector) *api.Prediction {
	return nil
}

func (d *YoloCV) CreateService(URL string, arguments api.ServiceOptions) error {
	return nil
}

func (d *YoloCV) DetectService(photo, service string) api.Prediction {
	return api.Prediction{}
}

func (d *YoloCV) Detect(photo string) api.Prediction {

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

	returnedResult := godd.PredictResult{}

	_ = json.Unmarshal(result.Bytes(), &returnedResult)
	fmt.Println("Returning", returnedResult)
	return api.Prediction{
		PredictResult: returnedResult,
		Url:           photo,
	}
}

func (d *YoloCV) WithService(s string) api.Detector {
	return d
}
