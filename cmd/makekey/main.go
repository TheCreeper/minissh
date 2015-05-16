package main

import (
	"flag"
	"log"

	"code.google.com/p/go.crypto/ssh"
)

var (

	// Debug/Verbose switch
	Verbose bool

	// Configuration filepath
	CreateKeyDir bool

	// Directory to place the generated keys
	KeyDir string
)

func init() {

	flag.BoolVar(&Verbose, "v", false, "debugging/verbose information")
	flag.BoolVar(&CreateKeyDir, "r", false, "Create key directory if it does not exist")
	flag.StringVar(&KeyDir, "d", "keydir", "Directory to place the generated keys")
	flag.Parse()
}

func main() {

}
