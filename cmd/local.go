package cmd

import (
	"log"

	"github.com/mudler/gluedd-cli/pkg/jobs"
	"github.com/mudler/gluedd-cli/pkg/resource"

	"github.com/mudler/gluedd/pkg/errand"
	"github.com/mudler/gluedd/pkg/predictor"

	"github.com/mudler/gluedd/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Starts gluedd in listen mode",
	Long:  `Starts gluedd with the configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		Service := viper.GetString("service")

		Server := viper.GetString("api_server")
		if len(args) == 0 {
			log.Fatalln("Insufficient arguments")
		}
		dd := api.NewDeepDetect(Server)
		if len(Service) > 0 {
			dd.WithService(Service)
		}

		//errandgen := errand.NewDefaultErrandGenerator()
		errandgen := generators.NewDummyGenerator()
		predictor := predictor.NewPredictor(dd, types.NewDummy(), errandgen)

		//predictor := resource.NewPredictor(dd, resource.NewLocalWatcher(args[0]))
		consumer := errand.NewErrandConsumer()

		consumer.Consume(predictor.Generate())

		for {
		}
	},
}

func init() {
	RootCmd.AddCommand(localCmd)
}
