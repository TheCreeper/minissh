package main

import (
	"encoding/json"
	"io/ioutil"
)

type ServerConfig struct {

	// Address and port to listen for connections on.
	Addr        string
	HostKeysDir string
}

func GenerateConfig() ([]byte, error) {
	cfg := &ServerConfig{
		Addr:        ":22",
		HostKeysDir: "",
	}
	return json.MarshalIndent(cfg, "", "	")
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
