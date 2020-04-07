package state

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"github.com/amukherj/goraft/pkg/rafttypes"
	"github.com/gogo/protobuf/proto"
)

type RaftInfo interface {
	GetTermInfo() (*rafttypes.TermInfo, error)
	UpdateTermInfo(info rafttypes.TermInfo) error
	GetAllLogs() ([]*rafttypes.LogEntry, error)
	AppendLog(entry *rafttypes.LogEntry) error
}

type RaftInfoReader interface {
	GetTermInfoReader() (io.Reader, error)
	GetLogReader() (io.Reader, error)
	GetLogWriter() (io.Writer, error)
}

type raftInfoFromFile struct {
	termInfoPath string
	logPath      string
}

func NewRaftInfoFromFile(termInfoPath string, logPath string) RaftInfo {
	return &raftInfoFromFile{
		termInfoPath: termInfoPath,
		logPath:      logPath,
	}
}

func (riff *raftInfoFromFile) GetTermInfoReader() (io.Reader, error) {
	data, err := ioutil.ReadFile(riff.termInfoPath)
	if err != nil {
		return nil, fmt.Errorf("Could not read term_info file %s: %v", riff.termInfoPath, err)
	}
	return bytes.NewReader(data), nil
}

func (riff *raftInfoFromFile) GetTermInfo() (*rafttypes.TermInfo, error) {
	reader, err := riff.GetTermInfoReader()
	if err != nil {
		return nil, err
	}
	var result rafttypes.TermInfo
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Could not read term_info data: %v", err)
	}
	if len(data) == 0 {
		result.CurrentTerm = 0
		result.VotedFor = 0
		return &result, nil
	}
	if err = proto.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("Could not unmarshal content in file %s: %v", riff.termInfoPath, err)
	}
	return &result, nil
}

func (riff *raftInfoFromFile) UpdateTermInfo(info rafttypes.TermInfo) error {
	data, err := proto.Marshal(&info)
	if err != nil {
		return fmt.Errorf("Could not marshal TermInfo: %v", err)
	}
	if err = ioutil.WriteFile(riff.termInfoPath, data, os.FileMode(int(0640))); err != nil {
		return fmt.Errorf("Could not write marshaled TermInfo to file %s: %v",
			riff.termInfoPath, err)
	}
	return nil
}

func (riff *raftInfoFromFile) GetLogReader() (io.Reader, error) {
	return os.Open(riff.logPath)
}

func (riff *raftInfoFromFile) GetLogWriter() (io.WriteCloser, error) {
	return os.Open(riff.logPath)
}

func (riff *raftInfoFromFile) GetAllLogs() ([]*rafttypes.LogEntry, error) {
	reader, err := riff.GetLogReader()
	if err != nil {
		return nil, fmt.Errorf("Could not open the log file for reading: %v", err)
	}

	result := []*rafttypes.LogEntry{}
	entry, err := riff.GetNextLogEntry(reader)
	for ; err == nil; entry, err = riff.GetNextLogEntry(reader) {
		result = append(result, entry)
	}
	return result, nil
}

func (riff *raftInfoFromFile) GetNextLogEntry(r io.Reader) (*rafttypes.LogEntry, error) {
	lenbuf := make([]byte, 4)
	_, err := io.ReadFull(r, lenbuf)
	if err != nil {
		return nil, fmt.Errorf("Failed to read record length: %v", err)
	}
	reclen := binary.LittleEndian.Uint32(lenbuf)
	msg := make([]byte, reclen)
	if _, err = io.ReadFull(r, msg); err != nil {
		return nil, fmt.Errorf("Failed to read record: %v", err)
	}
	var entry rafttypes.LogEntry
	if err = proto.Unmarshal(msg, &entry); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal record: %v", err)
	}
	return &entry, nil
}

func (riff *raftInfoFromFile) AppendLog(entry *rafttypes.LogEntry) error {
	data, err := proto.Marshal(entry)
	if err != nil {
		return fmt.Errorf("Failed marshal entry: %v", err)
	}
	reclen := len(data)
	lenbuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenbuf, uint32(reclen))

	w, err := riff.GetLogWriter()
	defer w.Close()

	if err != nil {
		return fmt.Errorf("Failed to get a log writer: %v", err)
	}

	bufs := net.Buffers{lenbuf, data}
	if _, err := bufs.WriteTo(w); err != nil {
		return fmt.Errorf("Failed to write record to log: %v", err)
	}
	return nil
}
