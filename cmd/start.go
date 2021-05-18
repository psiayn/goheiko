package cmd

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/psiayn/heiko/internal/config"
	"github.com/psiayn/heiko/internal/scheduler"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new heiko job",
	Run: func(cmd *cobra.Command, args []string) {
		task_arr := configuration.Jobs
		nodes := configuration.Nodes
		workDir := filepath.Join(
			viper.GetString("dataLocation"),
			viper.GetString("name"),
		)
		if viper.GetBool("daemon") {
			context := &daemon.Context{
				PidFileName: filepath.Join(
					workDir,
					"daemon.pid",
				),
				PidFilePerm: 0644,
				LogFileName: filepath.Join(
					workDir,
					"daemon.log",
				),
				LogFilePerm: 0644,
				WorkDir:     ".",
				Umask:       022,
				// Args:        nil,
				// Env:         nil,
			}

			d, err := context.Reborn()
			if err != nil {
				log.Fatalln("Error while daemonizing!", err)
				return
			}
			if d != nil {
				return
			}
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
