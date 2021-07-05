package cmd

import (
	"log"

	goDaemon "github.com/sevlyar/go-daemon"
	"github.com/psiayn/heiko/internal/daemon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops a running heiko daemon",
	Run: func(cmd *cobra.Command, args []string) {
		context := daemon.GetContext()

		// searching if daemon exists
		process, err := context.Search()
		if err != nil {
			log.Fatalf("Could not find daemon %s: %v", viper.GetString("name"), err)
		}
		err = goDaemon.SendCommands(process)
		if err != nil {
			log.Fatalf("Could not stop daemon %s: %v", viper.GetString("name"), err)
		}
	},
}
