// server.go
package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func runServer(s *[]byte) {
	ports := newPortList(s, knockSequenceLength, portRangeHigh, portRangeLow)
	startTime := time.Now()
	ports.update(startTime.Add(-refreshInterval))
	ports.update(startTime)

	wg := &sync.WaitGroup{}
	events := make(chan int)
	var currentListeners, previousListeners []net.Listener

	ticktock := time.Tick(refreshInterval)
	for true {
		select {
		case now := <-ticktock:
			ports.update(now)
			newListeners, err := openPorts(ports.current, events, wg)
			if err != nil {
				fmt.Println("Open ports failed.")
				fmt.Println(err)
				continue
			}
			closePorts(previousListeners)
			previousListeners = currentListeners
			currentListeners = newListeners

		case port := <-events:
			fmt.Printf("<%v\n", port)
		}
	}
}
