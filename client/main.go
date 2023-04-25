package main

import (
	pb "CS598FTS-Warmup/mwmr"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	replicas = []string{"128.110.217.160:50051", "128.110.217.137:50051", "128.110.217.131:50051"}

	socketPath  = "/tmp/sdn-uds.sock"
	cid         = flag.Int64("cid", 0, "the id of this client")
	f           = 1
	n           = 2*f + 1
	grpcClients = make([]pb.MWMRClient, n)
	ErrorLogger *log.Logger
)

func main() {
	flag.Parse()

	initGrpcConn()

	if err := udsHandler(); err != nil {
		fmt.Println("Error starting UDS server:", err)
		return
	}
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
