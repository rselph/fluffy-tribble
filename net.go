// net.go
package main

import (
	"fmt"
	"net"
	"sync"
)

func openPorts(ports []int, events chan int, wg *sync.WaitGroup) ([]net.Listener, error) {
	returnListeners := make([]net.Listener, len(ports))
	for i, port := range ports {
		addr := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			closePorts(returnListeners)
			return nil, err
		}
		returnListeners[i] = ln
		wg.Add(1)
		go portWatcher(ln, port, events, wg)
	}
	return returnListeners, nil
}

func closePorts(ports []net.Listener) {
	for _, port := range ports {
		if port != nil {
			port.Close()
		}
	}
}

func portWatcher(ln net.Listener, port int, events chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for true {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		conn.Close()
		events <- port
	}
}
