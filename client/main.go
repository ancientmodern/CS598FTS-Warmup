package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "CS598FTS-Warmup/mwmr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMWMRClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	getReply, err := c.GetPhase(ctx, &pb.GetRequest{
		Key:   "aaa",
		Value: "bbb",
		Ts: &pb.Timestamp{
			Ts:       0,
			ClientID: 1,
		},
	})
	if err != nil {
		log.Fatalf("could not getphase: %v", err)
	}
	log.Printf("GetReply: %s | %s | %d %d \n", getReply.GetKey(), getReply.GetValue(), getReply.GetTs().Ts, getReply.GetTs().ClientID)

	setReply, err := c.SetPhase(ctx, &pb.SetRequest{
		Key:   "aaa",
		Value: "bbb",
		Ts: &pb.Timestamp{
			Ts:       0,
			ClientID: 1,
		},
	})
	if err != nil {
		log.Fatalf("could not getphase: %v", err)
	}
	log.Printf("SetACK: %s | %d %d \n", setReply.GetKey(), setReply.GetTs().Ts, setReply.GetTs().ClientID)
}
