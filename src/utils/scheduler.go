package utils

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
)

func Schedule(tasks chan Task, nodes []Node, wg *sync.WaitGroup) {
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
