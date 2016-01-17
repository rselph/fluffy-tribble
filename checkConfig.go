// checkConfig.go
package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func isSafeConfig() bool {
	me, err := user.Current()
	if err != nil {
		fmt.Println(os.Stderr, err)
		return false
	}

	ftDir = filepath.Join(me.HomeDir, ftDir)

	dirInfo, err := os.Lstat(ftDir)
	if err != nil {
		fmt.Println(os.Stderr, err)
		fmt.Println("Creating directory.")
		err = os.MkdirAll(ftDir, 0700)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return false
		}
		dirInfo, err = os.Lstat(ftDir)
	}

	if !dirInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "%s is not a directory.\n", ftDir)
		return false
	}

	if dirInfo.Mode().Perm()&077 != 0 {
		fmt.Fprintf(os.Stderr, "Wrong permissions on %s.  Should be accessible only to owner.\n", ftDir)
		return false
	}

	ftSecretFile = filepath.Join(ftDir, ftSecretFile)
	secretInfo, err := os.Lstat(ftSecretFile)
	if err != nil {
		fmt.Println(os.Stderr, err)
		fmt.Println("Creating new secret.")

		secret := newSecret(defaultSecretSize)
		err = secret.saveTo(ftSecretFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return false
		}
		secretInfo, err = os.Lstat(ftSecretFile)
	}

	if secretInfo.Mode().Perm()&077 != 0 {
		fmt.Fprintf(os.Stderr, "Wrong permissions on %s.  Should be accessible only to owner.\n", ftSecretFile)
		return false
	}

	var scriptFile string
	if isServer {
		ftServerFile = filepath.Join(ftDir, ftServerFile)
		scriptFile = ftServerFile
	} else {
		ftClientFile = filepath.Join(ftDir, ftClientFile)
		scriptFile = ftClientFile
	}

	scriptInfo, err := os.Lstat(scriptFile)
	if err != nil {
		fmt.Println(os.Stderr, err)
		return false
	}

	if scriptInfo.Mode().Perm()&077 != 0 {
		fmt.Fprintf(os.Stderr, "Wrong permissions on %s.  Should be accessible only to owner.\n", scriptFile)
		return false
	}

	return true
}
