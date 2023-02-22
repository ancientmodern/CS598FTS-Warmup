all: build

OUTPUT_PATH=output
CLIENT_BINARY=$(OUTPUT_PATH)/client
REPLICA_BINARY=$(OUTPUT_PATH)/replica

CLIENT_FILE=client/main.go
REPLICA_FILE=server/main.go
PROTO_FILE=mwmr/mwmr.proto

build:
	go build -o $(CLIENT_BINARY) $(CLIENT_FILE)
	go build -o $(REPLICA_BINARY) $(REPLICA_FILE)

clean:
	go clean
	rm -f $(OUTPUT_PATH)/*

client:
	go run $(CLIENT_FILE)

replica:
	go run $(REPLICA_FILE)

generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $(PROTO_FILE)

.PHONY: all build clean client replica generate