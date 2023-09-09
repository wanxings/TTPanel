package terminal

import (
	"TTPanel/internal/global"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"errors"
	"github.com/creack/pty"
)

const (
	DefaultCloseSignal  = syscall.SIGINT
	DefaultCloseTimeout = 10 * time.Second
)

type LocalCommand struct {
	closeSignal  syscall.Signal
	closeTimeout time.Duration

	cmd       *exec.Cmd
	pty       *os.File
	ptyClosed chan struct{}
}

func NewCommand(commands string) (*LocalCommand, error) {
	cmd := exec.Command("sh", "-c", commands)

	ptyD, err := pty.Start(cmd)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("NewCommand->Start->failed to start command:%s", err.Error()))
	}
	ptyClosed := make(chan struct{})

	lCmd := &LocalCommand{
		closeSignal:  DefaultCloseSignal,
		closeTimeout: DefaultCloseTimeout,

		cmd:       cmd,
		pty:       ptyD,
		ptyClosed: ptyClosed,
	}

	return lCmd, nil
}

func (l *LocalCommand) Read(p []byte) (n int, err error) {
	return l.pty.Read(p)
}

func (l *LocalCommand) Write(p []byte) (n int, err error) {
	return l.pty.Write(p)
}

func (l *LocalCommand) Close() error {
	if l.cmd != nil && l.cmd.Process != nil {
		_ = l.cmd.Process.Signal(l.closeSignal)
	}
	for {
		select {
		case <-l.ptyClosed:
			return nil
		case <-l.closeTimeoutC():
			_ = l.cmd.Process.Signal(syscall.SIGKILL)
		}
	}
}

func (l *LocalCommand) ResizeTerminal(width int, height int) error {
	window := struct {
		row uint16
		col uint16
		x   uint16
		y   uint16
	}{
		uint16(height),
		uint16(width),
		0,
		0,
	}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		l.pty.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&window)),
	)
	if errno != 0 {
		return errno
	} else {
		return nil
	}
}

func (l *LocalCommand) Wait(quitChan chan bool) {
	if err := l.cmd.Wait(); err != nil {
		global.Log.Errorf("ssh session wait failed, err: %v", err)
		setQuit(quitChan)
	}
}

func (l *LocalCommand) closeTimeoutC() <-chan time.Time {
	if l.closeTimeout >= 0 {
		return time.After(l.closeTimeout)
	}

	return make(chan time.Time)
}
