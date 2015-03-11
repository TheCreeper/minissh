package main

import (
	"flag"
	"log"
	"os"
	"sync"
)

var (

	// Waitgroup for all go routines
	wg sync.WaitGroup

	// Debug/Verbose switch
	Verbose bool

	// Generate sameple configuration switch
	GenConfig bool

	// Configuration filepath
	ConfigFile string
)

func init() {

	flag.BoolVar(&Verbose, "v", false, "debugging/verbose information")
	flag.BoolVar(&GenConfig, "g", false, "Generate a configuration file")
	flag.StringVar(&ConfigFile, "f", "sshd.conf", "The configuration to parse")
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

	for _, v := range cfg.ListenAddress {

		if Verbose {

			log.Printf("Starting SSH Server on %s", v)
		}

		wg.Add(1)
		go cfg.NewListenServer(v)
	}

	wg.Wait()
}
