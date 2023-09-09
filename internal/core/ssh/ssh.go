package ssh

import (
	"TTPanel/pkg/util"
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	goSSH "golang.org/x/crypto/ssh"
)

type ConnInfo struct {
	User        string        `json:"user"`
	Address     string        `json:"address"`
	Port        int           `json:"port"`
	Password    string        `json:"password"`
	PrivateKey  []byte        `json:"privateKey"`
	PassPhrase  []byte        `json:"passPhrase"`
	DialTimeOut time.Duration `json:"dialTimeOut"`

	Client     *goSSH.Client  `json:"client"`
	Session    *goSSH.Session `json:"session"`
	LastResult string         `json:"lastResult"`
}

func (c *ConnInfo) NewClient() (*ConnInfo, error) {
	config := &goSSH.ClientConfig{}
	config.SetDefaults()
	addr := fmt.Sprintf("%s:%d", c.Address, c.Port)
	config.User = c.User
	if !util.StrIsEmpty(c.Password) {
		config.Auth = []goSSH.AuthMethod{goSSH.Password(c.Password)}
	} else {
		signer, err := makePrivateKeySigner(c.PrivateKey, c.PassPhrase)
		if err != nil {
			return nil, err
		}
		config.Auth = []goSSH.AuthMethod{goSSH.PublicKeys(signer)}
	}
	if c.DialTimeOut == 0 {
		c.DialTimeOut = 5 * time.Second
	}
	config.Timeout = c.DialTimeOut

	config.HostKeyCallback = goSSH.InsecureIgnoreHostKey()
	client, err := goSSH.Dial("tcp", addr, config)
	if nil != err {
		return c, err
	}
	c.Client = client
	return c, nil
}

func (c *ConnInfo) Run(shell string) (string, error) {
	if c.Client == nil {
		if _, err := c.NewClient(); err != nil {
			return "", err
		}
	}
	session, err := c.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer func(session *goSSH.Session) {
		_ = session.Close()
	}(session)
	buf, err := session.CombinedOutput(shell)

	c.LastResult = string(buf)
	return c.LastResult, err
}

func (c *ConnInfo) Close() {
	_ = c.Client.Close()
}

type ConnClient struct {
	StdinPipe   io.WriteCloser
	ComboOutput *wsBufferWriter
	Session     *goSSH.Session
}

func (c *ConnInfo) NewSshConn(cols, rows int) (*ConnClient, error) {
	sshSession, err := c.Client.NewSession()
	if err != nil {
		return nil, err
	}

	stdinP, err := sshSession.StdinPipe()
	if err != nil {
		return nil, err
	}

	comboWriter := new(wsBufferWriter)
	sshSession.Stdout = comboWriter
	sshSession.Stderr = comboWriter

	modes := goSSH.TerminalModes{
		goSSH.ECHO:          1,
		goSSH.TTY_OP_ISPEED: 14400,
		goSSH.TTY_OP_OSPEED: 14400,
	}
	if err := sshSession.RequestPty("xterm", rows, cols, modes); err != nil {
		return nil, err
	}
	if err := sshSession.Shell(); err != nil {
		return nil, err
	}
	return &ConnClient{StdinPipe: stdinP, ComboOutput: comboWriter, Session: sshSession}, nil
}

func (s *ConnClient) Close() {
	if s.Session != nil {
		_ = s.Session.Close()
	}
}

type wsBufferWriter struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (w *wsBufferWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

func makePrivateKeySigner(privateKey []byte, passPhrase []byte) (goSSH.Signer, error) {
	var signer goSSH.Signer
	if passPhrase != nil {
		s, err := goSSH.ParsePrivateKeyWithPassphrase(privateKey, passPhrase)
		if err != nil {
			return nil, fmt.Errorf("error parsing SSH key: '%v'", err)
		}
		signer = s
	} else {
		s, err := goSSH.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("error parsing SSH key: '%v'", err)
		}
		signer = s
	}

	return signer, nil
}
