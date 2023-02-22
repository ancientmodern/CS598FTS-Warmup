package common

type Timestamp struct {
	Time int64
	Cid  int64
}

type Pair struct {
	Value string
	Ts    Timestamp
}
