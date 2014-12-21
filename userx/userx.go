package main

import (
	"bytes"
	"io"
	"io/ioutil"
)

var (

	// Value for fields with no value
	Unknown = "Unknown"

	// Location of the password file
	PasswordFile = ""
)

func init() {

	os := getHostOS()
	switch os {

	case "linux":

		PasswordFile = "/etc/passwd"
	}
}

type User struct {
	Username string
	Password string
	UserID   int
	GroupID  int
	UserInfo string
	HomeDir  string
	Shell    string
}

func LookupUser(user string) (uacc *User, err error) {

	buf, err := ioutil.ReadFile(PasswordFile)
	if err != nil {

		return
	}

	bb := bytes.NewBuffer(buf)
	for {

		line, err = bb.ReadBytes('\n')
		if err == io.EOF {

			break
		}
		if err != nil {

			return nil, err
		}

		uacc, err = ParsePasswdLine(line)
		if err != nil {

			return
		}
	}

	return
}

func ParsePasswdLine(line []byte) (uacc *User, err error) {

	bb := bytes.NewBuffer(line)
	for {

		field, err := bb.ReadBytes(':')
		if err == io.EOF {

			break
		}
		if err != nil {

			return nil, err
		}
		if len(field) < 1 {

			field = []byte(Unknown)
		}

		if len(uacc.Username) < 1 {

			uacc.Username = string(field)
		}
		if len(uacc.Password) < 1 {

			uacc.Password = string(field)
		}
		if len(uacc.UserID) < 1 {

			uacc.UserID = string(field)
		}
		if len(uacc.GroupID) < 1 {

			uacc.GroupID = string(field)
		}
		if len(uacc.UserInfo) < 1 {

			uacc.UserInfo = string(field)
		}
		if len(uacc.HomeDir) < 1 {

			uacc.HomeDir = string(field)
		}
		if len(uacc.Shell) < 1 {

			uacc.Shell = string(field)
		}
	}

	return
}
