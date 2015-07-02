package userext

import (
	"bytes"
	"io"
	"io/ioutil"
	"os/user"
	"path"
	"strings"

	"code.google.com/p/go.crypto/ssh"
)

func ReadUserAuthKeys(u string) (publicKeys []ssh.PublicKey, err error) {

	uacc, err := user.Lookup(u)
	if err != nil {

		return
	}

	b, err := ioutil.ReadFile(path.Join(uacc.HomeDir, ".ssh/authorized_keys"))
	if err != nil {

		return
	}

	buf := bytes.NewBuffer(b)
	for {

		c, err := buf.ReadBytes('\n')
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
