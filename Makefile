all: bin/server

.PHONY: protos
protos:
	protoc -I protos/raft --go_out=. term_info.proto
	protoc -I protos/raft --go_out=plugins=grpc:. rpc.proto

.PHONY: bin/server
bin/server: protos
	go build -o bin/server github.com/amukherj/goraft/cmd/server
