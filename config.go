package main

import (
	"encoding/json"
	"io/ioutil"
)

type ServerConfig struct {

	// Array of addresses to listen on
	ListenAddress []string

	// Array of hostkeys that can be parsed
	HostKeys []string
}

func GenerateConfig() ([]byte, error) {

	cfg := &ServerConfig{

		ListenAddress: []string{":22"},
		HostKeys:      []string{"/etc/ssh/ssh_host_rsa_key", "/etc/ssh/ssh_host_ecdsa_key"},
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
