package config

import "golang.org/x/crypto/ssh"

type Node struct {
	Name     string
	Host     string
	Username string
	Auth     Auth
	Port     int `default:"22"`
}

type Auth struct {
	Method   string
	Password string
	Keys     SSHKeys
}

type SSHKeys struct {
	PublicKey  ssh.PublicKey
	PrivateKey ssh.Signer
	Path       string
}

type Task struct {
	Name     string
	Init     []string
	Commands []string
	Restart  bool `default:"false"`
}

type Config struct {
	Nodes []Node
	Jobs  []Task
}
