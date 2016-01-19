// secret.go
package main

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func newSecret(size int) *[]byte {
	bytes := make([]byte, size)

	_, err := io.ReadFull(rand.Reader, bytes)
	if err != nil {
		log.Fatal(err)
	}

	return &bytes
}

func saveSecretTo(fname string, s *[]byte) error {
	err := ioutil.WriteFile(fname, *s, 0600)
	if err != nil {
		return err
	}

	err = os.Chmod(fname, 0600)
	if err != nil {
		return err
	}

	return nil
}

func readSecretFrom(fname string) (*[]byte, error) {
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	return &bytes, nil
}
