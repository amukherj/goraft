This is an implementation of the [Raft Consensus Protocol](https://raft.github.io/raft.pdf) in Go.

## Code

## Server
The main server code with the `main` function is present in `cmd/server/server.go`.

### Config
1. The configuration structure is defined under `internal/config/config.go`.
2. The interface and code for reading config from a config store (such as a file) is in `internal/config/io.go`.


## Building
To build, run:

    make
