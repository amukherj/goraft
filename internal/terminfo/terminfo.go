package terminfo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/amukherj/goraft/pkg/rafttypes"
	"github.com/golang/protobuf/proto"
)

type TermInfoPath struct {
	path string
}

func NewTermInfoPath(path string) *TermInfoPath {
	return &TermInfoPath{
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
func (tip *TermInfoPath) Create() (*rafttypes.TermInfo, bool, error) {
	tif, err := os.OpenFile(tip.path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		// why send the second return arg true: assume the file isn't
		// there or isn't accessible.
		return nil, true, err
	}
	defer func() {
		if err := tif.Close(); err != nil {
			log.Printf("Closing open file %s failed: %v", tip.path, err)
		}
	}()

	content, err := ioutil.ReadAll(tif)
	if err != nil {
		// why send the second return arg false: assume the file is
		// there and is accessible.
		return nil, false, err
	}
	var termInfo rafttypes.TermInfo
	var absent bool = false
	if len(content) == 0 {
		absent = true
		termInfo.CurrentTerm = 0
		termInfo.VotedFor = -1
		content, err = proto.Marshal(&termInfo)
		if err != nil {
			// why send the second return arg false: assume the file is
			// there and is accessible.
			return nil, absent, fmt.Errorf("Failed to marshal terminfo: %v", err)
		}
		if written, err := tif.Write(content); err != nil {
			return nil, absent, fmt.Errorf("Failed to initialize TermInfo: %v", err)
		} else if written != len(content) {
			return nil, absent, fmt.Errorf("TermInfo content written is suspect")
		}
	} else {
		absent = false
		err = proto.Unmarshal(content, &termInfo)
		if err != nil {
			// why send the second return arg false: assume the file is
			// there and is accessible.
			return nil, absent, fmt.Errorf("Failed to read TermInfo: %v", err)
		}
	}
	return &termInfo, absent, nil
}

func (tip *TermInfoPath) Update(payload *rafttypes.TermInfo) error {
	if payload == nil {
		return fmt.Errorf("Null payload passed")
	}

	tif, err := os.OpenFile(tip.path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		// why send the second return arg true: assume the file isn't
		// there or isn't accessible.
		return err
	}
	defer func() {
		if err := tif.Close(); err != nil {
			log.Printf("Closing open file %s failed: %v", tip.path, err)
		}
	}()

	content, err := proto.Marshal(payload)
	if err != nil {
		// why send the second return arg false: assume the file is
		// there and is accessible.
		return fmt.Errorf("Failed to marshal terminfo: %v", err)
	}
	if written, err := tif.Write(content); err != nil {
		return fmt.Errorf("Failed to initialize TermInfo: %v", err)
	} else if written != len(content) {
		return fmt.Errorf("TermInfo content written is suspect")
	}

	return nil
}
