package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/psiayn/heiko/internal/config"
	"github.com/psiayn/heiko/internal/connection"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

func setAuth() error {
	// TODO
	// Check if Default Heiko Keys exist
	// If not, create
	// Set as default key
	defaultKeyPath := ""

	for i, node := range configuration.Nodes {
		switch strings.ToUpper(node.Auth.Method) {
		case "KEYS":
			publicKeyPath := ""
			privateKeyPath := ""

			// Use Heiko default key if path not specified
			if node.Auth.Keys.Path == "" {
				publicKeyPath = defaultKeyPath + ".pub"
				privateKeyPath = defaultKeyPath
			} else {
				publicKeyPath = node.Auth.Keys.Path + ".pub"
				privateKeyPath = node.Auth.Keys.Path
			}

			// Validate that keys exist
			if _, err := os.Stat(publicKeyPath); err != nil {
				return fmt.Errorf("init: SSH Key %v for node %v does not exist: %v", publicKeyPath, node.Name, err.Error())
			}

			if _, err := os.Stat(privateKeyPath); err != nil {
				return fmt.Errorf("init: SSH Key %v for node %v does not exist: %v", privateKeyPath, node.Name, err.Error())
			}

			// Read and set public key
			pub, err := ioutil.ReadFile(publicKeyPath)
			if err != nil {
				return fmt.Errorf("init: read public key: %v", err.Error())
			}

			publicKey, _, _, _, err := ssh.ParseAuthorizedKey(pub)
			if err != nil {
				return fmt.Errorf("init: parse public key: %v", err.Error())
			}
			configuration.Nodes[i].Auth.Keys.PublicKey = publicKey

			// Read and set private key
			priv, err := ioutil.ReadFile(privateKeyPath)
			if err != nil {
				return fmt.Errorf("init: read private key: %v", err.Error())
			}

			privateKey, err := ssh.ParsePrivateKey(priv)
			if err != nil {
				return fmt.Errorf("init: parse private key: %v", err.Error())
			}
			configuration.Nodes[i].Auth.Keys.PrivateKey = privateKey

		case "PASSWORD":
			// Nothing to do in this case right?
		}
	}
	return nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Runs initialization of Jobs",
	Run: func(cmd *cobra.Command, args []string) {

		err := setAuth()
		if err != nil {
			log.Fatalln(err)
		}

		jobs := configuration.Jobs
		nodes := configuration.Nodes

		var wg sync.WaitGroup

		for _, job := range jobs {
			// we run initialization of this job for each node
			//    for each node, we run it in a separate goroutine
			wg.Add(len(nodes))

			for _, node := range nodes {

				go func(job config.Task, node config.Node) {
					defer wg.Done()
					log.Printf("Running initialization for task %s on node %s",
						job.Name,
						node.Name)
					connection.RunTask(node, job.Name, job.Init)
				}(job, node)

			}
			wg.Wait()
		}
	},
}
