package main

import (
	"code.google.com/p/go.crypto/ssh"
)

type Auth struct{ *ssh.CertChecker }

func NewAuth() (a Auth) {

	return
}
