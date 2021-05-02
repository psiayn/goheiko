package main

import (
	"fmt"
	"sync"
)
func main() {
	task_arr := []Task{{"echo belieb"}, {"echo in"}, {"echo go"}, {"echo supremacy"}}
	nodes := []Node{
		Node {
			name: "node_1",
			host: "1.2.3.4",
			username: "belieb",
		},
		Node {
			name: "node_2",
			host: "1.2.3.4",
			username: "in",
		},
		Node {
			name: "node_3",
			host: "1.2.3.4",
			username: "go",
		},
		Node {
			name: "node_4",
			host: "1.2.3.4",
			username: "go",
		},
		Node {
			name: "node_5",
			host: "1.2.3.4",
			username: "go",
		},
		Node {
			name: "node_6",
			host: "1.2.3.4",
			username: "go",
		},
		Node {
			name: "node_7",
			host: "1.2.3.4",
			username: "go",
		},
	}
	tasks := make(chan Task)
	var wg sync.WaitGroup
	wg.Add(4)
	go schedule(tasks, nodes, &wg)
	for _, task := range task_arr {
		tasks <- task
	}
	wg.Wait()
	fmt.Println("Belieb in go")
}
