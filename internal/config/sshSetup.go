package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

func createKeyPair(privateKeyPath, publicKeyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Generate and write private key as PEM
	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	// Set permissions to private key
	if err = os.Chmod(privateKeyPath, 0400); err != nil {
		return err
	}

	// Generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(publicKeyPath, ssh.MarshalAuthorizedKey(pub), 0400)
	if err != nil {
		return err
	}

	return nil
}

func knownHost(host string, port int) bool {
	var output []byte
	var err error

	if port != 22 {
		output, err = exec.Command("ssh-keygen", "-H", "-F", fmt.Sprintf("%s:%d", host, port)).CombinedOutput()
	} else {
		output, err = exec.Command("ssh-keygen", "-H", "-F", host).CombinedOutput()
	}

	if err.Error() == "exit status 1" && len(output) > 0 {
		return true
	}

	return false
}

func transferKey(keyPath, username, host string, port int) error {
	command := exec.Command(
		"ssh-copy-id",
		"-i",
		keyPath,
		"-p",
		fmt.Sprintf("%d", port),
		fmt.Sprintf("%s@%s", username, host),
	)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func SetAuth(configuration *Config) error {
	// Get home directory
	homePath, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("init: UserHomeDir: %v", err.Error())
	}

	defaultKeyPath := filepath.Join(homePath, ".ssh/heiko")
	privateKeyPath := filepath.Join(defaultKeyPath, "key")
	publicKeyPath := privateKeyPath + ".pub"

	// Validate that keys exist
	_, err1 := os.Stat(publicKeyPath)
	_, err2 := os.Stat(privateKeyPath)

	// Create fresh keys if either don't exist
	if err1 != nil || err2 != nil {
		os.RemoveAll(defaultKeyPath)
		os.Mkdir(defaultKeyPath, 0755)

		err := createKeyPair(privateKeyPath, publicKeyPath)
		if err != nil {
			return err
		}
	}

	for i, node := range configuration.Nodes {
		switch strings.ToUpper(node.Auth.Method) {
		case "PASSWORD":
			// If Authentication method is password, no INIT required.

		default:
			// Set custom key if specified
			if node.Auth.Keys.Path != "" {
				publicKeyPath = node.Auth.Keys.Path + ".pub"
				privateKeyPath = node.Auth.Keys.Path
			} else {
				configuration.Nodes[i].Auth.Keys.Path = privateKeyPath
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

			// Transfer Key if node not is not a known host
			if !knownHost(node.Host, node.Port) {
				err := transferKey(privateKeyPath, node.Username, node.Host, node.Port)
				if err != nil {
					return fmt.Errorf("key transfer: %v", err)
				}
			}
		}
	}
	return nil
}
