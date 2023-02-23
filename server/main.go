package main

import (
	. "CS598FTS-Warmup/common"
	pb "CS598FTS-Warmup/mwmr"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port       = flag.Int("port", 50051, "The server port")
	kvStore    = make(map[string]*Pair)
	NUM_RECORD = 10
)

// server is used to implement mwmr.MWMRServer.
type server struct {
	pb.UnimplementedMWMRServer
	// m map[string]Pair
}

// GetPhase implements mwmr.MWMRServer.
func (s *server) GetPhase(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	// TODO: Server GetPhase implementation
	key := in.GetKey()
	log.Printf("Receive GetRequest: key: %s\n", key)
	return &pb.GetReply{
		Value: kvStore[key].Value,
		Time:  kvStore[key].Ts.Time,
		Cid:   kvStore[key].Ts.Cid,
	}, nil
}

// SetPhase implements mwmr.MWMRServer.
func (s *server) SetPhase(ctx context.Context, in *pb.SetRequest) (*pb.SetACK, error) {
	// TODO: Server SetPhase implementation
	log.Printf("Receive SetRequest: %s | %s | %d %d \n", in.GetKey(), in.GetValue(), in.GetTime(), in.GetCid())

	key := in.GetKey()
	val := in.GetValue()
	time := in.GetTime()
	cid := in.GetCid()

	if time < kvStore[key].Ts.Time || (time == kvStore[key].Ts.Time && cid < kvStore[key].Ts.Cid) {
		return &pb.SetACK{
			Applied: false,
		}, nil
	}

	kvStore[key] = &Pair{Value: val, Ts: Timestamp{Time: time, Cid: cid}}

	//for i := 0; i < NUM_RECORD; i++ {
	//	log.Printf("entry %d, value %s time %d, cid %d", i, kvStore[strconv.Itoa(i)].Value, kvStore[strconv.Itoa(i)].Ts.Time, kvStore[strconv.Itoa(i)].Ts.Cid)
	//}

	return &pb.SetACK{
		Applied: true,
	}, nil
}

func testOnly() {
	kvStore["111"] = &Pair{
		Value: "222",
		Ts: Timestamp{
			Time: 10,
			Cid:  2,
		},
	}
}

func main() {
	flag.Parse()

	testOnly()

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
