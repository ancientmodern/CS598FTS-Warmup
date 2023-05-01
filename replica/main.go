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
	"sync"
)

var (
	IP      = flag.String("ip", "0.0.0.0", "the replica ip address")
	PORT    = flag.Int("port", 50051, "The replica port")
	kvStore = sync.Map{}
	s       *grpc.Server
)

// replica is used to implement mwmr.MWMRServer.
type replica struct {
	pb.UnimplementedMWMRServer
}

// GetPhase implements mwmr.MWMRServer.
func (s *replica) GetPhase(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	key := in.GetKey()

	pair, ok := kvStore.Load(key)
	if !ok {
		fmt.Println("Key not found:", key)

		return &pb.GetReply{
			Value: 0xFF,
			Time:  0,
			Cid:   0,
		}, nil
	}

	res := pair.(*Pair)
	val := res.Value
	t := res.Ts.Time
	cid := res.Ts.Cid

	fmt.Println("Get the key:", key, "with value:", val)

	return &pb.GetReply{
		Value: val,
		Time:  t,
		Cid:   cid,
	}, nil
}

// SetPhase implements mwmr.MWMRServer.
func (s *replica) SetPhase(ctx context.Context, in *pb.SetRequest) (*pb.SetACK, error) {
	key := in.GetKey()
	val := in.GetValue()
	time := in.GetTime()
	cid := in.GetCid()

	pair, ok := kvStore.Load(key)
	if !ok {
		fmt.Println("Insert a new key:", key, "with value:", val)
		pair = &Pair{
			Value: val,
			Ts: Timestamp{
				Time: time,
				Cid:  cid,
			},
		}
		kvStore.Store(key, pair)

		return &pb.SetACK{
			Applied: true,
		}, nil
	}

	res := pair.(*Pair)
	timeStore := res.Ts.Time
	cidStore := res.Ts.Cid

	if time < timeStore || (time == timeStore && cid < cidStore) {
		return &pb.SetACK{
			Applied: false,
		}, nil
	}

	insert := &Pair{
		Value: val,
		Ts: Timestamp{
			Time: time,
			Cid:  cid,
		},
	}
	kvStore.Store(key, insert)

	fmt.Println("Update the key:", key, "with value:", val)

	return &pb.SetACK{
		Applied: true,
	}, nil
}

func main() {
	flag.Parse()

	s = grpc.NewServer()
	pb.RegisterMWMRServer(s, &replica{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("replica listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
