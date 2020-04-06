.PHONY: protos
protos:
	protoc -I protos/raft --go_out=. term_info.proto

.PHONY: bin/server
bin/server:
	go build -o bin/server github.com/amukherj/goraft/cmd/server
