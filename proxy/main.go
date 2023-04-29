package main

import (
	"flag"
	"fmt"
	"log"
)

type Proxy interface {
	Init()
	Run()
	Stop()
}

type ProxyFactory func(string) Proxy

// TODO: read from configuration file
var (
	socketPath     = "/tmp/sdn_uds.sock"
	proxyFactories = make(map[string]ProxyFactory)
	proxyName      = flag.String("type", "simple", "the name of the proxy to run")
	cid            = flag.Int64("cid", 0, "the id of this client")
	f              = 1
	n              = 2*f + 1
	ErrorLogger    *log.Logger
)

func main() {
	flag.Parse()

	factory, ok := proxyFactories[*proxyName]
	if !ok {
		ErrorLogger.Printf("Invalid proxy name: %s", *proxyName)
		return
	}

	proxy := factory(socketPath)

	proxy.Init()

	fmt.Printf("Run Proxy: cid = %d, type = %s", *cid, *proxyName)

	proxy.Run()
}

func registerProxyFactory(name string, factory ProxyFactory) {
	proxyFactories[name] = factory
}
