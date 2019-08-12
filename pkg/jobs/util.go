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
