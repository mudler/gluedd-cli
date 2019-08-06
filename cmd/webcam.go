package cmd

import (
	"github.com/mudler/gluedd-cli/pkg/jobs"
	"github.com/mudler/gluedd-cli/pkg/resource"
	live "github.com/saljam/mjpeg"

	"github.com/mudler/gluedd/pkg/api"
	"github.com/mudler/gluedd/pkg/errand"
	"github.com/mudler/gluedd/pkg/predictor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"time"
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
		//errandgen := errand.NewDefaultErrandGenerator()
		errandgen := generators.NewV4lGenerator(stream)
		predictor := predictor.NewPredictor(dd, types.NewV4lStreamer(deviceID, viper.GetString("base_url"), 800, 600, stream), errandgen)

		//predictor := resource.NewPredictor(dd, resource.NewwebcamWatcher(args[0]))
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
