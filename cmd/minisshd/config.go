package main

import (
	"encoding/json"
	"io/ioutil"
)

type ServerConfig struct {

	// Address and port to listen for connections on.
	Addr string
}

func GenerateConfig() ([]byte, error) {
	cfg := &ServerConfig{
		Addr: ":22",
	}
	return json.MarshalIndent(cfg, "", "	")
}

func (cfg *ServerConfig) Validate() (err error) {
	if len(cfg.Addr) == 0 {
		cfg.Addr = ":22"
	}
	return
}

func GetCFG(f string) (cfg *ServerConfig, err error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return
	}
	return
}
