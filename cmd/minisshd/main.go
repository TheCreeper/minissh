package main

import (
	"flag"
	"log"
	"os"
)

var (
	// Debug/Verbose switch
	Verbose bool

	// Log to syslog switch.
	Syslog bool

	// Generate sameple configuration switch.
	GenConfig bool

	// Configuration filepath.
	ConfigFile string
)

func init() {

	flag.BoolVar(&Verbose, "v", false, "debugging/verbose information")
	flag.BoolVar(&Syslog, "s", true, "log to syslog.")
	flag.BoolVar(&GenConfig, "g", false, "Generate a configuration file")
	flag.StringVar(&ConfigFile, "f", "/etc/minissh/ssh.conf", "The configuration to parse")
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

	if err := cfg.NewServer(); err != nil {

		log.Fatal(err)
	}
}
