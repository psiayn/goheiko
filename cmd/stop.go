package cmd

import (
	"log"
	"syscall"

	"github.com/psiayn/heiko/internal/daemon"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops a running heiko daemon",
	Run: func(cmd *cobra.Command, args []string) {
		context := daemon.GetContext()
		// searching if daemon exists
		process, err := context.Search()
		if err != nil {
			log.Fatalln(err)
		}
		process.Signal(syscall.SIGTERM)
		defer context.Release()

	},
}
