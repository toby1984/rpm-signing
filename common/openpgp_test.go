package common

import (
	"fmt"
	"os"
	"testing"
)

func TestLoadPublicKey(t *testing.T) {

	fileName := "testdata/my-signing-key.pub.gpg.asc"
	var err error
	var file *os.File
	if file, err = os.Open(fileName); err != nil {
		t.Fatal(err)
	}
	entityList, err := LoadPublicKey(file)
	if err != nil {
		t.Fatal(err)
	}
	CloseQuietly(file)
	if len(entityList) == 0 {
		t.Fatal("Failed to load PGP key")
	}
	for _, value := range entityList {
		fmt.Printf("OpenPGP public key loaded: %v\n", value.Identities)
	}
}

func TestLoadPrivateKey(t *testing.T) {

	fileName := "testdata/my-signing-key.priv.gpg.asc"
	var err error
	var file *os.File
	if file, err = os.Open(fileName); err != nil {
		t.Fatal(err)
	}

	reader := NewSeekableReader(NewIOReaderWithOffset(file))
	privateKey, err := LoadPrivateKey(reader, []byte{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("OpenPGP private key loaded: %v\n", privateKey.Fingerprint)
}
