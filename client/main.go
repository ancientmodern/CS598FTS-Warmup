package main

import (
	pb "CS598FTS-Warmup/mwmr"
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
	"fmt"
  	"math/rand"
	"strconv"
)

const (
	defaultNumRead = 10
	defaultNumWrite = 10
	defaultNumInitial = 10
	defaultCid = 0
)

var (
	replicas = []string{"localhost:50051", "localhost:50052", "localhost:50053", "localhost:50053", "localhost:50053"}
	cid      = flag.Int64("cid", defaultCid, "the id of this client")
	numRead = flag.Int("numRead", defaultNumRead, "Number of Reads")
	numWrite = flag.Int("numWrite", defaultNumWrite, "Number of Writes")
	numInitial = flag.Int("numInitial", defaultNumInitial, "Number of Initialized Pairs")
	f        = 0
	total_sets = 0
	total_gets = 0
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
	// read("111")
	// write("111", "333")
	// read("111")
	// rand.Seed(strconv.Atoi(*numSeed))
	// for i := 0; i < strconv.Atoi(*numRead); i++ {
	// 	read(strconv.Itoa(rand.Intn(1000000)))
	// }
	// for i := 0; i < strconv.Atoi(*numWrite); i++ {
	// 	write(strconv.Itoa(rand.Intn(1000000)))
	// }
	rand.Seed(*cid)
	start_time := time.Now()
	for i := 0; i < *numRead; i++ {
		read(strconv.Itoa(rand.Intn(*numInitial)))
	}
	for i := 0; i < *numWrite; i++ {
		write(strconv.Itoa(rand.Intn(*numInitial)), strconv.Itoa(rand.Intn(*numInitial)))
	}

	end_time := time.Now()
	used_time := end_time.Sub(start_time)
	fmt.Println("Number", *cid, "client start time:", start_time)
	fmt.Println("Number", *cid, "client end time:", end_time)
	fmt.Println("Number", *cid, "client used time:", used_time)
	log.Printf("Number %d #total_sets done: %d\n", *cid, total_sets)
	log.Printf("Number %d #total_gets done: %d\n", *cid, total_gets)
}

