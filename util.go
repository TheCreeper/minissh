package main

import (
	"encoding/binary"
	"io"
	"log"
)

func CopyRequest(in io.ReadCloser, out io.WriteCloser) {

	written, err := io.Copy(out, in)
	if Verbose {

		log.Printf("%d bytes written", written)
	}
	if err != nil {

		log.Printf("CopyRequest: io.Copy: %s\n", err)
	}
	defer in.Close()
}

func CopyResponse(out io.ReadCloser, in io.WriteCloser) {

	written, err := io.Copy(in, out)
	if Verbose {

		log.Printf("%d bytes written", written)
	}
	if err != nil {

		log.Printf("CopyResponse: io.Copy: %s\n", err)
	}
	defer in.Close()
}

func ParseTerminalDims(b []byte) (uint32, uint32) {

	w := binary.BigEndian.Uint32(b)
	h := binary.BigEndian.Uint32(b[4:])
	return w, h
}
