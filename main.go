// fluffy-tribble project main.go
package main

import (
	"flag"
	"fmt"
	"os"
)

var isServer bool
var ftDir string = ".fluffy-tribble"
var ftSecretFile = "secret"
var ftServerFile = "server"
var ftClientFile = "client"
var defaultSecretSize = 1024

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

	if isServer {
		runServer()
	} else {
		runClient()
	}
}

func usage() {
	fmt.Println("Usage instructions here.")
}
