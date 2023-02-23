package common

type Timestamp struct {
	Time int64
	Cid  int64
}

type Pair struct {
	Value string
	Ts    Timestamp
}

func LessTimestamp(lhs, rhs Timestamp) bool {
	return lhs.Time < rhs.Time || (lhs.Time == rhs.Time && lhs.Cid < rhs.Cid)
}
