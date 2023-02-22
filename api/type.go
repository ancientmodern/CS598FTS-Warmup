package api

type Timestamp struct {
	ts       int
	clientID int
}

type Message struct {
	key   string
	value string
	ts    Timestamp
}
