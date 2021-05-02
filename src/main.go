package main

import (
	"os"
	"sync"
	"github.com/psiayn/goheiko/utils"
)

func main() {
	confPath := os.Args[1]
	config := utils.ReadConfig(confPath)
	task_arr := config.Jobs
	nodes := config.Nodes
	tasks := make(chan utils.Task)
	var wg sync.WaitGroup
	wg.Add(len(task_arr))
	go utils.Schedule(tasks, nodes, &wg)
	for _, task := range task_arr {
		tasks <- task
	}
	wg.Wait()
}
