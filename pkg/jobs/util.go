package generators

import (
	"fmt"
	"strings"

	"github.com/mudler/gluedd/pkg/api"
)

type Category struct {
	Person  bool
	Animal  bool
	Vehicle bool
}

func PredictionToCategory(pre api.Prediction) Category {

	cat := Category{Person: false, Animal: false, Vehicle: false}

	for _, c := range pre.Body.Predictions[0].Classes {

		localCat := DecodeCat(c.Cat)
		if localCat.Animal {
			cat.Animal = true
		} else if localCat.Person {
			cat.Person = true
		} else if localCat.Vehicle {
			cat.Vehicle = true
		}

	}

	return cat
}

func DecodeCat(encoded string) Category {
	cat := Category{}
	catPred := strings.ToLower(encoded)
	if strings.Contains(catPred, "van") || strings.Contains(catPred, "truck") || strings.Contains(catPred, "car") {
		cat.Vehicle = true
		fmt.Println("Vehicle detected")
	}
	if strings.Contains(catPred, "man") || strings.Contains(catPred, "person") || strings.Contains(catPred, "face") || strings.Contains(catPred, "woman") {
		cat.Person = true
		fmt.Println("Person detected")
	}

	if strings.Contains(catPred, "animal") || strings.Contains(catPred, "cat") || strings.Contains(catPred, "dog") {
		cat.Animal = true
		fmt.Println("Animal detected")
	}
	return cat
}
