package cmd

import (

	"github.com/mudler/gluedd-cli/pkg/jobs"
	"github.com/mudler/gluedd-cli/pkg/resource"

	"github.com/mudler/gluedd/pkg/errand"
	"github.com/mudler/gluedd/pkg/predictor"
"time"
	"github.com/mudler/gluedd/pkg/api"
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
	
		dd := api.NewDeepDetect(Server)
		if len(Service) > 0 {
			dd.WithService(Service)
		}
		
		//errandgen := errand.NewDefaultErrandGenerator()
		errandgen := generators.NewOpenHabGenerator(viper.GetString("openhab_url"),viper.GetString("vehicle_item"),viper.GetString("person_item"))
		predictor := predictor.NewPredictor(dd, types.NewJpegStreamer(viper.GetString("stream_url"), viper.GetString("base_url"), viper.GetString("asset_dir")), errandgen)

		//predictor := resource.NewPredictor(dd, resource.NewopenhabWatcher(args[0]))
		consumer := errand.NewErrandConsumer()

		consumer.Consume(predictor.Generate())	
		for {
			time.Sleep(1*time.Second)
		}
	},
}

func init() {
	RootCmd.AddCommand(openhabCmd)
}
