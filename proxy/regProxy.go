package main

import (
	pb "CS598FTS-Warmup/mwmr"
	"encoding/binary"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
)

type RegProxy struct {
	replicas      []string
	grpcClients   []pb.MWMRClient
	socketAddress string
	running       bool
}

func NewRegProxy(socketAddr string) Proxy {
	return &RegProxy{
		replicas:      []string{"node-1:50051", "node-2:50051", "node-3:50051"},
		grpcClients:   make([]pb.MWMRClient, n),
		socketAddress: socketAddr,
		running:       true,
	}
}

func (s *RegProxy) Stop() {
	s.running = false
}

func (s *RegProxy) Init() {
	for rid := 0; rid < n; rid++ {
		conn, err := grpc.Dial(s.replicas[rid], grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		s.grpcClients[rid] = pb.NewMWMRClient(conn)
	}
}

func (s *RegProxy) handleConnection(conn net.Conn) {
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

	keyUint64 := binary.BigEndian.Uint64(key)
	valUint32 := binary.BigEndian.Uint32([]byte{val})

	if setByte == 0x00 {
		// Get request
		getVal := s.read(keyUint64) // 0xFF means key does not exist
		fmt.Printf("GET: dpid = %d, mac_address = %s, get_val: %d\n", dpid, macAddress, getVal)
		msg := make([]byte, 4)
		binary.BigEndian.PutUint32(msg, getVal)

		conn.Write(msg)
	} else {
		// Set request
		fmt.Printf("SET: dpid = %d, mac_address = %s, set_val: %d\n", dpid, macAddress, val)
		s.write(keyUint64, valUint32)
	}
}

func (s *RegProxy) Run() {
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

func init() {
	registerProxyFactory("reg", NewRegProxy)
	registerProxyFactory("register", NewRegProxy)
}
