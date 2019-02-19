package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type TelnetServer struct {
	c *ShellConfig
}

func newTelnetServer(c *ShellConfig) (*TelnetServer, error) {
	return &TelnetServer{c}, nil
}

func (t *TelnetServer) Start() error {
	ln, err := net.Listen("tcp", t.c.Listen)
	if err != nil {
		return fmt.Errorf("failed to listen: %s", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept: %s", err)
			continue
		}
		r := bufio.NewReader(conn)
		go func() {
			d := &DummyShell{
				Prompt:   t.c.Prompt,
				W:        conn,
				R:        &ConnReadLiner{r},
				Commands: t.c.Commands,
			}
			d.Start()
		}()
	}
}

type ConnReadLiner struct {
	r *bufio.Reader
}

func (c *ConnReadLiner) ReadLine() (string, error) {
	b, _, err := c.r.ReadLine()
	return string(b), err
}
