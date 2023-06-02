package sshd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func DefaultClientConfig() *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: "low",
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func DefaultClient(t *testing.T) *ssh.Client {
	client, dErr := ssh.Dial("tcp", "localhost:2222", DefaultClientConfig())
	assert.Nil(t, dErr)
	return client
}

func KeepAlive(client *ssh.Client) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_, _, err := client.SendRequest("keepalive@golang.org", true, nil)
			if err != nil {
				return
			}
		}
	}
}
