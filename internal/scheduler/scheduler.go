// package scheduler implements different schedulers
// for selecting appropriate nodes to run tasks
package scheduler

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
	"github.com/psiayn/heiko/internal/config"
	"github.com/psiayn/heiko/internal/connection"
)

func Schedule(tasks chan config.Task, nodes []config.Node, wg *sync.WaitGroup) {
	rand.Seed(time.Now().Unix())
	for {
		task := <- tasks
		go func() {
			node := nodes[rand.Intn(len(nodes))]
			fmt.Println(task, " on node ", node.Name)
			time.Sleep(10)
			connection.Connect(node, task)
			wg.Done()
		}()
	}
}
