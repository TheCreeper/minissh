package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"code.google.com/p/go.crypto/ssh"
	//"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/TheCreeper/MiniSSH/userx"
	"github.com/creack/termios/win"
	"github.com/kr/pty"
)

// Some errors
var (
	ErrKeyMatch = errors.New("Authorised Key does not match given key")
)

func handlePublicKeyCallback(conn ssh.ConnMetadata, key ssh.PublicKey) (perm *ssh.Permissions, err error) {

	publicKeys, err := userx.ReadUserAuthKeys(conn.User())
	if err != nil {

		return
	}
	for _, v := range publicKeys {

		if bytes.Compare(v.Marshal(), key.Marshal()) != 0 { // Use certchecker

			continue
		}
	}

	perm = &ssh.Permissions{}

	return
}

// Log all authentication attempts
func handleAuthLogCallback(conn ssh.ConnMetadata, method string, err error) {

	log.Printf("Login attempt for user %s from %s with client %s",
		conn.User(),
		conn.RemoteAddr(),
		string(conn.ClientVersion()))
	return
}

func handleServerConn(conn ssh.ConnMetadata, chans <-chan ssh.NewChannel) {

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
		go s.handleChannelRequest(channel, requests)
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

func (s *Session) handleChannelRequest(channel ssh.Channel, in <-chan *ssh.Request) {

	for req := range in {

		switch req.Type {

		case "shell":

			if len(req.Payload) == 0 {

				if err := req.Reply(true, nil); err != nil {

					log.Printf("Unable to reply to channel request (%s)", err)
					continue
				}
				go s.handleNewTerminal(channel)
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

func (s *Session) handleNewTerminal(channel ssh.Channel) {

	cmd := exec.Command("/bin/bash")
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: 1000, Gid: 1000}

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
