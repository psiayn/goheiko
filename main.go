package main

import (
	"os"
	"sync"
	"fmt"
	"github.com/psiayn/heiko/internal/scheduler"
	"github.com/psiayn/heiko/internal/config"
)

func main() {
	confPath := os.Args[1]
	configuration := config.ReadConfig(confPath)
	task_arr := configuration.Jobs
	nodes := configuration.Nodes
	fmt.Println("len of nodes = ", len(task_arr))
	tasks := make(chan config.Task)
	var wg sync.WaitGroup
	wg.Add(len(task_arr))
	go scheduler.Schedule(tasks, nodes, &wg)
	for _, task := range task_arr {
		tasks <- task
	}
	wg.Wait()
}
