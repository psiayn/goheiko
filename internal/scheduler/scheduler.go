// package scheduler implements different schedulers
// for selecting appropriate nodes to run tasks
package scheduler

import (
	"github.com/psiayn/heiko/internal/config"
	"github.com/psiayn/heiko/internal/connection"
	"log"
	"math/rand"
	"sync"
	"time"
)

func RandomScheduler(tasks chan config.Task, nodes []config.Node, wg *sync.WaitGroup) {
	rand.Seed(time.Now().Unix())
	for {
		task := <-tasks
		go func() {
			node := nodes[rand.Intn(len(nodes))]
			log.Printf("Running task %s on node %s", task.Name, node.Name)

			err := connection.RunTask(node, task.Name, task.Commands)

			// if command errored out or is set to Restart, try running it again
			if err != nil || task.Restart {
				tasks <- task
				wg.Add(1)
			}

			wg.Done()
		}()
	}
}
