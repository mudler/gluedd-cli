package generators

import (
	"fmt"
	"strings"
)

type Category struct {
	Person  bool
	Animal  bool
	Vehicle bool
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
