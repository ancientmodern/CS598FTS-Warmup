package main

import (
	. "CS598FTS-Warmup/common"
	pb "CS598FTS-Warmup/mwmr"
	"context"
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
	var wg sync.WaitGroup
	wg.Add(n)

	ch := make(chan Pair, n)

	for i := 0; i < n; i++ {
		go func(rid int) {
			// Contact the replica and print out its response.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			getReply, err := grpcClients[rid].GetPhase(ctx, &pb.GetRequest{
				Key: key,
			})
			if err != nil {
				ErrorLogger.Printf("Reader %d getPhase from replica %d failed: %v", *cid, rid, err)
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

	go func() {
		wg.Wait()
		close(ch)
	}()

	done := make([]Pair, 0, f+1)

	for pair := range ch {
		done = append(done, pair)
		if len(done) >= f+1 {
			break
		}
	}

	return getReadPair(done)
}

func readerSetPhase(key string, pair Pair) {
	var wg sync.WaitGroup
	wg.Add(n)

	ch := make(chan bool, n)

	for i := 0; i < n; i++ {
		go func(rid int) {
			// Contact the replica and print out its response.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			setReply, err := grpcClients[rid].SetPhase(ctx, &pb.SetRequest{
				Key:   key,
				Value: pair.Value,
				Time:  pair.Ts.Time,
				Cid:   pair.Ts.Cid,
			})
			if err != nil {
				ErrorLogger.Printf("Reader %d setPhase from replica %d failed: %v", *cid, rid, err)
			} else {
				ch <- setReply.GetApplied()
			}
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	done := 0

	for range ch {
		done++
		if done >= f+1 {
			break
		}
	}
}

func read(key string) string {
	readPair := readerGetPhase(key)
	readerSetPhase(key, readPair)
	return readPair.Value
}
