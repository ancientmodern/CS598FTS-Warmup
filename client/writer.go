package main

import (
	. "CS598FTS-Warmup/common"
	pb "CS598FTS-Warmup/mwmr"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

func getNewTimestamp(arr []Timestamp) Timestamp {
	var maxTime int64 = 0
	for _, ts := range arr {
		if ts.Time > maxTime {
			maxTime = ts.Time
		}
	}
	return Timestamp{
		Time: maxTime + 1,
		Cid:  *cid,
	}
}

func writerGetPhase(key string) Timestamp {
	f := 2
	n := 2*f + 1

	var wg sync.WaitGroup
	wg.Add(n)

	ch := make(chan Timestamp, n)

	for i := 0; i < n; i++ {
		go func(rid int) {
			conn, err := grpc.Dial(replicas[rid], grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewMWMRClient(conn)

			// Contact the server and print out its response.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			getReply, err := c.GetPhase(ctx, &pb.GetRequest{
				Key: key,
			})
			if err != nil {
				log.Printf("Writer %d getphase from replica %d failed: %v", cid, rid, err)
			} else {
				temp := Timestamp{
					Time: getReply.GetTime(),
					Cid:  getReply.GetCid(),
				}
				ch <- temp
			}
		}(i)
	}

	done := make([]Timestamp, 0, f+1)
	for len(done) < f+1 {
		select {
		case ts := <-ch:
			done = append(done, ts)
			wg.Done()
		default:

		}
	}

	wg.Wait()

	return getNewTimestamp(done)
}

func write(key, value string) bool {
	newTimestamp := writerGetPhase(key)
	// TODO: writerSetPhase
	return true
}
