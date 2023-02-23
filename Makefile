all: build

OUTPUT_PATH=output
CLIENT_BINARY=$(OUTPUT_PATH)/client
REPLICA_BINARY=$(OUTPUT_PATH)/replica

CLIENT_DIR=./client
REPLICA_DIR=./server
PROTO_FILE=mwmr/mwmr.proto

build:
	go build -race -o $(CLIENT_BINARY) $(CLIENT_DIR)
	go build -race -o $(REPLICA_BINARY) $(REPLICA_DIR)

clean:
	go clean
	rm -f $(OUTPUT_PATH)/*

client:
	go run -race $(CLIENT_DIR)

replica:
	go run -race $(REPLICA_DIR)

generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $(PROTO_FILE)

.PHONY: all build clean client replica generate