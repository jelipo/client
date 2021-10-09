package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var GlobalConfig *Config

//
type Config struct {
	Server Server `json:"server"`
	Local  Local  `json:"local"`
}

type Server struct {
	RunnerId  string `json:"runnerId"`
	Address   string `json:"address"`
	Token     string `json:"token"`
	MaxJobNum int    `json:"maxWorkerNum"`
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

func InitConfig(path string) error {
	reader := newConfigReader(path)
	config, err := reader.readNewConfig()
	if err != nil {
		return err
	}
	GlobalConfig = config
	return nil
}

func WriteConfigFile(config *Config, path string) error {
	configJsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(configJsonBytes)
	if err != nil {
		return err
	}
	return nil
}
