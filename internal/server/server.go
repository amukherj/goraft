package server

import (
	"fmt"
	"log"

	"github.com/amukherj/goraft/internal/config"
	"github.com/amukherj/goraft/internal/raftlog"
	"github.com/amukherj/goraft/internal/terminfo"
	"github.com/amukherj/goraft/pkg/rafttypes"
)

type ServerMode int

const (
	_ = iota
	modeFollower
	modeCandidate
	modeLeader
)

type RemoteSyncState struct {
	matchIndex int32
	nextIndex  int32
}

type RaftServer struct {
	config           *config.RaftConfig
	mode             ServerMode
	currentTerm      int32
	votedFor         int32
	commitIndex      int32
	lastApplied      int32
	remoteSyncStates map[int32]RemoteSyncState
	log              []*rafttypes.LogEntry
}

func NewRaftServer(config *config.RaftConfig) (*RaftServer, error) {
	rs := &RaftServer{
		config: config,
		mode:   modeFollower,
	}

	// initialize term info
	if err := rs.initTermInfo(); err != nil {
		return nil, err
	}

	if err := rs.initLog(); err != nil {
		return nil, err
	}

	return rs, nil
}

func (rs *RaftServer) initTermInfo() error {
	termInfoPath := terminfo.NewTermInfoPath(rs.config.Config.TermInfoPath)
	termInfo, absent, err := termInfoPath.Create()
	if err != nil {
		return fmt.Errorf("Error creating / reading terminfo: %v", err)
	}
	if absent {
		log.Printf("Created new termInfo")
	}
	log.Printf("Current termInfo: %v, %v",
		termInfo.GetCurrentTerm(), termInfo.GetVotedFor())

	rs.currentTerm = termInfo.GetCurrentTerm()
	rs.votedFor = termInfo.GetVotedFor()
	return nil
}

func (rs *RaftServer) initLog() error {
	raftLogPath := raftlog.NewRaftLogPath(rs.config.Config.LogPath)
	entries, _, err := raftLogPath.Create()
	if err != nil {
		return fmt.Errorf("Error creating / reading logs: %v", err)
	}
	rs.log = entries
	rs.commitIndex = 0
	rs.lastApplied = 0
	return nil
}

func (rs *RaftServer) Run() error {
	return nil
}
