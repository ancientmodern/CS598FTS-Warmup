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

func (s *RegProxy) readerGetPhase(key uint64) (Pair, error) {
	var wg sync.WaitGroup
	wg.Add(n)

	ch := make(chan Pair, n)
	errCh := make(chan error, n)

	for i := 0; i < n; i++ {
		go func(rid int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			getReply, err := s.grpcClients[rid].GetPhase(ctx, &pb.GetRequest{
				Key: key,
			})
			if err != nil {
				// ErrorLogger.Printf("Reader %d getPhase from replica %d failed: %v", *cid, rid, err)
				errCh <- err
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
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
		close(errCh)
	}()

	done := make([]Pair, 0, f+1)

	for {
		select {
		case pair, ok := <-ch:
			if ok {
				done = append(done, pair)
				if len(done) >= f+1 {
					return getReadPair(done), nil
				}
			}
		case err, ok := <-errCh:
			if ok && err == ErrKeyNotFound {
				return Pair{}, ErrKeyNotFound
			}
		default:
			if len(done) == 0 && len(ch) == 0 && len(errCh) == 0 {
				return Pair{}, nil
			}
		}
	}
}

func (s *RegProxy) readerSetPhase(key uint64, pair Pair) {
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
				// ErrorLogger.Printf("Reader %d setPhase from replica %d failed: %v", *cid, rid, err)
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

func (s *RegProxy) read(key uint64) uint32 {
	readPair, err := s.readerGetPhase(key)
	if err != nil {
		return 0xFF // 0xFF means key does not exist
	}
	s.readerSetPhase(key, readPair)
	return readPair.Value
}
