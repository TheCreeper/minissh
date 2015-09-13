package osext

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
)

var (
	// Match private host keys
	ValidPrivateHostKey = regexp.MustCompile(`^ssh_host_(\w)+_key$`)

	// Match public host keys
	ValidPublicHostKey = regexp.MustCompile(`^ssh_host_(\w)+_key.pub$`)
)

func ReadHostKeys(filename string) (keys [][]byte, err error) {
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

		b, err := ioutil.ReadFile(filepath.Join(filename, v.Name()))
		if err != nil {
			return nil, err
		}
		keys = append(keys, b)
	}
	return
}
