package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigReader interface {
	Read() (*RaftConfig, error)
}

type fileConfigReader struct {
	fileName string
}

func NewFileConfigReader(fileName string) *fileConfigReader {
	return &fileConfigReader{
		fileName: fileName,
	}
}

func (fcr *fileConfigReader) Read() (*RaftConfig, error) {
	if _, err := os.Stat(fcr.fileName); err != nil {
		return nil, fmt.Errorf("Could not stat config file %s: %v", fcr.fileName, err)
	}
	data, err := ioutil.ReadFile(fcr.fileName)
	if err != nil {
		return nil, fmt.Errorf("Could not read config file %s: %v", fcr.fileName, err)
	}
	config, err := fcr.readYAML(data)
	if err != nil {
		config, err = fcr.readJSON(data)
	}

	return config, err
}

func (fcr *fileConfigReader) readYAML(data []byte) (*RaftConfig, error) {
	var config RaftConfig
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (fcr *fileConfigReader) readJSON(data []byte) (*RaftConfig, error) {
	var config RaftConfig
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
