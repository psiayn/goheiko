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

// TODO: return err from here, to see if it has to be run again
// TODO: refactor this into functions
func Connect(node config.Node, task config.Task) {
	sshConfig := &ssh.ClientConfig{
		User: node.Username,
		Auth: []ssh.AuthMethod{ssh.Password(node.Password)},
	}

	f_name := filepath.Join(viper.GetString("dataLocation"), viper.GetString("name"), "out_"+task.Name)
	f, err := os.OpenFile(f_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERROR while opening output file: ", err)
		return
	}
	defer f.Close()

	log.Println("Connecting to node .....")
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", node.Host, node.Port), sshConfig)
	if err != nil {
		log.Println("ERROR while connecting to node: ", err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Println("ERROR while creating SSH session: ", err)
		return
	}

	combinedCommand := strings.Join(task.Commands, "; ")

	out, err := session.CombinedOutput(combinedCommand)
	if err != nil {
		log.Println("ERROR while running command: ", err)
		return
	}
	_, err = f.Write(out)
	if err != nil {
		log.Println("ERROR while writing output to file: ", err)
		return
	}
}
