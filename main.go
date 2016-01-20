// fluffy-tribble project main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var isServer bool
var ftDir string = ".fluffy-tribble"
var ftSecretFile = "secret"
var ftServerFile = "server"
var ftClientFile = "client"
var defaultSecretSize = 1024
var knockSequenceLength = 10
var portRangeLow = 20000
var portRangeHigh = 21000
var refreshInterval = 1 * time.Second
var remoteHost = "localhost"

func main() {
	fmt.Println("fluffy-tribble")

	flag.Usage = usage
	flag.CommandLine.SetOutput(os.Stdout)
	flag.BoolVar(&isServer, "s", false, "Run as server daemon")
	flag.Parse()

	if !isSafeConfig() {
		fmt.Fprintln(os.Stderr, "Exiting.")
		os.Exit(1)
	}

	s, err := readSecretFrom(ftSecretFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if isServer {
		runServer(s)
	} else {
		runClient(s)
	}
}

func usage() {
	fmt.Println("Usage instructions here.")
}
