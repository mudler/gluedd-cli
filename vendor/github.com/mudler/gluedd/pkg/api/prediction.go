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
	"bytes"
	"fmt"

	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"strconv"

	"github.com/CorentinB/gobbox"
	"github.com/jolibrain/godd"
	jpeg "github.com/pixiv/go-libjpeg/jpeg"
)

type Prediction struct {
	godd.PredictResult
	Error error
	Url   string
}

func (p Prediction) Explain() {
	fmt.Println("")

	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	fmt.Println("Predictions:", len(p.Body.Predictions))
	for _, i := range p.Body.Predictions {
		fmt.Println("URI", i.URI)
		for _, c := range i.Classes {
			fmt.Println("Category:", c.Cat+" [prob "+strconv.FormatFloat(c.Prob, 'f', 6, 64)+"]")
		}
	}
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
}

// Parts extracted from: https://github.com/jolibrain/livedetect

func (p Prediction) ToByte() ([]byte, error) {
	unbased, err := base64.StdEncoding.DecodeString(p.Url)
	if err != nil {
		return []byte{}, err
	}
	img, err := jpeg.Decode(bytes.NewReader(unbased), &jpeg.DecoderOptions{})
	if err != nil {
		return []byte{}, err
	}
	if len(
		p.Body.Predictions) == 0 {
		return []byte{}, err
	}
	// Convert to RGBA
	b := img.Bounds()
	imgRGBA := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(), img, b.Min, draw.Src)

	for i, _ := range p.Body.Predictions[0].Classes {
		imgRGBA = p.writeBoundingBox(img, i)
		img = imgRGBA
	}

	buf := new(bytes.Buffer)
	if imgRGBA != nil {
		err := jpeg.Encode(buf, imgRGBA, &jpeg.EncoderOptions{Quality: 50})
		if err == nil {
			return buf.Bytes(), nil
		} else {
			return []byte{}, nil
		}
	}

	return []byte{}, nil
}

func (p Prediction) writeBoundingBox(img image.Image, class int) (imgRGBA *image.RGBA) {
	// Convert to RGBA
	b := img.Bounds()
	imgRGBA = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(), img, b.Min, draw.Src)

	// Set colors for bbox
	red := color.RGBA{255, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}

	// Set coordinates for bbox
	x1 := int(p.Body.Predictions[0].Classes[class].Bbox.Xmin)
	x2 := int(p.Body.Predictions[0].Classes[class].Bbox.Xmax)
	y1 := int(p.Body.Predictions[0].Classes[class].Bbox.Ymin)
	y2 := int(p.Body.Predictions[0].Classes[class].Bbox.Ymax)

	// Draw the bounding box
	gobbox.DrawBoundingBox(imgRGBA, p.Body.Predictions[0].Classes[class].Cat, x1, x2, y2, y1, red, white)

	return imgRGBA
}
