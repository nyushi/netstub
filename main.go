package main

import (
	"log"
	"os"

	"github.com/k0kubun/pp"
)

func main() {
	path := os.Args[1]
	conf, err := readConf(path)
	if err != nil {
		log.Fatalf("failed to load `%s`: %s", path, err)
	}
	pp.Println(conf)
	var (
		s DummyServer
	)
	switch conf.Shell.Type {
	case "telnet":
		s, err = newTelnetServer(conf.Shell)
	case "ssh":
		s, err = newSSHServer(conf.Shell)
	default:
		log.Fatalf("invalid type %s", conf.Shell.Type)
	}
	if err != nil {
		log.Fatalf("invalid ssh configuration: %s", err)
	}
	s.Start()
}
