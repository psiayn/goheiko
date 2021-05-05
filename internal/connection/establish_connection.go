package connection

import (
	"fmt"
	"github.com/psiayn/heiko/internal/config"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

func Connect(node config.Node, task config.Task) {
	sshConfig := &ssh.ClientConfig{
		User: node.Username,
		Auth: []ssh.AuthMethod{ssh.Password(node.Password)},
	}

	f_name := "out_" + task.Name
	f, err := os.OpenFile(f_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("ERROR while opening output file: ", err)
		return
	}
	defer f.Close()

	fmt.Println("Connecting to node .....")
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", node.Host, node.Port), sshConfig)
	if err != nil {
		fmt.Println("ERROR while connecting to node: ", err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Println("ERROR while creating SSH session: ", err)
		return
	}

	combinedCommand := strings.Join(task.Commands, "; ")

	out, err := session.CombinedOutput(combinedCommand)
	if err != nil {
		fmt.Println("ERROR while running command: ", err)
		return
	}
	_, err = f.Write(out)
	if err != nil {
		fmt.Println("ERROR while writing output to file: ", err)
		return
	}
}
