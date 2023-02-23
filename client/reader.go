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

func getReadPair(arr []Pair) Pair {
	res := arr[0]
	for _, pair := range arr {
		if LessTimestamp(res.Ts, pair.Ts) {
			res = pair
		}
	}
	return res
}

func readerGetPhase(key string) Pair {
	n := 2*f + 1

	var wg sync.WaitGroup
	wg.Add(n)

	ch := make(chan Pair, n)

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
				log.Printf("Reader %d getphase from replica %d failed: %v", *cid, rid, err)
			} else {
				temp := Pair{
					Value: getReply.GetValue(),
					Ts: Timestamp{
						Time: getReply.GetTime(),
						Cid:  getReply.GetCid(),
					},
				}
				ch <- temp
			}
			wg.Done()
		}(i)
	}

	done := make([]Pair, 0, f+1)
	for len(done) < f+1 {
		select {
		case pair := <-ch:
			done = append(done, pair)
		default:

		}
	}

	close(ch)
	wg.Wait()

	return getReadPair(done)
}

func readerSetPhase(key string, pair Pair) {
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
				Value: pair.Value,
				Time:  pair.Ts.Time,
				Cid:   pair.Ts.Cid,
			})
			if err != nil {
				log.Printf("Reader %d setphase from replica %d failed: %v", *cid, rid, err)
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

	close(ch)
	wg.Wait()
}

func read(key string) string {
	readPair := readerGetPhase(key)
	log.Printf("Reader getPhase done, pair: {value: %s, time: %d, cid: %d}\n", readPair.Value, readPair.Ts.Time, readPair.Ts.Cid)
	readerSetPhase(key, readPair)
	log.Printf("Reader setPhase done, value: %s\n", readPair.Value)
	log.Println("====================================================================================================")
	return readPair.Value
}
