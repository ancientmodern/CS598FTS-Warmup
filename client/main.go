package main

import (
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	defaultNumRead    = 10
	defaultNumWrite   = 10
	defaultNumInitial = 10
	defaultCid        = 0
)

var (
	// replicas   = []string{"172.16.50.1:50051", "172.16.50.2:50051", "172.16.50.3:50051", "172.16.50.4:50051", "172.16.50.5:50051"}
	replicas   = []string{"localhost:50051", "localhost:50052", "localhost:50053", "localhost:50054", "localhost:50055"}
	cid        = flag.Int64("cid", defaultCid, "the id of this client")
	numRead    = flag.Int("numRead", defaultNumRead, "Number of Reads")
	numWrite   = flag.Int("numWrite", defaultNumWrite, "Number of Writes")
	numInitial = flag.Int("numInitial", defaultNumInitial, "Number of Initialized Pairs")
	f          = 2
	total_sets = 0
	total_gets = 0
)

func main() {
	flag.Parse()
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
	log.Println("====================================================================================================")
	log.Println("Number", *cid, "client start time:", start_time)
	log.Println("Number", *cid, "client end time:", end_time)
	log.Println("Number", *cid, "client used time:", used_time)
	log.Printf("Number %d #total_sets done: %d\n", *cid, total_sets)
	log.Printf("Number %d #total_gets done: %d\n", *cid, total_gets)
}
