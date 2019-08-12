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

package api

import (
	"errors"

	"github.com/jolibrain/godd"
)

type Options struct {
	Width, Height, Best  int
	Detection, Mask, GPU bool
	Confidence           float64
}

type ServiceOptions struct {
	Create                bool
	GPU                   bool
	Nclasses              int
	Template              string
	ModelRepository       string
	ServiceDescription    string
	Mllib                 string
	Connector             string
	Init                  string
	MlLibDataType         string
	MlLibMaxBatchSize     int
	MlLibMaxWorkspaceSize int
	Width                 int
	Height                int
	Service               string
	Mask                  bool
	Extensions            *[]string
	Mean                  *[]float64
}
type DeepDetect struct {
	Server  string
	service string
	Options *Options
}

func NewDeepDetect(server string, opts *Options) Detector {
	return &DeepDetect{Server: server, Options: opts}
}

func (d *DeepDetect) CreateService(URL string, arguments ServiceOptions) error {

	// Create service struct for service creation
	// parameters
	var service godd.ServiceRequest

	// Fill the service structure
	service.Name = arguments.Service
	service.Mllib = arguments.Mllib
	service.Parameters.Input.Connector = arguments.Connector
	service.Parameters.Input.Width = arguments.Width
	service.Parameters.Input.Height = arguments.Height
	service.Model.Repository = arguments.ModelRepository
	service.Model.Init = arguments.Init
	service.Parameters.Mllib.Nclasses = arguments.Nclasses
	service.Parameters.Mllib.GPU = arguments.GPU

	if len(arguments.MlLibDataType) > 0 {
		service.Parameters.Mllib.Datatype = arguments.MlLibDataType
	}

	if arguments.MlLibMaxBatchSize != -1 {
		service.Parameters.Mllib.MaxBatchSize = arguments.MlLibMaxBatchSize
	}

	if arguments.MlLibMaxWorkspaceSize != -1 {
		service.Parameters.Mllib.MaxWorkspaceSize = arguments.MlLibMaxWorkspaceSize
	}

	if service.Model.Init != "" {
		service.Model.CreateRepository = true
	}

	// Mask support
	if arguments.Mask == true {
		if len(*arguments.Extensions) == 0 {
			return errors.New("You need to specify at least one extension for mask")
		}
		arguments.Connector = "image"
		arguments.Mllib = "caffe2"
		service.Parameters.Input.Mean = *arguments.Mean
		service.Model.Extensions = *arguments.Extensions
	}

	// Send the service creation request
	creationResult, err := godd.CreateService(URL, &service)
	if err != nil {
		return err
	}

	if creationResult.Status.Code != 201 {
		return errors.New(creationResult.Status.Msg)
	}

	return nil
}

func (d *DeepDetect) DetectService(photo, service string) Prediction {

	// Create predict structure for request parameters
	var predict godd.PredictRequest

	predict.Service = service
	predict.Data = append(predict.Data, photo)
	predict.Parameters.Output.Bbox = true
	predict.Parameters.Output.ConfidenceThreshold = 0.1

	if d.Options != nil {
		predict.Parameters.Input.Width = d.Options.Width
		predict.Parameters.Input.Height = d.Options.Height
		predict.Parameters.Output.Best = d.Options.Best
		predict.Parameters.Output.Bbox = d.Options.Detection
		predict.Parameters.Output.ConfidenceThreshold = d.Options.Confidence
		predict.Parameters.Mllib.GPU = d.Options.GPU
		predict.Parameters.Output.Mask = d.Options.Mask
	}

	predictResult, err := godd.Predict(d.Server, &predict)
	if err != nil {
		return Prediction{Error: err}
	}

	return Prediction{PredictResult: predictResult, Url: photo}
}

func (d *DeepDetect) Detect(photo string) Prediction {

	service := "detection_600"
	if len(d.service) > 0 {
		service = d.service
	}
	return d.DetectService(photo, service)
}

func (d *DeepDetect) WithService(s string) Detector {
	d.service = s
	return d
}
