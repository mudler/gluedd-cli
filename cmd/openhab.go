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

package cmd

import (
	generators "github.com/mudler/gluedd-cli/pkg/jobs"
	types "github.com/mudler/gluedd-cli/pkg/resource"
	live "github.com/saljam/mjpeg"

	"time"

	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
	"github.com/mudler/gluedd/pkg/predictor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var openhabCmd = &cobra.Command{
	Use:   "openhab",
	Short: "Starts gluedd in listen mode",
	Long:  `Redirect jpeg stream detection to openhab item states`,
	Run: func(cmd *cobra.Command, args []string) {
		Service := viper.GetString("service")

		Server := viper.GetString("api_server")

		dd := api.NewDeepDetect(Server, &api.Options{
			Width:      viper.GetInt("image_width"),
			Height:     viper.GetInt("image_height"),
			Detection:  true,
			Confidence: viper.GetFloat64("confidence"),
		})
		if len(Service) > 0 {
			dd.WithService(Service)
		}
		stream := live.NewStream()

		opts := types.JpegStreamerOptions{
			ListeningURL:     viper.GetString("base_url"),
			LiveStreamingURL: viper.GetString("stream_url"),
			Stream:           stream,
			LivePreview:      viper.GetBool("preview"),
			Buffer:           viper.GetInt("buffer_size"),
			Width:            uint(viper.GetInt("image_width")),
			Height:           uint(viper.GetInt("image_height")),
			Resize:           viper.GetBool("resize"),
			Crop:             viper.GetBool("crop"),
			CropMode:         viper.GetString("crop_mode"),
			CropAnchor:       viper.GetBool("crop_anchor"),
			Approx:           viper.GetBool("approx"),
			CropAnchorX:      viper.GetInt("crop_anchor_x"),
			CropAnchorY:      viper.GetInt("crop_anchor_y"),
			Timeout:          viper.GetInt("client_timeout"),

			CropWidth:  viper.GetInt("crop_width"),
			CropHeight: viper.GetInt("crop_height"),
		}

		openhabOptions := &generators.OpenHabGeneratorOptions{
			APIURL:      viper.GetString("openhab_url"),
			VehicleItem: viper.GetString("vehicle_item"),
			AnimalItem:  viper.GetString("animal_item"),
			Stream:      stream,
			Live:        viper.GetBool("preview"),
		}

		errandgen := generators.NewOpenHabGenerator(openhabOptions)
		predictor := predictor.NewPredictor(dd, types.NewJpegStreamer(opts), errandgen)
		consumer := errand.NewErrandConsumer()

		consumer.Consume(predictor.Generate())
		for {
			time.Sleep(1 * time.Second)
		}
	},
}

func init() {
	RootCmd.AddCommand(openhabCmd)
}
