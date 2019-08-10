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
	Prediction                         api.Prediction
	URL                                string
	VehicleItem, HumanItem, AnimalItem string
	Stream                             *live.Stream
	Live                               bool
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

	var cat Category
	for _, c := range e.Prediction.Body.Predictions[0].Classes {

		localCat := DecodeCat(c.Cat)
		if localCat.Animal {
			cat.Animal = true
		} else if localCat.Person {
			cat.Person = true
		} else if localCat.Vehicle {
			cat.Vehicle = true
		}

	}

	go func() {
		if cat.Vehicle {
			e.UpdateItem(e.VehicleItem, "ON")
		} else {
			e.UpdateItem(e.VehicleItem, "OFF")
		}
		if cat.Person {
			e.UpdateItem(e.HumanItem, "ON")
		} else {
			e.UpdateItem(e.HumanItem, "OFF")
		}
		if cat.Animal {
			e.UpdateItem(e.AnimalItem, "ON")
		} else {
			e.UpdateItem(e.AnimalItem, "OFF")
		}
	}()
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
func (e *OpenHabErrand) Generate(d api.Detector) *api.Prediction {
	return nil
}

type OpenHabGenerator struct {
	URL                                string
	VehicleItem, HumanItem, AnimalItem string
	Stream                             *live.Stream
	Live                               bool
}

func NewOpenHabGenerator(OpenHabURL, VehicleItem, HumanItem, AnimalItem string, stream *live.Stream, live bool) errand.ErrandGenerator {
	return &OpenHabGenerator{URL: OpenHabURL, VehicleItem: VehicleItem, HumanItem: HumanItem, AnimalItem: AnimalItem, Stream: stream, Live: live}
}
func (l *OpenHabGenerator) GenerateErrand(p api.Prediction) errand.Errand {
	return &OpenHabErrand{Prediction: p, URL: l.URL, VehicleItem: l.VehicleItem, HumanItem: l.HumanItem, AnimalItem: l.AnimalItem, Stream: l.Stream, Live: l.Live}
}
