package main

import (
	pb "CS598FTS-Warmup/mwmr"
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	replicas = []string{"localhost:50051", "localhost:50052", "localhost:50053"}
	cid      = flag.Int64("id", 0, "the id of this client")
	f        = 1
)

func reference() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(replicas[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMWMRClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// GetPhase example
	getReply, err := c.GetPhase(ctx, &pb.GetRequest{
		Key: "aaa",
	})
	if err != nil {
		log.Fatalf("could not getphase: %v", err)
	}
	log.Printf("GetReply: value: %s, timestamp: %d %d\n", getReply.GetValue(), getReply.GetTime(), getReply.GetCid())
	// ---------------------------------------------------------------

	// SetPhase example
	setReply, err := c.SetPhase(ctx, &pb.SetRequest{
		Key:   "aaa",
		Value: "bbb",
		Time:  0,
		Cid:   0,
	})
	if err != nil {
		log.Fatalf("could not getphase: %v", err)
	}
	log.Printf("SetACK: applied: %t\n", setReply.GetApplied())
	// ---------------------------------------------------------------
}

func main() {
	flag.Parse()
	read("111")
	write("111", "333")
	read("111")
}
