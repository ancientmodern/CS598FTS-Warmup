package main

import (
	. "CS598FTS-Warmup/common"
	pb "CS598FTS-Warmup/mwmr"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"google.golang.org/grpc"
)

var (
	PORT       = flag.Int("port", 50051, "The server port")
	SIZE_STORE = flag.Int("size", 10, "num of entries in the store")
	kvStore    = make(map[string]*Pair_m)
	s          *grpc.Server
)

// server is used to implement mwmr.MWMRServer.
type server struct {
	pb.UnimplementedMWMRServer
}

// GetPhase implements mwmr.MWMRServer.
func (s *server) GetPhase(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	key := in.GetKey()
	kvStore[key].Mtx.RLock()
	val := kvStore[key].Value
	t := kvStore[key].Ts.Time
	cid := kvStore[key].Ts.Cid
	kvStore[key].Mtx.RUnlock()
	return &pb.GetReply{
		Value: val,
		Time:  t,
		Cid:   cid,
	}, nil
}

// SetPhase implements mwmr.MWMRServer.
func (s *server) SetPhase(ctx context.Context, in *pb.SetRequest) (*pb.SetACK, error) {
	key := in.GetKey()
	val := in.GetValue()
	time := in.GetTime()
	cid := in.GetCid()

	kvStore[key].Mtx.RLock()
	time_store := kvStore[key].Ts.Time
	cid_store := kvStore[key].Ts.Cid
	kvStore[key].Mtx.RUnlock()

	if time < time_store || (time == time_store && cid < cid_store) {
		return &pb.SetACK{
			Applied: false,
		}, nil
	}

	kvStore[key].Mtx.Lock()
	kvStore[key].Value = val
	kvStore[key].Ts.Time = time
	kvStore[key].Ts.Cid = cid
	kvStore[key].Mtx.Unlock()

	return &pb.SetACK{
		Applied: true,
	}, nil
}

func main() {
	flag.Parse()

	// initialize the kvStore
	for i := 0; i < *SIZE_STORE; i++ {
		kvStore[strconv.Itoa(i)] = &Pair_m{Value: "init", Ts: Timestamp{Time: -1, Cid: -1}, Mtx: sync.RWMutex{}}
	}

	s = grpc.NewServer()
	pb.RegisterMWMRServer(s, &server{})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
