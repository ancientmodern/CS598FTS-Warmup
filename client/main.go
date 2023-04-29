package main

import (
	"flag"
	"log"
)

// TODO: read from configuration file
var (
	socketPath  = "/tmp/sdn_uds.sock"
	cid         = flag.Int64("cid", 0, "the id of this client")
	f           = 1
	n           = 2*f + 1
	ErrorLogger *log.Logger
)

type Server interface {
	Init()
	Run()
	Stop()
}

func main() {
	flag.Parse()

	server := NewRegServer(socketPath)

	server.Init()

	server.Run()
}
