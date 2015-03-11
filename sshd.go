package main

import (
	"log"
	"net"

	"code.google.com/p/go.crypto/ssh"
	"github.com/TheCreeper/MiniSSH/userx"
)

func (cfg *ServerConfig) NewListenServer(addr string) {

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{

		PublicKeyCallback: handlePublicKeyCallback,
		AuthLogCallback:   handleAuthLogCallback,
	}

	// Parse the host keys
	for _, v := range cfg.HostKeys {

		key, err := userx.ReadHostKey(v)
		if err != nil {

			log.Fatalf("ReadHostKeys: %s, %s", err, cfg.HostKeys)
		}

		config.AddHostKey(key)
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

			log.Printf("conn.Accept: %s", err)
			continue
		}

		sshconn, chans, reqs, err := ssh.NewServerConn(sConn, config)
		if err != nil {

			log.Printf("ssh.NewServerConn: Failed to handshake: %s", err)
			continue
		}

		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)

		// Handle the incomming request
		go handleServerConn(sshconn, chans)
	}

	wg.Done()
}
