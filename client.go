// client.go
package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func runClient(s *[]byte) {
	ports := newPortList(s, knockSequenceLength, 3, portRangeHigh, portRangeLow)

	nowInterval := interval(time.Now())
	ports.update(nowInterval - 3)
	ports.update(nowInterval - 2)

	succeeded := false
	var err error
	for !succeeded {
		for delta := int64(-1); delta < 2; delta += 1 {
			err = tryInterval(nowInterval+delta, ports)
			if err == nil {
				succeeded = true
				break
			}
			fmt.Printf("FAIL: %v\n", err)
		}

		if !succeeded {
			os.Exit(1)
			time.Sleep(250 * time.Millisecond)
			nowInterval = interval(time.Now())
		}
	}

	fmt.Println("CONNECT SUCCESSFUL")
}

func tryInterval(i int64, ports *PortList) error {
	ports.update(i)

	fmt.Println(ports.current)
	for _, port := range ports.current {
		dialString := net.JoinHostPort(remoteHost, fmt.Sprintf("%d", port))
		conn, err := net.Dial("tcp", dialString)
		if err != nil {
			return err
		}
		conn.Close()
	}

	return nil
}
