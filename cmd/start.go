package cmd

import (
	"log"
	"sync"

	"github.com/psiayn/heiko/internal/config"
	"github.com/psiayn/heiko/internal/scheduler"
	"github.com/psiayn/heiko/internal/daemon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new heiko job",
	Run: func(cmd *cobra.Command, args []string) {
		task_arr := configuration.Jobs
		nodes := configuration.Nodes

		// handle daemonizing now if required
		if viper.GetBool("daemon") {
			context := daemon.GetContext()

			// this essentially forks the program from this point
			// child is running as a daemon!
			// d is nil in child and non-nil in parent
			d, err := context.Reborn()
			if err != nil {
				log.Fatalln("Error while daemonizing!", err)
				return
			}

			if d != nil {
				// exit parent
				return
			}

			// important - this releases the pidfile
			//             once the program completes
			defer context.Release()
			log.Print("- - - - - - - - - - - - - - -")
			log.Print("daemon started")
		}

		log.Println("len of nodes = ", len(task_arr))
		tasks := make(chan config.Task)

		var wg sync.WaitGroup
		wg.Add(len(task_arr))
		go scheduler.RandomScheduler(tasks, nodes, &wg)
		for _, task := range task_arr {
			tasks <- task
		}
		wg.Wait()

	},
}
