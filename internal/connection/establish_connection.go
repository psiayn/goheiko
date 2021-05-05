package connection

import (
	"fmt"
	"os"
	"golang.org/x/crypto/ssh"
	"github.com/psiayn/heiko/internal/config"
)

func Connect(node config.Node, task config.Task) {
	sshConfig := &ssh.ClientConfig{
		User: node.Username,
		Auth: []ssh.AuthMethod{ssh.Password(node.Password)},
	}
	f_name := "out_" + node.Name + "," + task.Name
	f, e := os.OpenFile(f_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if e != nil {
		fmt.Println("ERROR ", e)
		return
	}
	defer f.Close()
	fmt.Println("Connecting to node .....")
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, err := ssh.Dial("tcp", node.Host, sshConfig)
	if err != nil {
		fmt.Println("ERROR ", err)
		return
	}

	session, error := client.NewSession()
	if error != nil {
		fmt.Println("ERROR ", error)
		client.Close()
		return
	}
	command := task.Commands[1]
	out, er := session.CombinedOutput(command)
	if er != nil {
		fmt.Println("ERROR ", er)
		client.Close()
		return
	}
	_, errr := f.Write(out)
	if errr != nil {
		fmt.Println("ERROR ", er)
		client.Close()
		return
	}
	client.Close()
}
