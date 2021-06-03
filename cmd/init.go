package cmd

import (
	"log"
	"sync"

	"github.com/psiayn/heiko/internal/config"
	"github.com/psiayn/heiko/internal/connection"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Runs initialization of Jobs",
	Run: func(cmd *cobra.Command, args []string) {

		err := config.SetAuth(&configuration)
		if err != nil {
			log.Fatalln(err)
		}

		jobs := configuration.Jobs
		nodes := configuration.Nodes

		var wg sync.WaitGroup

		for _, job := range jobs {
			// we run initialization of this job for each node
			//    for each node, we run it in a separate goroutine
			wg.Add(len(nodes))

			for _, node := range nodes {

				go func(job config.Task, node config.Node) {
					defer wg.Done()
					log.Printf("Running initialization for task %s on node %s",
						job.Name,
						node.Name)
					connection.RunTask(node, job.Name, job.Init)
				}(job, node)

			}
			wg.Wait()
		}
	},
}
