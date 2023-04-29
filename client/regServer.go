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

type RegServer struct {
	replicas      []string
	grpcClients   []pb.MWMRClient
	serverAddress string
	running       bool
}

func NewRegServer(serverAddr string) Server {
	return &RegServer{
		replicas:      []string{"node-1:50051", "node-2:50051", "node-3:50051"},
		grpcClients:   make([]pb.MWMRClient, n),
		serverAddress: serverAddr,
		running:       true,
	}
}

func (s *RegServer) Stop() {
	s.running = false
}

func (s *RegServer) Init() {
	for rid := 0; rid < n; rid++ {
		conn, err := grpc.Dial(s.replicas[rid], grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		s.grpcClients[rid] = pb.NewMWMRClient(conn)
	}
}

func (s *RegServer) handleConnection(conn net.Conn) {
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

func (s *RegServer) Run() {
	// Clean up the socket file if it exists
	if _, err := os.Stat(s.serverAddress); !os.IsNotExist(err) {
		os.Remove(s.serverAddress)
	}

	listener, err := net.Listen("unix", s.serverAddress)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("Listen on %s\n", s.serverAddress)

	for s.running {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go s.handleConnection(conn)
	}
}
