package sshutil

import (
	"io/ioutil"

	"code.google.com/p/go.crypto/ssh"
)

func ReadHostKey(filename string) (ssh.Signer, error) {

	b, err := ioutil.ReadFile(filename)
	if err != nil {

		return nil, err
	}

	return ssh.ParsePrivateKey(b)
}
