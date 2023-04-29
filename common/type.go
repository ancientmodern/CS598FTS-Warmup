package common

import (
	"errors"
	"sync"
)

var ErrKeyNotFound = errors.New("key not found")

type Timestamp struct {
	Time int64
	Cid  int64
}

type Pair struct {
	Value uint32
	Ts    Timestamp
}

type PairMutex struct {
	Value uint32
	Ts    Timestamp
	Mtx   sync.RWMutex
}

func LessTimestamp(lhs, rhs Timestamp) bool {
	return lhs.Time < rhs.Time || (lhs.Time == rhs.Time && lhs.Cid < rhs.Cid)
}
