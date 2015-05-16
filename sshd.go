package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"code.google.com/p/go.crypto/ssh"
	//"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/TheCreeper/MiniSSH/sshutil"
	"github.com/TheCreeper/MiniSSH/userx"
	"github.com/creack/termios/win"
	"github.com/kr/pty"
)

func (cfg *ServerConfig) NewListenServer(addr string) {

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{

		PublicKeyCallback: HandlePublicKeyCallback,
		AuthLogCallback:   HandleAuthLogCallback,
	}

	// Parse the host keys
	for _, v := range cfg.HostKeys {

		key, err := sshutil.ReadHostKey(v)
		if err != nil {

			log.Printf("ReadHostKey: %s, %s", err, cfg.HostKeys)
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
		go HandleServerConn(sshconn, chans)
	}

	wg.Done()
}

func HandlePublicKeyCallback(conn ssh.ConnMetadata, key ssh.PublicKey) (perm *ssh.Permissions, err error) {

	publicKeys, err := userx.ReadUserAuthKeys(conn.User())
	if err != nil {

		return
	}
	for _, v := range publicKeys {

		certchecker := &ssh.CertChecker{}

		p, err := certchecker.Authenticate(conn, v)
		if err != nil {

			log.Print(err)
			continue
		}

		return p, nil
	}

	return
}

// Log all authentication attempts
func HandleAuthLogCallback(conn ssh.ConnMetadata, method string, err error) {

	log.Printf("Login attempt for user %s from %s with client %s",
		conn.User(),
		conn.RemoteAddr(),
		string(conn.ClientVersion()))
}

func HandleServerConn(conn ssh.ConnMetadata, chans <-chan ssh.NewChannel) {

	for newChannel := range chans {

		if t := newChannel.ChannelType(); t != "session" {

			newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {

			log.Printf("Could not accept channel (%s)", err)
			continue
		}

		s := &Session{

			Conn: conn,
		}
		go s.HandleChannelRequest(channel, requests)
	}
}

type Session struct {

	// SSH Conn Meta Data
	Conn ssh.ConnMetadata

	// pty
	pty *os.File

	// Win Size
	THeight uint16
	TWidth  uint16
}

func (s *Session) HandleChannelRequest(channel ssh.Channel, in <-chan *ssh.Request) {

	for req := range in {

		switch req.Type {

		case "shell":

			if len(req.Payload) == 0 {

				if err := req.Reply(true, nil); err != nil {

					log.Printf("Unable to reply to channel request (%s)", err)
					continue
				}
				go s.HandleNewTerminal(channel)
			}

		case "pty-req":

			termLen := req.Payload[3]
			w, h := ParseTerminalDims(req.Payload[termLen+4:])
			s.THeight = uint16(h)
			s.TWidth = uint16(w)

			if err := req.Reply(true, nil); err != nil {

				log.Printf("Unable to reply to channel request (%s)", err)
			}

		case "window-change":

			w := &win.Winsize{

				Height: uint16(binary.BigEndian.Uint32(req.Payload[4:])),
				Width:  uint16(binary.BigEndian.Uint32(req.Payload)),
			}
			if err := win.SetWinsize(s.pty.Fd(), w); err != nil {

				log.Printf("Unable to set window-change (%s)", err)
			}

		case "close":

			if err := channel.Close(); err != nil {

				log.Printf("Unable to close channel (%s)", err)
			}

		case "default":

			log.Printf("Invalid Request Type (%s)", req.Type)

		}
	}
}

func (s *Session) HandleNewTerminal(channel ssh.Channel) {

	cmd := exec.Command("/bin/bash")
	cmd.Dir = "/home/dc0"

	f, err := pty.Start(cmd)
	if err != nil {

		log.Printf("Could not start pty (%s)", err)
	}
	s.pty = f

	w := &win.Winsize{

		Height: s.THeight,
		Width:  s.TWidth,
	}
	if err := win.SetWinsize(s.pty.Fd(), w); err != nil {

		log.Printf("Unable to set window-change (%s)", err)
	}

	go CopyRequest(channel, s.pty)
	go CopyResponse(s.pty, channel)
}
