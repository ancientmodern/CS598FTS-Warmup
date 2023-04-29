package main

import (
	pb "CS598FTS-Warmup/mwmr"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// TODO: read from configuration file
var (
	replicas = []string{"node-1:50051", "node-2:50051", "node-3:50051"}

	socketPath  = "/tmp/sdn-uds.sock"
	cid         = flag.Int64("cid", 0, "the id of this client")
	f           = 1
	n           = 2*f + 1
	grpcClients = make([]pb.MWMRClient, n)
	ErrorLogger *log.Logger
)

func main() {
	flag.Parse()

	// initGrpcConn()

	server := NewSimpleServer(socketPath)

	server.Run()
}

func initGrpcConn() {
	for rid := 0; rid < n; rid++ {
		conn, err := grpc.Dial(replicas[rid], grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		grpcClients[rid] = pb.NewMWMRClient(conn)
	}
}
