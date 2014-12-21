package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
)

var ConfigFile string

type ClientConfig struct {
	Servers []struct {
		Addr string
	}
}

func (cfg *ClientConfig) Flush() (err error) {

	file, err := os.OpenFile(ConfigFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {

		return
	}
	defer file.Close()

	b, err := json.MarshalIndent(cfg, "", "	")
	if err != nil {

		return
	}

	buf := bufio.NewWriter(file)
	defer buf.Flush()

	_, err = buf.Write(b)
	if err != nil {

		return
	}

	return
}

func (cfg *ClientConfig) Validate() (err error) {

	return
}

func GetCFG(f string) (cfg *ClientConfig, err error) {

	b, err := ioutil.ReadFile(f)
	if err != nil {

		return
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {

		return
	}

	err = cfg.Validate()
	if err != nil {

		return
	}

	return
}
