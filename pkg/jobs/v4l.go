package generators

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/draw"

	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
	jpeg "github.com/pixiv/go-libjpeg/jpeg"
	live "github.com/saljam/mjpeg"

	"github.com/CorentinB/gobbox"
)

// Parts of this errand are extracted from https://github.com/jolibrain/livedetect/blob/master/writeBoundingBox.go

type V4lErrand struct {
	Prediction api.Prediction
	Stream     *live.Stream
}

func PredictionToByte(e api.Prediction) ([]byte, error) {
	unbased, err := base64.StdEncoding.DecodeString(e.Url)
	if err != nil {
		return []byte{}, err
	}
	img, err := jpeg.Decode(bytes.NewReader(unbased), &jpeg.DecoderOptions{})
	if err != nil {
		return []byte{}, err
	}
	if len(e.Body.Predictions) != 0 {
		// Convert to RGBA
		b := img.Bounds()
		imgRGBA := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(imgRGBA, imgRGBA.Bounds(), img, b.Min, draw.Src)

		for i, _ := range e.Body.Predictions[0].Classes {
			imgRGBA = writeBoundingBox(img, e, i)
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

	}

	return []byte{}, nil
}

func (e *V4lErrand) Apply() error {
	if e.Prediction.Error != nil || len(e.Prediction.Body.Predictions) == 0 {
		return e.Prediction.Error
	}

	if len(e.Prediction.Body.Predictions) != 0 {
		b, err := PredictionToByte(e.Prediction)
		if err != nil {
			return err
		}
		go e.Stream.UpdateJPEG(b)
	} else {
		return errors.New("Can't encode frame to live stream.")
	}

	return nil
}

// Generate makes the errand generates new prediction if needed
func (e *V4lErrand) Generate(d api.Detector) *api.Prediction {
	return nil
}

type V4lGenerator struct{ Stream *live.Stream }

func NewV4lGenerator(Stream *live.Stream) errand.ErrandGenerator {
	return &V4lGenerator{Stream: Stream}
}
func (l *V4lGenerator) GenerateErrand(p api.Prediction) errand.Errand {
	return &V4lErrand{Prediction: p, Stream: l.Stream}
}

func writeBoundingBox(img image.Image, result api.Prediction, class int) (imgRGBA *image.RGBA) {
	// Convert to RGBA
	b := img.Bounds()
	imgRGBA = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(), img, b.Min, draw.Src)

	// Set colors for bbox
	red := color.RGBA{255, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}

	// Set coordinates for bbox
	x1 := int(result.Body.Predictions[0].Classes[class].Bbox.Xmin)
	x2 := int(result.Body.Predictions[0].Classes[class].Bbox.Xmax)
	y1 := int(result.Body.Predictions[0].Classes[class].Bbox.Ymin)
	y2 := int(result.Body.Predictions[0].Classes[class].Bbox.Ymax)

	// Draw the bounding box
	gobbox.DrawBoundingBox(imgRGBA, result.Body.Predictions[0].Classes[class].Cat, x1, x2, y2, y1, red, white)

	return imgRGBA
}
