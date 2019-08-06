package generators

import (
	"bytes"
	"fmt"
	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type OpenHabErrand struct {
	WithRemoval                      bool
	Prediction                       api.Prediction
	URL                              string
	AssetDir, VehicleItem, HumanItem string
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
	r, err := http.NewRequest("GET", e.Prediction.Url, nil)
	if err != nil {
		return err
	}
	if e.WithRemoval {
		file := path.Base(r.URL.Path)
		os.RemoveAll(filepath.Join(e.AssetDir, file))
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

	if vehicle {
		err := e.UpdateItem(e.VehicleItem, "ON")
		if err != nil {
			return err
		}
	} else {
		err := e.UpdateItem(e.VehicleItem, "OFF")
		if err != nil {
			return err
		}
	}
	if person {
		err := e.UpdateItem(e.HumanItem, "ON")
		if err != nil {
			return err
		}
	} else {
		err := e.UpdateItem(e.HumanItem, "OFF")
		if err != nil {
			return err
		}
	}
	return nil
}

// Generate makes the errand generates new prediction if needed
func (e *OpenHabErrand) Generate(d api.Detector) *api.Prediction {
	return nil
}

type OpenHabGenerator struct {
	URL                              string
	AssetDir, VehicleItem, HumanItem string
	WithRemoval                      bool
}

func NewOpenHabGenerator(OpenHabURL, VehicleItem, HumanItem, AssetDir string, Remove bool) errand.ErrandGenerator {
	return &OpenHabGenerator{URL: OpenHabURL, VehicleItem: VehicleItem, HumanItem: HumanItem, AssetDir: AssetDir, WithRemoval: Remove}
}
func (l *OpenHabGenerator) GenerateErrand(p api.Prediction) errand.Errand {
	return &OpenHabErrand{Prediction: p, URL: l.URL, VehicleItem: l.VehicleItem, HumanItem: l.HumanItem, AssetDir: l.AssetDir, WithRemoval: l.WithRemoval}
}
