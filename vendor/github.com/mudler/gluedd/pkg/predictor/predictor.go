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

package predictor

import (
	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
	"github.com/mudler/gluedd/pkg/resource"
)

type DefaultPredictor struct {
	API       api.Detector
	Resource  resource.Resource
	Generator errand.ErrandGenerator
}

func NewPredictor(API api.Detector, res resource.Resource, err errand.ErrandGenerator) Predictor {
	return &DefaultPredictor{API: API, Resource: res, Generator: err}
}

func (p *DefaultPredictor) Generate() chan errand.Errand {
	jobs := make(chan errand.Errand)
	predictions := p.Predict()
	go func() {
		for pre := range predictions {

			errand := p.Generator.GenerateErrand(pre)
			extraPrediction := errand.Generate(p.API)
			if extraPrediction != nil {
				jobs <- p.Generator.GenerateErrand(*extraPrediction)
			}
			jobs <- errand
		}
	}()

	return jobs
}

func (p *DefaultPredictor) Predict() chan api.Prediction {

	predictions := make(chan api.Prediction)
	uris := p.Resource.Listen()
	go func() {
		for s := range uris {
			predictions <- p.API.Detect(s)
		}
	}()

	return predictions
}
