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

	wg := &sync.WaitGroup{}
	events := make(chan int)
	var currentListeners, previousListeners []net.Listener

	updateListenState := func(now time.Time) error {
		ports.update(now)
		fmt.Println(ports.current)
		newListeners, err := openPorts(ports.current, events, wg)
		if err != nil {
			return err
		}
		closePorts(previousListeners)
		previousListeners = currentListeners
		currentListeners = newListeners
		return nil
	}

	startTime := time.Now()
	updateListenState(startTime.Add(-refreshInterval))
	updateListenState(startTime)

	ticktock := time.Tick(refreshInterval)

	for true {
		select {
		case now := <-ticktock:
			err := updateListenState(now)
			if err != nil {
				fmt.Println("Open ports failed.")
				fmt.Println(err)
			}

		case port := <-events:
			fmt.Printf(">%v\n", port)
			if ports.checkFull(port) {
				// do server stuff!!
				fmt.Println("KNOCK SUCCEEDED")
			}
		}
	}
}
