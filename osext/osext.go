package osext

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"code.google.com/p/go.crypto/ssh"
)

var (
	// Match private host keys
	ValidPrivateHostKey = regexp.MustCompile(`^ssh_host_(\w)+_key$`)
)

func ReadHostKey(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func ReadHostKeys(filename string) (keys []ssh.Signer, err error) {
	files, err := ioutil.ReadDir(filename)
	if err != nil {
		return
	}
	for _, v := range files {
		if v.IsDir() {
			continue
		}

		if !(ValidPrivateHostKey.MatchString(v.Name())) {
			continue
		}

		b, err := ReadHostKey(filepath.Join(filename, v.Name()))
		if err != nil {
			return nil, err
		}

		signer, err := ssh.ParsePrivateKey(b)
		if err != nil {
			return nil, err
		}
		keys = append(keys, signer)
	}

	if len(keys) < 1 {
		return nil, errors.New("No Host Keys Found!")
	}
	return
}
