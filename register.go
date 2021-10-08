package main

import (
	"client/config"
	"fmt"
	"log"
)

func RegisterServer(initValue *InitValue) error {

}

func readServer(server *config.Server, initValue *InitValue) error {
	var host string
	var err error
	if len(initValue.ServerHost) == 0 {
		host, err = scanln("please enter dps server host")
		if err != nil {
			return err
		}
	} else {
		host = initValue.ServerHost
	}
	var serverPath string
	if len(initValue.ServerPath) == 0 {
		serverPath, err = scanln("please enter dps server path")
		if err != nil {
			return err
		}
	} else {
		serverPath = initValue.ServerHost
	}
	server.Address = host + serverPath

}

func scanln(tip string) (string, error) {
	log.Println(tip)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	log.Println("your input \"" + input + "\"")
	return input, nil
}
