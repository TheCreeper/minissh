package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os/exec"
	"sync"

	"code.google.com/p/go.crypto/ssh"
	"github.com/TheCreeper/MiniSSH/userx"
)

var Verbose bool

func (cfg *ClientConfig) NewListenServer(wg sync.WaitGroup, addr string) {

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{

		PublicKeyCallback: HandlePublicKeyCallback,
		AuthLogCallback:   HandleAuthLogCallback,
	}

	s, err := ReadHostKeys("/home/dc0/.ssh")
	if err != nil {

		log.Fatal(err)
	}
	for _, v := range s {

		config.AddHostKey(v)
	}

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	conn, err := net.Listen("tcp", addr)
	if err != nil {

		log.Fatal("Failed to listen for connection: ", err)
	}
	for {

		sConn, err := conn.Accept()
		if err != nil {

			log.Fatal(err)
		}

		_, chans, reqs, err := ssh.NewServerConn(sConn, config)
		if err != nil {

			log.Print("Failed to handshake: ", err)
		}

		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)

		// Handle the incomming request
		go handleServerConn(chans)
	}

	wg.Done()
}

func handleServerConn(chans <-chan ssh.NewChannel) {

	for newChannel := range chans {

		if newChannel.ChannelType() != "session" {

			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {

			return
		}

		go handleChannelRequest(channel, requests)
	}
}

func handleChannelRequest(channel ssh.Channel, in <-chan *ssh.Request) {

	for req := range in {

		switch req.Type {

		case "shell":

			go handleNewTerminal(channel)
		}
	}
}

func handleNewTerminal(channel ssh.Channel) {

	cmd := exec.Command("/usr/bin/zsh", "")
	inpipe, err := cmd.StdinPipe()
	if err != nil {

		return
	}

	outpipe, err := cmd.StdoutPipe()
	if err != nil {

		return
	}

	go CopyRequest(channel, inpipe)
	go CopyResponse(outpipe, channel)
}

func CopyRequest(in io.ReadCloser, out io.WriteCloser) {

	written, err := io.Copy(out, in)
	if Verbose {

		log.Printf("%d bytes written", written)
	}
	if err != nil {

		log.Printf("io.Copy() %s\n", err)
		return
	}

	return
}

func CopyResponse(out io.ReadCloser, in io.WriteCloser) {

	written, err := io.Copy(in, out)
	if Verbose {

		log.Printf("%d bytes written", written)
	}
	if err != nil {

		log.Printf("io.Copy() %s\n", err)
		return
	}

	return
}

func init() {

	flag.StringVar(&ConfigFile, "config", "./config.json", "The configuration file location")
	flag.BoolVar(&Verbose, "v", false, "Be more verbose")
	flag.Parse()
}

func main() {

	var wg sync.WaitGroup

	cfg, err := GetCFG(ConfigFile)
	if err != nil {

		log.Fatal(err)
	}

	for _, v := range cfg.Servers {

		log.Printf("Starting Server on %s", v.Addr)

		wg.Add(1)
		go cfg.NewListenServer(wg, v.Addr)
	}

	wg.Wait()
}
