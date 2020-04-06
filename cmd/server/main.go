package main

import (
	"log"
	"os"
	"path"

	raftconfig "github.com/amukherj/goraft/internal/config"
	"github.com/amukherj/goraft/internal/server"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicf("Could not determine working directory: %v", err)
	}

	yamlConfig := path.Join(cwd, "raft.yaml")
	configReader := raftconfig.NewFileConfigReader(yamlConfig)
	config, err := configReader.Read()
	if err != nil {
		log.Printf("Could not read raft config from %s: %v", yamlConfig, err)
		jsonConfig := path.Join(cwd, "raft.json")
		configReader = raftconfig.NewFileConfigReader(yamlConfig)
		config, err = configReader.Read()
		if err != nil {
			log.Panicf("Could not read raft config from %s: %v", jsonConfig, err)
		}
	}
	if err = config.Validate(); err != nil {
		log.Panicf("Invalid configuration: %v", err)
	}

	raftServer, err := server.NewRaftServer(config)
	if err != nil {
		log.Panicf("Failed to initialize server: %v", err)
	}

	err = raftServer.Run()
	if err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
}
