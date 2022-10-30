// client.go
package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func runClient(s []byte, cmdLine []string) {
	startedAt := time.Now()

	succeeded := false
	for !succeeded {
		ports := newPortList(s, knockSequenceLength, 3, portRangeHigh, portRangeLow)
		nowInterval := interval(time.Now())
		ports.update(nowInterval - 3)
		ports.update(nowInterval - 2)

		for delta := int64(-1); delta < 2; delta += 1 {
			err := tryInterval(nowInterval+delta, ports)
			if err == nil {
				succeeded = true
				break
			}
		}

		if !succeeded {
			if time.Since(startedAt) >= connectTimeout {
				fmt.Println("CONNECT FAILED")
				os.Exit(1)
			}
			time.Sleep(250 * time.Millisecond)
		}
	}

	//	fmt.Println("CONNECT SUCCESSFUL")
	cmd := exec.Command(cmdLine[0], cmdLine[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = append(os.Environ(), fmt.Sprintf("FLUFFY_TRIBBLE_HOST=%s", remoteHost))
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func tryInterval(i int64, ports *PortList) error {
	ports.update(i)

	buf := make([]byte, 1)
	for _, port := range ports.current {
		dialString := net.JoinHostPort(remoteHost, strconv.Itoa(port))
		conn, err := net.Dial("tcp", dialString)
		if err != nil {
			return err
		}

		for err == nil {
			_, err = conn.Read(buf)
		}
		if err != io.EOF {
			return err
		}
		conn.Close()
	}

	return nil
}
