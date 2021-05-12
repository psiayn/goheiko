package connection

import (
	"fmt"
	"github.com/psiayn/heiko/internal/config"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Connect(node config.Node) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: node.Username,
		Auth: []ssh.AuthMethod{ssh.Password(node.Password)},
	}

	log.Println("Connecting to node .....")
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", node.Host, node.Port), sshConfig)
	if err != nil {
		log.Println("ERROR while connecting to node: ", err)
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
		log.Println("ERROR while opening output file for task: ", err)
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
		log.Println("ERROR while creating SSH session: ", err)
		return err
	}

	// concatenate commands with semicolon
	combinedCommand := strings.Join(commands, "; ")

	out, err := session.CombinedOutput(combinedCommand)
	if err != nil {
		log.Println("ERROR while running command: ", err)
		return err
	}
	_, err = f.Write(out)
	if err != nil {
		log.Println("ERROR while writing output to file: ", err)
		return err
	}

	return nil
}
