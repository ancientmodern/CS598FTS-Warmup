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
	defaultNumRead = "5000"
	defaultNumWrite = "5000"
	defaultClientNum = "0"
)

var (
	replicas = []string{"localhost:50051", "localhost:50052", "localhost:50053", "localhost:50053", "localhost:50053"}
	cid      = flag.Int64("id", 0, "the id of this client")
	f        = 2
	numRead = flag.String("numRead", defaultNumRead, "Number of Reads")
	numWrite = flag.String("numWrite", defaultNumWrite, "Number of Writes")
	clientNum = flag.String("clientNum", defaultClientNum, "Client Number")
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
	int_numSeed, _ := strconv.ParseInt(*clientNum, 10, 64)
	rand.Seed(int_numSeed)
	int_numRead, _ := strconv.Atoi(*numRead)
	int_numWrite, _ := strconv.Atoi(*numWrite)
	start_time := time.Now()
	for i := 0; i < int_numRead; i++ {
		read(strconv.Itoa(rand.Intn(1000000)))
	}
	for i := 0; i < int_numWrite; i++ {
		write(strconv.Itoa(rand.Intn(1000000)), strconv.Itoa(rand.Intn(1000000)))
	}
	end_time := time.Now()
	used_time := end_time.Sub(start_time)
	fmt.Println("Number", *clientNum, "client start time:", start_time)
	fmt.Println("Number", *clientNum, "client end time:", end_time)
	fmt.Println("Number", *clientNum, "client used time:", used_time)
	log.Printf("Number %s #total_sets done: %d\n", *clientNum, total_sets)
	log.Printf("Number %s #total_gets done: %d\n", *clientNum, total_gets)
}

