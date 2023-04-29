package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type SimpleProxy struct {
	macToPort     map[string]byte
	socketAddress string
	running       bool
	mu            sync.Mutex
}

func NewSimpleProxy(socketAddr string) Proxy {
	return &SimpleProxy{
		macToPort:     make(map[string]byte),
		socketAddress: socketAddr,
		running:       true,
	}
}

func (s *SimpleProxy) Init() {

}

func (s *SimpleProxy) Stop() {
	s.running = false
}

func (s *SimpleProxy) handleConnection(conn net.Conn) {
	defer conn.Close()

	response := make([]byte, 10)
	_, err := conn.Read(response)
	if err != nil {
		return
	}

	setByte := response[0]
	dpid := binary.BigEndian.Uint16(response[1:3])
	macAddress := decodeMacAddress(response[3:9])
	val := response[9]
	key := response[1:9]

	keyStr := string(key)

	s.mu.Lock()
	defer s.mu.Unlock()

	if setByte == 0x00 {
		// Get request
		val, ok := s.macToPort[keyStr]
		fmt.Printf("GET: dpid = %d, mac_address = %s, get_val: %d\n", dpid, macAddress, val)

		if !ok {
			val = 0xFF // 0xFF means key does not exist
		}
		conn.Write([]byte{val})
	} else {
		// Set request
		fmt.Printf("SET: dpid = %d, mac_address = %s, set_val: %d\n", dpid, macAddress, val)
		s.macToPort[keyStr] = val
	}
}

func (s *SimpleProxy) Run() {
	// Clean up the socket file if it exists
	if _, err := os.Stat(s.socketAddress); !os.IsNotExist(err) {
		os.Remove(s.socketAddress)
	}

	listener, err := net.Listen("unix", s.socketAddress)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("Listen on %s\n", s.socketAddress)

	for s.running {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go s.handleConnection(conn)
	}
}

func decodeMacAddress(macBytes []byte) string {
	hexMac := hex.EncodeToString(macBytes)
	return strings.Join([]string{hexMac[:2], hexMac[2:4], hexMac[4:6], hexMac[6:8], hexMac[8:10], hexMac[10:12]}, ":")
}

func init() {
	registerProxyFactory("simple", NewSimpleProxy)
}
