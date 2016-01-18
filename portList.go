// portList.go
package main

import (
	"crypto/sha512"
	"encoding/binary"
	"time"
)

type PortList struct {
	ports      [][]int
	knockState []int
}

func newPortList() *PortList {
	retval := &PortList{}
	retval.ports = make([][]int, 2)
	retval.knockState = make([]int, knockSequenceLength)
	return retval
}

func (p *PortList) update(t time.Time, secret *Secret) {
	epochTime := t.Unix() / int64(refreshInterval.Seconds())
	p.ports[1] = p.ports[0]
	p.ports[0] = make([]int, knockSequenceLength)

	hasher := sha512.New()

	hasher.Reset()
	binary.Write(hasher, binary.LittleEndian, epochTime)
	result := hasher.Sum(nil)

	hasher.Reset()
	hasher.Write(result)
	hasher.Write([]byte(*secret))
	master := hasher.Sum(nil)

	for _, i := range p.ports[0] {
		hasher.Reset()
		hasher.Write(master)
		binary.Write(hasher, binary.LittleEndian, i)
		finalHash := hasher.Sum(nil)
	}
}
