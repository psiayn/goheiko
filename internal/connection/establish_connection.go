package connection

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/psiayn/heiko/internal/config"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

func publicKey(path string) []ssh.AuthMethod {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		path = filepath.Join(home, ".ssh/heiko/key")
	}

	key, err := ioutil.ReadFile(path)
	if err != nil {
		return []ssh.AuthMethod{}
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return []ssh.AuthMethod{}
	}

	return []ssh.AuthMethod{ssh.PublicKeys(signer)}
}

func Connect(node config.Node) (*ssh.Client, error) {

	sshConfig := &ssh.ClientConfig{
		User: node.Username,
	}

	switch strings.ToUpper(node.Auth.Method) {
	case "PASSWORD":
		sshConfig.Auth = []ssh.AuthMethod{ssh.Password(node.Auth.Password)}

	case "KEYS":
		sshConfig.Auth = publicKey(node.Auth.Keys.Path)
	}

	log.Printf("Connecting to node %s .....", node.Name)
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", node.Host, node.Port), sshConfig)
	if err != nil {
		log.Printf("ERROR while connecting to node %s: %v", node.Name, err)
		return nil, err
	}
	return client, nil
}

func RunTask(node config.Node, name string, commands []string) error {
	// by default, this will be ~/.heiko/<name>/out_<task>
	f_name := filepath.Join(
		viper.GetString("dataLocation"),
		viper.GetString("name"),
		"out_"+name,
	)
	f, err := os.OpenFile(f_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf(
			"ERROR (task: %s, node: %s) while opening output file: %v",
			name,
			node.Name,
			err,
		)
		return err
	}
	defer f.Close()

	client, err := Connect(node)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Printf(
			"ERROR (task: %s, node: %s) while creating SSH session: %v",
			name,
			node.Name,
			err,
		)
		return err
	}

	// concatenate commands with semicolon
	combinedCommand := strings.Join(commands, "; ")

	out, err := session.CombinedOutput(combinedCommand)
	if err != nil {
		log.Printf(
			"ERROR (task: %s, node: %s) while running command: %v",
			name,
			node.Name,
			err,
		)
		return err
	}
	_, err = f.Write(out)
	if err != nil {
		log.Printf(
			"ERROR (task: %s, node: %s) while writing output to file: %v",
			name,
			node.Name,
			err,
		)
		return err
	}

	return nil
}
