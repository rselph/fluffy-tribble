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
var refreshInterval = 10 * time.Second
var remoteHost = "localhost"
var connectTimeout = 5 * time.Second

func main() {
	flag.Usage = usage
	flag.CommandLine.SetOutput(os.Stdout)
	flag.BoolVar(&isServer, "s", false, "Run as server daemon")
	flag.Parse()
	if len(flag.Args()) > 0 {
		remoteHost = flag.Arg(0)
	}

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
		var cmdLine []string
		if len(flag.Args()) > 1 {
			cmdLine = flag.Args()[1:]
		} else {
			cmdLine = []string{ftClientFile}
		}

		runClient(s, cmdLine)
	}
}

func usage() {
	fmt.Println("Usage instructions here.")
}
