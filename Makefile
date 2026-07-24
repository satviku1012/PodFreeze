.PHONY: deps proto build-proxy build-daemon clean

BIN_DIR := bin

all: build-proxy build-daemon

deps:
	go mod tidy

proto:
	mkdir -p $(BIN_DIR) api/proto/podfreeze
	protoc --go_out=. --go-grpc_out=. --go_opt=module=github.com/podfreeze/core --go-grpc_opt=module=github.com/podfreeze/core api/proto/podfreeze.proto

build-proxy: proto
	go build -o $(BIN_DIR)/proxy ./cmd/proxy

build-daemon: proto
	go build -o $(BIN_DIR)/daemon ./cmd/daemon

clean:
	rm -rf $(BIN_DIR)
