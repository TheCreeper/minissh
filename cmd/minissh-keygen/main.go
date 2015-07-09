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

var (
	KeySize    int
	KeyType    string
	Password   string
	Comment    string
	OutputFile string
)

func init() {
	flag.IntVar(&KeySize, "b", 256, "Specifies the number of bits in the key to create. For ecdsa keys the possible values are 256, 384, and 521.")
	flag.StringVar(&KeyType, "t", "ecdsa", "Specifies the type of key to create. Only option is ecdsa.")
	flag.StringVar(&Password, "P", "", "The password used to encrypt the private key. AES-128 is used.")
	flag.StringVar(&OutputFile, "f", os.ExpandEnv("$HOME/.ssh/id_ecdsa"), "Specifies the filename of the key file. This will be where private/public keypair files are writen.")
	flag.Parse()
}

func main() {
	// Check if the private key already exists at the specified location.
	if _, err := os.Stat(OutputFile); err == nil {
		log.Fatal("Keyfile already exists at " + OutputFile)
	}

	// Check if public key already exits.
	if _, err := os.Stat(OutputFile + ".pub"); err == nil {
		log.Fatal("Keyfile already exists at " + OutputFile + ".pub")
	}

	public, private, err := GenerateKeyPair(KeyType, KeySize, Password)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(os.ExpandEnv(OutputFile), private, 0600); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(os.ExpandEnv(OutputFile+".pub"), public, 0640); err != nil {
		log.Fatal(err)
	}
}

func GenerateKeyPair(keytype string, keysize int, password string) (public, private []byte, err error) {
	// Add more key types in future
	switch keytype {
	case "ecdsa":
		return GenerateECDSAKeyPair(keysize, password)
	}
	return
}

func GenerateECDSAKeyPair(keysize int, password string) (public, private []byte, err error) {
	var curve elliptic.Curve
	switch keysize {
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

	// Marshal the public key
	sshPubKey, err := ssh.NewPublicKey(&prvKey.PublicKey)
	if err != nil {
		return
	}
	public = ssh.MarshalAuthorizedKey(sshPubKey)

	// Marshal the private key
	prvKeyDer, err := x509.MarshalECPrivateKey(prvKey)
	if err != nil {
		return
	}
	block := &pem.Block{Type: "EC PRIVATE KEY", Bytes: prvKeyDer}

	// Encrypt the private key
	if len(password) != 0 {
		// AES-128 is the only option for private key encryption just like in ssh-keygen.
		block, err = x509.EncryptPEMBlock(rand.Reader,
			"EC PRIVATE KEY",
			prvKeyDer,
			[]byte(password),
			x509.PEMCipherAES128)
		if err != nil {
			return
		}
	}

	private = pem.EncodeToMemory(block)
	return
}
