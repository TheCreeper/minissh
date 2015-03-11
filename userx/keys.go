package userx

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"os/user"
	"path"
	"regexp"
	"strings"

	"code.google.com/p/go.crypto/ssh"
)

var (

	// Match private host keys
	ValidPrivateHostKey = regexp.MustCompile(`^ssh_host_(\w)+_key$`)
)

var (
	ErrNoHostKeys = errors.New("No Host Keys Found!")
)

func GenHostKeyRSA() []byte {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {

		panic(err)
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{

		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}
	return pem.EncodeToMemory(&privateKeyBlock)
}

func ReadHostKey(p string) (privateKey ssh.Signer, err error) {

	buf, err := ioutil.ReadFile(p)
	if err != nil {

		return
	}

	return ssh.ParsePrivateKey(buf)
}

func ReadHostKeys(p string) (keys []ssh.Signer, err error) {

	fileinfo, err := ioutil.ReadDir(p)
	if err != nil {

		return
	}
	for _, v := range fileinfo {

		if v.IsDir() {

			continue
		}

		if !(ValidPrivateHostKey.MatchString(v.Name())) {

			continue
		}

		s, err := ReadHostKey(path.Join(p, v.Name()))
		if err != nil {

			return nil, err
		}
		keys = append(keys, s)
	}

	if len(keys) < 1 {

		return nil, ErrNoHostKeys
	}

	return
}

func ReadUserAuthKeys(u string) (publicKeys []ssh.PublicKey, err error) {

	uacc, err := user.Lookup(u)
	if err != nil {

		return
	}

	buf, err := ioutil.ReadFile(path.Join(uacc.HomeDir, ".ssh/authorized_keys"))
	if err != nil {

		return
	}

	bb := bytes.NewBuffer(buf)
	for {

		c, err := bb.ReadBytes('\n')
		if err == io.EOF {

			break
		}
		if err != nil {

			return nil, err
		}
		if strings.HasPrefix(string(c), "#") {

			continue
		}

		publicKey, _, _, _, err := ssh.ParseAuthorizedKey(c)
		if err != nil {

			return nil, err
		}
		publicKeys = append(publicKeys, publicKey)
	}

	return
}
