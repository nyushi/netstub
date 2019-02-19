package main

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type SSHServer struct {
	c *ShellConfig
}

func newSSHServer(c *ShellConfig) (*SSHServer, error) {
	return &SSHServer{c}, nil
}
func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func (t *SSHServer) Start() error {
	s := &ssh.Server{
		Addr: t.c.Listen,
		Handler: func(s ssh.Session) {
			term := terminal.NewTerminal(s, "")
			s.Pty()
			d := &DummyShell{
				Prompt:   t.c.Prompt,
				W:        term,
				R:        term,
				Commands: t.c.Commands,
			}
			d.Start()
		},
		PasswordHandler: func(ctx ssh.Context, password string) bool { return true },
	}
	return s.ListenAndServe()
}
