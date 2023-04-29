all: build

OUTPUT_PATH=output
PROXY_BINARY=$(OUTPUT_PATH)/proxy
REPLICA_BINARY=$(OUTPUT_PATH)/replica

PROXY_DIR=./proxy
REPLICA_DIR=./replica
PROTO_FILE=mwmr/mwmr.proto

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(PROXY_BINARY) $(PROXY_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(REPLICA_BINARY) $(REPLICA_DIR)

clean:
	go clean
	rm -f $(OUTPUT_PATH)/*

client:
	go run -race $(PROXY_DIR)

replica:
	go run -race $(REPLICA_DIR)

generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $(PROTO_FILE)

.PHONY: all build clean client replica generate