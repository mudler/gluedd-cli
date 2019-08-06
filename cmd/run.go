package cmd

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mudler/gluedd-cli/pkg/jobs"
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
		dd := api.NewDeepDetect(Server, nil)
		if len(Service) > 0 {
			dd.WithService(Service)
		}
		image, err := DownloadEncode(args[0])
		if err != nil {
			log.Fatalln(err.Error())
		}
		pred := dd.Detect(image)
		pred.Explain()
		b, err := generators.PredictionToByte(pred)
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
