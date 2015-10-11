package main

import (
	"log"

	"github.com/TheCreeper/MiniSSH/minissh"
)

func main() {
	if err := minissh.NewServer(); err != nil {
		log.Fatal(err)
	}
}
