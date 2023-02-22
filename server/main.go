package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "CS598FTS-Warmup/mwmr"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement mwmr.MWMRServer.
type server struct {
	pb.UnimplementedMWMRServer
}

// GetPhase implements mwmr.MWMRServer.
func (s *server) GetPhase(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	// TODO: Server GetPhase implementation
	log.Printf("Receive GetRequest: %s | %s | %d %d \n", in.GetKey(), in.GetValue(), in.GetTs().Ts, in.GetTs().ClientID)
	return &pb.GetReply{
		Key:   in.GetKey(),
		Value: in.GetValue(),
		Ts:    in.GetTs(),
	}, nil
}

// SetPhase implements mwmr.MWMRServer.
func (s *server) SetPhase(ctx context.Context, in *pb.SetRequest) (*pb.SetACK, error) {
	// TODO: Server SetPhase implementation
	log.Printf("Receive SetRequest: %s | %s | %d %d \n", in.GetKey(), in.GetValue(), in.GetTs().Ts, in.GetTs().ClientID)
	return &pb.SetACK{
		Key: in.GetKey(),
		Ts:  in.GetTs(),
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMWMRServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
