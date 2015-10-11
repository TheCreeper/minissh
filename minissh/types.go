package minissh

import (
	"os"

	"code.google.com/p/go.crypto/ssh"
)

type Session struct {
	// SSH Conn Meta Data
	Conn ssh.ConnMetadata

	// pty
	pty *os.File

	// Win Size
	WinH uint16
	WinW uint16
}
