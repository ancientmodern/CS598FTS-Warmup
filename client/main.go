package main

import (
	pb "CS598FTS-Warmup/mwmr"
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultNumRead    = 10
	defaultNumWrite   = 10
	defaultNumInitial = 10
	defaultCid        = 0
)

var (
	replicas = []string{"128.110.217.160:50051", "128.110.217.137:50051", "128.110.217.131:50051", "128.110.217.155:50051", "128.110.217.120:50051"}

	cid        = flag.Int64("cid", defaultCid, "the id of this client")
	numRead    = flag.Int("numRead", defaultNumRead, "Number of Reads")
	numWrite   = flag.Int("numWrite", defaultNumWrite, "Number of Writes")
	numInitial = flag.Int("numInitial", defaultNumInitial, "Number of Initialized Pairs")
	f          = 2
	n          = 2*f + 1
	grpcClient = make([]pb.MWMRClient, n)
	totalSets  = 0
	totalGets  = 0
)

func main() {
	flag.Parse()
	rand.Seed(*cid)
	for rid := 0; rid < n; rid++ {
		conn, err := grpc.Dial(replicas[rid], grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		grpcClient[rid] = pb.NewMWMRClient(conn)
	}
	startTime := time.Now()
	for i := 0; i < *numRead; i++ {
		read(strconv.Itoa(rand.Intn(*numInitial)))
	}
	for i := 0; i < *numWrite; i++ {
		write(strconv.Itoa(rand.Intn(*numInitial)), strconv.Itoa(rand.Intn(*numInitial)))
	}

	endTime := time.Now()
	usedTime := endTime.Sub(startTime)
	log.Println("====================================================================================================")
	log.Println("Number", *cid, "client start time:", startTime)
	log.Println("Number", *cid, "client end time:", endTime)
	log.Println("Number", *cid, "client used time:", usedTime)
	log.Printf("Number %d #total_sets done: %d\n", *cid, totalSets)
	log.Printf("Number %d #total_gets done: %d\n", *cid, totalGets)
}
