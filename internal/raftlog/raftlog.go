package raftlog

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/amukherj/goraft/pkg/rafttypes"
	"github.com/golang/protobuf/proto"
)

type RaftLogPath struct {
	path string
}

func NewRaftLogPath(path string) *RaftLogPath {
	return &RaftLogPath{
		path: path,
	}
}

// Creates if not present, and populates with initial values.
// Returns:
//   (TermInfo, true, nil): if not present.
//   (TermInfo, true, err): if not present and creation failed.
//   (TermInfo, false, nil): if already present.
//   (TermInfo, false, err): if already present and read failed.
//
func (rlp *RaftLogPath) Create() ([]*rafttypes.LogEntry, bool, error) {
	rlf, err := os.OpenFile(rlp.path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		// why send the second return arg true: assume the file isn't
		// there or isn't accessible.
		return nil, true, err
	}
	defer func() {
		if err := rlf.Close(); err != nil {
			log.Printf("Closing open file %s failed: %v", rlp.path, err)
		}
	}()

	content, err := ioutil.ReadAll(rlf)
	if err != nil {
		// why send the second return arg false: assume the file is
		// there and is accessible.
		return nil, false, err
	}
	var absent bool = false
	if len(content) == 0 {
		absent = true
		return []*rafttypes.LogEntry{}, absent, nil
	}
	var result []*rafttypes.LogEntry
	lrs := NewLogRecordStream(content)
	for {
		record, err1 := lrs.unmarshalNextRecord()
		if record == nil || err1 != nil {
			err = err1
			break
		}
		result = append(result, record)
	}
	return result, absent, err
}

// Append a series of log records to the existing log path.
func (rlp *RaftLogPath) Append(payload []*rafttypes.LogEntry) error {
	if payload == nil {
		return fmt.Errorf("Null payload passed")
	}

	rlf, err := os.OpenFile(rlp.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		// why send the second return arg true: assume the file isn't
		// there or isn't accessible.
		return err
	}
	defer func() {
		if err := rlf.Close(); err != nil {
			log.Printf("Closing open file %s failed: %v", rlp.path, err)
		}
	}()

	var buflist net.Buffers = make([][]byte, 0, len(payload))
	contentLen := 0
	for i, entry := range payload {
		buf, err := marshalRecord(entry)
		if err != nil {
			return fmt.Errorf("Error marshaling record at index %d: %v", i, err)
		}
		contentLen += len(buf)
		buflist = append(buflist, buf)
	}
	if written, err := buflist.WriteTo(rlf); err != nil {
		return fmt.Errorf("Failed to append log entries: %v", err)
	} else if written != int64(contentLen) {
		return fmt.Errorf("Log records appended are suspect")
	}

	return nil
}

// Type that encapsulates a stream of length-prefix log entries.
// Not concurrency safe because of the currentOffset member which maintains
// the offset at which the next read will happen
type LogRecordStream struct {
	buf           []byte
	length        int
	currentOffset int
	valid         bool
}

func NewLogRecordStream(buf []byte) *LogRecordStream {
	return &LogRecordStream{
		buf:           buf,
		length:        len(buf),
		currentOffset: 0,
		valid:         true,
	}
}

// Reads and returns the next record, plus the offset to start reading the
// record following the one returned.
//
func (lrs *LogRecordStream) unmarshalNextRecord() (*rafttypes.LogEntry, error) {
	if lrs.currentOffset == lrs.length {
		return nil, nil
	}
	if !lrs.valid {
		return nil, fmt.Errorf("Invalid record stream state")
	}
	prefixSize := 4
	sizeBuf := lrs.buf[lrs.currentOffset : lrs.currentOffset+prefixSize]
	recSize := int(binary.LittleEndian.Uint32(sizeBuf))
	recStart := lrs.currentOffset + prefixSize
	recBuf := lrs.buf[recStart : recStart+recSize]

	var logEntry rafttypes.LogEntry
	var err error
	if err = proto.Unmarshal(recBuf, &logEntry); err != nil {
		lrs.valid = false
		return nil, err
	}
	lrs.currentOffset += prefixSize + recSize
	return &logEntry, err
}

func marshalRecord(entry *rafttypes.LogEntry) ([]byte, error) {
	buf, err := proto.Marshal(entry)
	if err != nil {
		return nil, err
	}
	prefixLen := 4
	recSize := len(buf)
	result := make([]byte, prefixLen, prefixLen+recSize)
	binary.LittleEndian.PutUint32(result, uint32(recSize))
	result = append(result, buf...)
	return result, nil
}
