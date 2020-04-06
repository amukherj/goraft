.PHONY: protos
protos:
	protoc -I protos/raft --go_out=. persistent.proto

.PHONY: bin/server
bin/server:
	go build -o bin/server github.com/amukherj/goraft/cmd/server
