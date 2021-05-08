package config

type Node struct {
	Name string
	Host string
	Username string
	Password string
	Port int `default:"22"`
}

type Task struct {
	Name string
	Init []string
	Commands []string
	Restart bool `default:"false"`
}

type Config struct {
	Nodes []Node
	Jobs  []Task
}
