package cmd

import (
	"log"

	"github.com/mudler/gluedd/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
		dd.Detect(args[0]).Explain()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
