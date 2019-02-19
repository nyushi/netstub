package main

import (
	"io"
	"log"
	"net"
	"time"
)

type Readliner interface {
	ReadLine() (string, error)
}
type DummyShell struct {
	Prompt   string
	W        io.Writer
	R        Readliner
	Commands Commands

	showPrompt bool
	tmperr     int
	continues  []string
}

func (d *DummyShell) Start() {
	d.showPrompt = true
	for {
		if d.showPrompt {
			d.W.Write([]byte(d.Prompt))
		}
		l, err := d.R.ReadLine()
		if err != nil {
			oe, ok := err.(*net.OpError)
			if !ok {
				log.Printf("error at read: %s", err)
				return
			}
			if oe.Temporary() {
				if d.tmperr == 3 {
					log.Printf("tmp error occured many times: %s", oe)
					return
				}
				log.Printf("temporary error, retrying: %s", oe)
				time.Sleep(time.Second)
				d.tmperr++
			}
			log.Printf("network error: %s", err)
			return
		}
		if len(d.continues) > 0 {
			d.W.Write([]byte(d.continues[0]))
			if len(d.continues) == 1 {
				d.continues = nil
			} else {
				d.continues = d.continues[1:]
			}
			continue
		}
		d.tmperr = 0
		cmd := d.Commands.Match(l)
		if cmd == nil {
			continue
		}
		d.showPrompt = !cmd.NoPrompt
		if cmd.ChangePrompt != nil {
			d.Prompt = *cmd.ChangePrompt
		}
		d.W.Write([]byte(cmd.Output + "\n"))
		if cmd.Continues != nil {
			d.continues = cmd.Continues
		}
	}
}
