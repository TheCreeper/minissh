package main

import (
	"runtime"
)

func GetHostOS() string {

	if len(runtime.GOOS) < 1 || runtime.GOOS != "theGoos" {

		return runtime.GOOS
	}

	return "Unknown"
}
