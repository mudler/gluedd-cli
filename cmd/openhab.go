package cmd

import (
	"github.com/mudler/gluedd-cli/pkg/jobs"
	"github.com/mudler/gluedd-cli/pkg/resource"
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

		dd := api.NewDeepDetect(Server, nil)
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
			Approx:           viper.GetBool("approx"),
		}

		//errandgen := errand.NewDefaultErrandGenerator()
		errandgen := generators.NewOpenHabGenerator(viper.GetString("openhab_url"), viper.GetString("vehicle_item"), viper.GetString("person_item"), viper.GetString("animal_item"), stream, true)
		predictor := predictor.NewPredictor(dd, types.NewJpegStreamer(opts), errandgen)

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
