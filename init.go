package main

import (
	"client/util"
	_ "embed"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
	"os"
)
import nested "github.com/antonfisher/nested-logrus-formatter"

func InitContext() (*InitValue, error) {
	initLogContext()
	printBanner()
	value, err := readEnv()
	if err != nil {
		return nil, err
	}
	if len(value.ConfigPath) == 0 {
		currentPath := util.GetCurrentPath()
		value.ConfigPath = currentPath + "/dps_runner/config.json"
	}
	return value, nil
}

//go:embed banner.txt
var bannerStr string

func printBanner() {
	fmt.Println(bannerStr)
}

func initLogContext() {
	logrus.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		HideKeys:        true,
	})
}

func ConfigExist(configPath string) (bool, error) {
	file, err := os.Stat(configPath)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if file.IsDir() {
		return false, errors.New("config is a dir")
	}
	return true, nil
}

func readEnv() (*InitValue, error) {
	initValue := InitValue{}
	if err := env.Parse(&initValue); err != nil {
		return nil, errors.New("read env failed")
	}
	return &initValue, nil
}

type InitValue struct {
	ConfigPath   string `env:"CONFIG_PATH"`
	ServerHost   string `env:"SERVER_HOST"`
	ServerPath   string `env:"SERVER_PATH" envDefault:"/cicdengine/api/v1/runner"`
	OneTimeToken string `env:"ONE_TIME_TOKEN"`
	WorkDir      string `env:"WORK_DIR"`
}
