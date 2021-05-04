// package scheduler implements different schedulers
// for selecting appropriate nodes to run tasks
package scheduler

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
	"github.com/psiayn/heiko/internal/config"
)

func Schedule(tasks chan config.Task, nodes []config.Node, wg *sync.WaitGroup) {
	rand.Seed(time.Now().Unix())
	for {
		task := <- tasks
		go func() {
			fmt.Println(task, " on node ", nodes[rand.Intn(len(nodes))].Name)
			time.Sleep(10)
			wg.Done()
		}()
	}
}
