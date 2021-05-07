package config

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/creasty/defaults"
)

type Node struct {
	Name string
	Host string
	Username string
	Password string
	Port int
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

func ReadConfig(configPath string) Config {
	dat, e := ioutil.ReadFile(configPath)
	if e != nil {
		fmt.Println("Failed to read file!")
		panic(e)
	}

	config := Config{}
	err := yaml.UnmarshalStrict(dat, &config)
	if err != nil {
		fmt.Println("YABE! Failed to unmarshal YAML")
		panic(err)
	}

	if err := defaults.Set(&config); err != nil {
		panic(fmt.Errorf("panik: Could not set defaults %v", err))
	}

	return config
}
