package main

import (
	"bytes"
	"errors"
	"log"

	"code.google.com/p/go.crypto/ssh"
)

// Some errors
var (
	ErrKeyMatch = errors.New("Authorised Key does not match given key")
)

func HandlePublicKeyCallback(conn ssh.ConnMetadata, key ssh.PublicKey) (perm *ssh.Permissions, err error) {

	publicKeys, err := ReadUserAuthKeys(conn.User())
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
func HandleAuthLogCallback(conn ssh.ConnMetadata, method string, err error) {

	log.Printf("Login attempt for user %s from %s with client %s",
		conn.User(),
		conn.RemoteAddr().String(),
		string(conn.ClientVersion()))
	return
}
