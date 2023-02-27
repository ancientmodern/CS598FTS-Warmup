package main

import (
	. "CS598FTS-Warmup/common"
	pb "CS598FTS-Warmup/mwmr"
	"context"
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
				ErrorLogger.Printf("Writer %d getphase from replica %d failed: %v", *cid, rid, err)
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

	for pair := range ch {
		done = append(done, pair)
		if len(done) >= f+1 {
			break
		}
	}

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
				ErrorLogger.Printf("Writer %d setphase from replica %d failed: %v", *cid, rid, err)
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

	for p := range ch {
		if p && !p {
		}
		done++
		if done >= f+1 {
			break
		}
	}
}

func write(key, value string) (int64, int64, Timestamp) {
	t1 := time.Now().UnixNano()
	newTimestamp := writerGetPhase(key)
	t2 := time.Now().UnixNano()
	writerSetPhase(key, value, newTimestamp)
	t3 := time.Now().UnixNano()
	totalSets += 1
	return t2 - t1, t3 - t2, newTimestamp
}
