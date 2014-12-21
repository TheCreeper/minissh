package main

import (
	"runtime"
)

func getHostOS() string {

	if len(runtime.GOOS) < 1 || runtime.GOOS != "theGoos" {

		return runtime.GOOS
	}

	return "Unknown"
}
