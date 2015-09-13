package main

import (
	"flag"
	"log"
	"os"

	"github.com/TheCreeper/MiniSSH/minissh"
)

var (
	// Generate sameple configuration switch.
	GenConfig bool

	// Configuration filepath.
	ConfigFile string
)

func init() {
	flag.BoolVar(&Verbose, "v", false, "debugging/verbose information")
	flag.BoolVar(&GenConfig, "genconf", false, "Generate a configuration file")
	flag.StringVar(&ConfigFile, "f", "ssh.conf", "The configuration to parse")
	flag.Parse()

	if GenConfig {
		cfg, err := GenerateConfig()
		if err != nil {
			log.Fatal(err)
		}

		_, err = os.Stdout.Write(cfg)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
}

func main() {
	cfg, err := GetCFG(ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := minissh.NewServer(); err != nil {
		log.Fatal(err)
	}
}
