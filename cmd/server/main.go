package main

import (
	"fmt"
	"log"
	"os"
	"path"

	raftconfig "github.com/amukherj/goraft/internal/config"
	raftstate "github.com/amukherj/goraft/internal/state"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicf("Could not determine working directory: %v", err)
	}

	// Read the config
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

	// Update data from previous runs if any
	if config.Config.LogPath == "" {
		log.Panicf("No log_path set in the config: %v", err)
	}
	if config.Config.TermInfoPath == "" {
		log.Panicf("No term_info_path set in the config: %v", err)
	}

	fmt.Printf("Config read: %+v", config)
	raftInfo := raftstate.NewRaftInfoFromFile(config.Config.TermInfoPath,
		config.Config.LogPath)
	termInfo, err := raftInfo.GetTermInfo()
	if err != nil {
		log.Panicf("Could not read termInfo: %v", err)
	}
	if termInfo.GetCurrentTerm() == 0 {
		log.Printf("Server starting up. Current term is 0")
	} else {
		log.Printf("Server starting up. Current term is %d",
			termInfo.GetCurrentTerm())
	}
}
