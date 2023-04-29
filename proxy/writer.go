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

func (s *RegProxy) writerGetPhase(key uint64) Timestamp {
	var wg sync.WaitGroup
	wg.Add(n)

	ch := make(chan Timestamp, n)

	for i := 0; i < n; i++ {
		go func(rid int) {
			// Contact the replica and print out its response.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			getReply, err := s.grpcClients[rid].GetPhase(ctx, &pb.GetRequest{
				Key: key,
			})
			if err != nil {
				ErrorLogger.Printf("Writer %d getPhase from replica %d failed: %v", *cid, rid, err)
				temp := Timestamp{
					Time: 0,
					Cid:  getReply.GetCid(),
				}
				ch <- temp
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

func (s *RegProxy) writerSetPhase(key uint64, pair Pair) {
	var wg sync.WaitGroup
	wg.Add(n)

	ch := make(chan bool, n)

	for i := 0; i < n; i++ {
		go func(rid int) {
			// Contact the replica and print out its response.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			setReply, err := s.grpcClients[rid].SetPhase(ctx, &pb.SetRequest{
				Key:   key,
				Value: pair.Value,
				Time:  pair.Ts.Time,
				Cid:   pair.Ts.Cid,
			})
			if err != nil {
				ErrorLogger.Printf("Writer %d setPhase from replica %d failed: %v", *cid, rid, err)
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

func (s *RegProxy) write(key uint64, value uint32) {
	newTimestamp := s.writerGetPhase(key)
	s.writerSetPhase(key, Pair{
		Value: value,
		Ts:    newTimestamp,
	})
}
