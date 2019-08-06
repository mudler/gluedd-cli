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
	"time"
)

var openhabCmd = &cobra.Command{
	Use:   "openhab",
	Short: "Starts gluedd in listen mode",
	Long:  `Redirect jpeg stream detection to openhab item states`,
	Run: func(cmd *cobra.Command, args []string) {
		Service := viper.GetString("service")

		Server := viper.GetString("api_server")

		dd := api.NewDeepDetect(Server, nil)
		if len(Service) > 0 {
			dd.WithService(Service)
		}
		stream := live.NewStream()

		//errandgen := errand.NewDefaultErrandGenerator()
		errandgen := generators.NewOpenHabGenerator(viper.GetString("openhab_url"), viper.GetString("vehicle_item"), viper.GetString("person_item"), stream, true)
		predictor := predictor.NewPredictor(dd, types.NewJpegStreamer(viper.GetString("stream_url"), viper.GetString("base_url"), stream, true), errandgen)

		//predictor := resource.NewPredictor(dd, resource.NewopenhabWatcher(args[0]))
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
