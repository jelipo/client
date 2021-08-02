package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var GlobalConfig *Config

//
type Config struct {
	Server Server `json:"server"`
	Local  Local  `json:"local"`
}

type Server struct {
	ClientId     string `json:"clientId"`
	Address      string `json:"address"`
	Token        string `json:"token"`
	MaxWorkerNum int    `json:"maxWorkerNum"`
}

type Local struct {
	ClientWorkDir string `json:"client_work_dir"`
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

func IninConfig(path string) error {
	reader := newConfigReader(path)
	config, err := reader.readNewConfig()
	if err != nil {
		return err
	}
	GlobalConfig = config
	return nil
}
