// server.go
package main

import (
	"time"
)

func runServer(s *[]byte) {
	ports := newPortList()
	startTime := time.Now()
	ports.update(startTime.Add(-refreshInterval), s)
	ports.update(startTime, s)

	ticktock := time.Tick(refreshInterval * time.Second)
	for true {
		select {
		case <-ticktock:

		}
	}
}
