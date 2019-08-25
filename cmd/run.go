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
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pixiv/go-libjpeg/jpeg"

	"github.com/mudler/gluedd/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DownloadEncode(url string) (string, error) {
	timeout := time.Duration(2 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	// Get the data
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	img, err := jpeg.DecodeIntoRGBA(resp.Body, &jpeg.DecoderOptions{})
	if err != nil {
		return "", err
	}

	// Encode as base64
	buffer64 := new(bytes.Buffer)
	err = jpeg.Encode(buffer64, img, &jpeg.EncoderOptions{Quality: 100})
	if err != nil {
		return "", err
	}
	imageBase64 := base64.StdEncoding.EncodeToString(buffer64.Bytes())
	return imageBase64, nil
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts gluedd",
	Long:  `Starts gluedd with the configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		Service := viper.GetString("service")

		Server := viper.GetString("api_server")
		if len(args) == 0 {
			log.Fatalln("Insufficient arguments")
		}
		dd := api.NewDeepDetect(Server, &api.Options{
			Width:      viper.GetInt("image_width"),
			Height:     viper.GetInt("image_height"),
			Detection:  true,
			Confidence: viper.GetFloat64("confidence"),
		})
		if len(Service) > 0 {
			dd.WithService(Service)
		}
		image, err := DownloadEncode(args[0])
		if err != nil {
			log.Fatalln(err.Error())
		}
		pred := dd.Detect(image)
		pred.Explain()
		b, err := pred.ToByte()
		if err != nil {
			log.Fatalln(err.Error())
		}

		out := "out.jpeg"
		if len(args) == 2 && len(args[1]) > 0 {
			out = args[1]
		}
		err = ioutil.WriteFile(out, b, 0644)
		if err != nil {
			log.Fatalln(err.Error())
		}

	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
