package main

import (
	pb "CS598FTS-Warmup/mwmr"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/montanaflynn/stats"
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

	// replicas      = []string{"localhost:50051"}
	cid           = flag.Int64("cid", defaultCid, "the id of this client")
	numRead       = flag.Int("numRead", defaultNumRead, "Number of Reads")
	numWrite      = flag.Int("numWrite", defaultNumWrite, "Number of Writes")
	numInitial    = flag.Int("numInitial", defaultNumInitial, "Number of Initialized Pairs")
	f             = 2
	n             = 2*f + 1
	grpcClient    = make([]pb.MWMRClient, n)
	totalSets     = 0
	totalGets     = 0
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	read_ts       []float64
	read_get_ts   []float64
	read_set_ts   []float64
	write_ts      []float64
	write_get_ts  []float64
	write_set_ts  []float64
)

func main() {

	flag.Parse()
	rand.Seed(*cid)

	file, err := os.OpenFile(fmt.Sprintf("logs/logs_client%d.txt", *cid), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	for rid := 0; rid < n; rid++ {
		conn, err := grpc.Dial(replicas[rid], grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		grpcClient[rid] = pb.NewMWMRClient(conn)
	}

	start := time.Now().UnixNano()
	for i := 0; i < *numRead; i++ {
		_, read_get_t, read_set_t := read(strconv.Itoa(rand.Intn(*numInitial)))
		read_get_ts = append(read_get_ts, float64(read_get_t))
		read_set_ts = append(read_set_ts, float64(read_set_t))
		read_ts = append(read_ts, float64(read_get_t+read_set_t))
		// InfoLogger.Printf("ith read: %d, read_get_t: %d ns, read_set_t: %d ns, read_t: %d ns", i, read_get_t, read_set_t, read_get_t+read_set_t)

	}
	for i := 0; i < *numWrite; i++ {
		write_get_t, write_set_t := write(strconv.Itoa(rand.Intn(*numInitial)), strconv.Itoa(rand.Intn(*numInitial)))
		write_get_ts = append(write_get_ts, float64(write_get_t))
		write_set_ts = append(write_set_ts, float64(write_set_t))
		write_ts = append(write_ts, float64(write_get_t+write_set_t))
		// InfoLogger.Printf("ith write: %d, write_get_t: %d ns, write_set_t: %d ns, write_t: %d ns", i, write_get_t, write_set_t, write_get_t+write_set_t)
	}
	end := time.Now().UnixNano()

	InfoLogger.Println("====================================================================================================")
	InfoLogger.Printf("Number %d #total_sets done: %d\n", *cid, totalSets)
	InfoLogger.Printf("Number %d #total_gets done: %d\n", *cid, totalGets)
	logMean()
	logPercentile(50)
	logPercentile(99)
	InfoLogger.Printf("start time: %d, end time: %d \n", start, end)

	if err := file.Close(); err != nil {
		ErrorLogger.Fatal(err)
	}
}

func calcMean(lst []float64) float64 {
	if len(lst) == 0 {
		return 0
	}
	res, err := stats.Mean(lst)
	if err != nil {
		ErrorLogger.Fatal(err)
	}
	return res
}

func logMean() {
	InfoLogger.Println("====================================================")

	InfoLogger.Printf("Mean latency of read: %f ns", calcMean(read_ts))
	InfoLogger.Printf("Mean latency of read get: %f ns", calcMean(read_get_ts))
	InfoLogger.Printf("Mean latency of read set: %f ns", calcMean(read_set_ts))

	InfoLogger.Printf("Mean latency of write: %f ns", calcMean(write_ts))
	InfoLogger.Printf("Mean latency of write get: %f ns", calcMean(write_get_ts))
	InfoLogger.Printf("Mean latency of write get: %f ns", calcMean(write_set_ts))
}

func calcPercentile(lst []float64, percentile float64) float64 {
	if len(lst) == 0 {
		return 0
	}
	res, err := stats.Percentile(lst, percentile)
	if err != nil {
		ErrorLogger.Fatal(err)
	}
	return res
}

func logPercentile(percentile float64) {
	InfoLogger.Println("====================================================")
	InfoLogger.Printf("Percentile %f latency of read: %f ns\n", percentile, calcPercentile(read_ts, percentile))
	InfoLogger.Printf("Percentile %f latency of read get: %f ns\n", percentile, calcPercentile(read_get_ts, percentile))
	InfoLogger.Printf("Percentile %f latency of read set: %f ns\n", percentile, calcPercentile(read_set_ts, percentile))

	InfoLogger.Printf("Percentile %f latency of write: %f ns\n", percentile, calcPercentile(write_ts, percentile))
	InfoLogger.Printf("Percentile %f latency of write get: %f ns\n", percentile, calcPercentile(write_get_ts, percentile))
	InfoLogger.Printf("Percentile %f latency of write get: %f ns\n", percentile, calcPercentile(write_set_ts, percentile))
}
