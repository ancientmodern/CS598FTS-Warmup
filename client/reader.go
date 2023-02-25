package main

import (
	. "CS598FTS-Warmup/common"
	pb "CS598FTS-Warmup/mwmr"
	"context"
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
			// Contact the server and print out its response.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			getReply, err := grpcClient[rid].GetPhase(ctx, &pb.GetRequest{
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

	go func() {
		wg.Wait()
		close(ch)
	}()

	done := make([]Pair, 0, f+1)
	// for len(done) < f+1 {
	// 	log.Printf("here")
	// 	select {
	// 	case pair := <-ch:
	// 		done = append(done, pair)
	// 	default:
	// 	}
	// }

	for pair := range ch {
		done = append(done, pair)
		if len(done) >= f+1 {
			break
		}
	}

	return getReadPair(done)
}

func readerSetPhase(key string, pair Pair) {
	n := 2*f + 1

	var wg sync.WaitGroup
	wg.Add(n)

	ch := make(chan bool, n)

	for i := 0; i < n; i++ {
		go func(rid int) {
			// Contact the server and print out its response.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			setReply, err := grpcClient[rid].SetPhase(ctx, &pb.SetRequest{
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

	go func() {
		wg.Wait()
		close(ch)
	}()

	done := 0
	// for done < f+1 {
	// 	select {
	// 	case _ = <-ch:
	// 		done++
	// 	default:

	// 	}
	// }

	for p := range ch {
		if p && !p {
		}
		done++
		if done >= f+1 {
			break
		}
	}

	// wg.Wait()
	// close(ch)
}

func read(key string) string {
	readPair := readerGetPhase(key)
	readerSetPhase(key, readPair)
	totalGets += 1
	return readPair.Value
}
