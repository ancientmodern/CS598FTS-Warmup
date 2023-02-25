package main

import (
	. "CS598FTS-Warmup/common"
	pb "CS598FTS-Warmup/mwmr"
	"context"
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
	n := 2*f + 1

	var wg sync.WaitGroup
	wg.Add(n)

	ch := make(chan Timestamp, n)

	for i := 0; i < n; i++ {
		go func(rid int) {
			// Contact the server and print out its response.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			getReply, err := grpcClient[rid].GetPhase(ctx, &pb.GetRequest{
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

	go func() {
		wg.Wait()
		close(ch)
	}()

	done := make([]Timestamp, 0, f+1)
	// for len(done) < f+1 {
	// 	select {
	// 	case ts := <-ch:
	// 		done = append(done, ts)
	// 	default:

	// 	}
	// }

	for pair := range ch {
		done = append(done, pair)
		if len(done) >= f+1 {
			break
		}
	}

	// wg.Wait()
	// close(ch)

	return getNewTimestamp(done)
}

func writerSetPhase(key, value string, ts Timestamp) {
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

func write(key, value string) {
	newTimestamp := writerGetPhase(key)
	writerSetPhase(key, value, newTimestamp)
	totalSets += 1
}
