package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

func createKeyPair(privateKeyPath, publicKeyPath string) error {
	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("cannot generate RSA key: %s", err)
	}
	publickey := &privatekey.PublicKey

	// dump private key to file
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	privatePem, err := os.Create(privateKeyPath)
	if err != nil {
		return fmt.Errorf("create private key: %s", err)
	}

	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		return fmt.Errorf("encode private key: %s", err)
	}

	// dump public key to file
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return fmt.Errorf("dumping publickey: %s", err)
	}

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicPem, err := os.Create(publicKeyPath)
	if err != nil {
		return fmt.Errorf("create public key: %s", err)
	}

	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		return fmt.Errorf("encode public key: %s", err)
	}

	return nil
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
		case "KEYS":

			// Set custom key if specified
			if node.Auth.Keys.Path != "" {
				publicKeyPath = node.Auth.Keys.Path + ".pub"
				privateKeyPath = node.Auth.Keys.Path
			} else {
				configuration.Nodes[i].Auth.Keys.Path = defaultKeyPath
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
