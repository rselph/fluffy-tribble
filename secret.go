// secret.go
package main

import (
	"crypto/rand"
	"io"
	"log"
	"os"
)

type Secret []byte

func newSecret(size int) *Secret {
	bytes := Secret(make([]byte, size))

	bytesRead := 0
	for bytesRead < size {
		bytesLeft := bytes[bytesRead:]
		n, err := rand.Read(bytesLeft)
		if err != nil {
			log.Fatal(err)
		}
		bytesRead += n
	}

	return &bytes
}

func (s *Secret) save(w io.Writer) (err error) {
	bytesWritten := 0

	var n int
	for bytesWritten < len(*s) && err == nil {
		bytesLeft := (*s)[bytesWritten:]
		n, err = w.Write(bytesLeft)
		bytesWritten += n
	}

	return
}

func (s *Secret) saveTo(fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	err = os.Chmod(fname, 0600)
	if err != nil {
		return err
	}

	s.save(f)

	return nil
}
