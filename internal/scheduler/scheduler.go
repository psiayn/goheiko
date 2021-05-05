// package scheduler implements different schedulers
// for selecting appropriate nodes to run tasks
package scheduler

import (
	"fmt"
	"github.com/psiayn/heiko/internal/config"
	"github.com/psiayn/heiko/internal/connection"
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
			fmt.Println(task, " on node ", node.Name)
			connection.Connect(node, task)
			wg.Done()
		}()
	}
}
