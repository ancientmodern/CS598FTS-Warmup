package main

import (
	. "CS598FTS-Warmup/common"
	pb "CS598FTS-Warmup/mwmr"
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
				log.Printf("Writer %d getphase from replica %d failed: %v", *cid, rid, err)
			} else {
				temp := Timestamp{
					Time: getReply.GetTime(),
					Cid:  getReply.GetCid(),
				}
				ch <- temp
			}
			wg.Done()
		}(i)
	}

	done := make([]Timestamp, 0, f+1)
	for len(done) < f+1 {
		select {
		case ts := <-ch:
			done = append(done, ts)
		default:

		}
	}

	wg.Wait()
	close(ch)

	return getNewTimestamp(done)
}

func writerSetPhase(key, value string, ts Timestamp) {
	n := 2*f + 1

	var wg sync.WaitGroup
	wg.Add(n)

	ch := make(chan bool, n)

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

			setReply, err := c.SetPhase(ctx, &pb.SetRequest{
				Key:   key,
				Value: value,
				Time:  ts.Time,
				Cid:   ts.Cid,
			})
			if err != nil {
				log.Printf("Writer %d setphase from replica %d failed: %v", *cid, rid, err)
			} else {
				ch <- setReply.GetApplied()
			}
			wg.Done()
		}(i)
	}

	done := 0
	for done < f+1 {
		select {
		case _ = <-ch:
			done++
		default:

		}
	}

	wg.Wait()
	close(ch)
}

func write(key, value string) {
	newTimestamp := writerGetPhase(key)
	writerSetPhase(key, value, newTimestamp)
	total_sets += 1
}
