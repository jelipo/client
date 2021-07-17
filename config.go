package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//
type Config struct {
	Server Server `json:"server"`
}

type Server struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

type ConfigReader struct {
	configPath string
}

func newConfigReader(path string) ConfigReader {
	return ConfigReader{configPath: path}
}

func (reader *ConfigReader) readNewConfig() (*Config, error) {
	open, err := os.Open(reader.configPath)
	if err != nil {
		return nil, err
	}
	all, err := ioutil.ReadAll(open)
	var config Config
	jsonErr := json.Unmarshal(all, &config)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return &config, nil
}
