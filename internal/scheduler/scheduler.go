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

// these channels (one for each running scheduler) tell the scheduler to exit.
//   when you pass a struct{}{} to this channel, the scheduler(s) exits
var Stops = make([]chan struct{}, 0)
// these channels are used to signal that the scheduler exited,
//   this is to ensure that any cleanup required by the scheduler is done.
//   the scheduler sends a struct{}{} over this channel when exiting
var Dones = make([]chan struct{}, 0)
// more about this empty struct: https://dave.cheney.net/2014/03/25/the-empty-struct


func RandomScheduler(tasks chan config.Task, stop chan struct{}, done chan struct{}, nodes []config.Node, wg *sync.WaitGroup) {
	rand.Seed(time.Now().Unix())

// this is a label for the loop below
//  we use this because using a "break"
//  inside the switch statement will only
//  break from the switch and not the loop
//  more info: https://forum.golangbridge.org/t/is-using-continue-label-or-break-label-good-practice/8345
LOOP:
	for {
		select {
		case task := <-tasks:  // got a new task
			go func() {
				node := nodes[rand.Intn(len(nodes))]
				log.Printf("Running task %s on node %s", task.Name, node.Name)

				err := connection.RunTask(node, task.Name, task.Commands)

				// if command errored out or is set to Restart, try running it again
				if err != nil || task.Restart {
					tasks <- task
					wg.Add(1)
				}

				log.Printf("Task %s completed on node %s", task.Name, node.Name)

				wg.Done()
			}()

			case <-stop:  // got the signal to stop :(

			// breaks out of the loop
			// this means that the goroutines run above are not waited for
			//   so, when the daemon exits the SSH connections are broken
			//   with no chance for graceful exit

			break LOOP
		}
	}

	// signal that this scheduler is exiting
	done <- struct{}{}
}
