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
	"github.com/mudler/gluedd-cli/pkg/jobs"
	"github.com/mudler/gluedd-cli/pkg/resource"
	live "github.com/saljam/mjpeg"

	"log"
	"strconv"
	"time"

	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
	"github.com/mudler/gluedd/pkg/predictor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var webcamCmd = &cobra.Command{
	Use:   "webcam",
	Short: "Starts gluedd in listen mode",
	Long:  `Starts gluedd with the configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		Service := viper.GetString("service")

		Server := viper.GetString("api_server")

		dd := api.NewDeepDetect(Server, nil)
		if len(Service) > 0 {
			dd.WithService(Service)
		}
		if len(args) == 0 {
			log.Fatalln("Insufficient arguments")
		}

		deviceID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Invalid devide ID")
		}
		stream := live.NewStream()

		opts := &types.V4lStreamerOptions{
			Device:    deviceID,
			Stream:    stream,
			Buffer:    viper.GetInt("buffer_size"),
			StreamURL: viper.GetString("base_url"),
			Width:     viper.GetInt("webcam_width"),
			Height:    viper.GetInt("webcam_width"),
		}

		errandgen := generators.NewV4lGenerator(stream)
		predictor := predictor.NewPredictor(dd, types.NewV4lStreamer(opts), errandgen)
		consumer := errand.NewErrandConsumer()

		consumer.Consume(predictor.Generate())
		for {
			time.Sleep(1 * time.Second)
		}
	},
}

func init() {
	RootCmd.AddCommand(webcamCmd)
}
