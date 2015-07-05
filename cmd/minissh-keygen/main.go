package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

const (
	PubKeyPostfix = ".pub"
)

var (
	KeySize    int
	KeyType    string
	Passphrase string
	Comment    string
	OutputFile string
)

func init() {
	flag.IntVar(&KeySize, "b", 256, "")
	flag.StringVar(&KeyType, "t", "ecdsa", "")
	flag.StringVar(&Passphrase, "p", "", "")
	flag.StringVar(&OutputFile, "f", "", "")
	flag.Parse()

	if len(OutputFile) == 0 {
		OutputFile = os.ExpandEnv("$HOME/.ssh/id_ecdsa")
	}
}

func main() {
	public, private, err := GenerateECDSAKeyPair(KeySize)
	if err != nil {
		log.Fatal(err)
	}

	if err = ioutil.WriteFile(os.ExpandEnv(OutputFile), private, 0600); err != nil {
		log.Fatal(err)
	}

	if err = ioutil.WriteFile(os.ExpandEnv(OutputFile+PubKeyPostfix), public, 0600); err != nil {
		log.Fatal(err)
	}
}

func GenerateECDSAKeyPair(size int) (public, private []byte, err error) {
	var curve elliptic.Curve
	switch size {
	case 256:
		curve = elliptic.P256()
	case 384:
		curve = elliptic.P384()
	case 521:
		curve = elliptic.P521()
	default:
		return
	}

	// Generate the public/private key pair
	prvKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return
	}

	// Marshal the private key
	prvKeyDer, err := x509.MarshalECPrivateKey(prvKey)
	if err != nil {
		return
	}

	prvKeyBlock := &pem.Block{Type: "EC PRIVATE KEY", Bytes: prvKeyDer}
	private = pem.EncodeToMemory(prvKeyBlock)

	// Marshal the public key
	sshPubKey, err := ssh.NewPublicKey(&prvKey.PublicKey)
	if err != nil {
		return
	}
	public = ssh.MarshalAuthorizedKey(sshPubKey)

	return
}
