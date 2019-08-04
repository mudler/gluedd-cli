package generators

import (
	"fmt"
	"strings"

	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
)

type DummyErrand struct {
	Prediction api.Prediction
}

func (e *DummyErrand) Apply() error {
	if e.Prediction.Error != nil || len(e.Prediction.Body.Predictions) == 0 {
		return e.Prediction.Error
	}
	e.Prediction.Explain()
	for _, c := range e.Prediction.Body.Predictions[0].Classes {
		if strings.Contains(c.Cat, "Van") || strings.Contains(c.Cat, "Truck") || strings.Contains(c.Cat, "Car") {
			fmt.Println("Vehicle detected")
		}
	}
	return nil
}

type DummyGenerator struct{}

func NewDummyGenerator() errand.ErrandGenerator {
	return &DummyGenerator{}
}
func (l *DummyGenerator) GenerateErrand(p api.Prediction) errand.Errand {
	return &DummyErrand{Prediction: p}
}
