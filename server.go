// server.go
package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func runServer(s *[]byte) {
	fmt.Println("fluffy-tribble")

	historyLength := 3
	ports := newPortList(s, knockSequenceLength, historyLength, portRangeHigh, portRangeLow)

	wg := &sync.WaitGroup{}
	events := make(chan net.Conn)
	var listeners [][]net.Listener = make([][]net.Listener, historyLength)

	updateListenState := func(intervalNum int64) error {
		ports.update(intervalNum)
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

		case conn := <-events:
			laddr := conn.LocalAddr().String()
			portString := laddr[strings.LastIndex(laddr, ":")+1:]
			port, _ := strconv.Atoi(portString)
			if ports.checkFull(port) {
				//				fmt.Println("KNOCK SUCCEEDED")
				cmd := exec.Command(ftServerFile)
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				err := cmd.Start()
				time.Sleep(100 * time.Millisecond)
				conn.Close()
				if err != nil {
					fmt.Fprint(os.Stderr, err)
				} else {
					err = cmd.Wait()
					if err != nil {
						fmt.Fprint(os.Stderr, err)
					}
				}
			} else {
				conn.Close()
			}
		}
	}
}
