// server.go
package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func runServer(s *[]byte) {
	historyLength := 3
	ports := newPortList(s, knockSequenceLength, historyLength, portRangeHigh, portRangeLow)

	wg := &sync.WaitGroup{}
	events := make(chan int)
	var listeners [][]net.Listener = make([][]net.Listener, historyLength)

	updateListenState := func(intervalNum int64) error {
		ports.update(intervalNum)
		fmt.Println(ports.current)
		newListeners, err := openPorts(ports.current, events, wg)
		if err != nil {
			return err
		}
		closePorts(listeners[0])
		listeners = append(listeners[1:], newListeners)
		return nil
	}

	last := interval(time.Now())
	updateListenState(last + 1)

	ticktock := time.Tick(refreshInterval / 2)

	for true {
		select {
		case now := <-ticktock:
			if present := interval(now); present != last {
				err := updateListenState(present + 1)
				if err != nil {
					fmt.Println("Open ports failed.")
					fmt.Println(err)
				}
				last = present
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
