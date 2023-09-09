package service

import (
	"TTPanel/internal/model"
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

type WebSSHService struct{}

func (s *WebSSHService) NewSshClient(h *model.Host) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            h.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 5,
	}
	if len(h.PrivateKey) > 10 {
		key, err := ssh.ParsePrivateKey([]byte(h.PrivateKey))
		if err != nil {
			return nil, err
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(key)}
	} else {
		config.Auth = []ssh.AuthMethod{ssh.Password(h.Password)}
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", h.Address, h.Port), config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
