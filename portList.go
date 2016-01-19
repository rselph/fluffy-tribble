// portList.go
package main

import (
	"bytes"
	"crypto"
	_ "crypto/sha256"
	"encoding/binary"
	"time"
)

type PortList struct {
	current        []int
	previous       []int
	secret         *[]byte
	sequenceLength int
	mask           int
	hi, lo         int
}

func newPortList(secret *[]byte, sequenceLength, hi, lo int) *PortList {
	retval := &PortList{
		secret:         secret,
		sequenceLength: sequenceLength,
		hi:             hi,
		lo:             lo,
	}

	mask := 1
	for mask <= (hi - lo) {
		mask = mask << 1
	}
	mask -= 1
	retval.mask = mask

	return retval
}

func (p *PortList) update(t time.Time) {
	epochTime := t.Unix() / int64(refreshInterval.Seconds())
	p.current = make([]int, p.sequenceLength)

	hasher := crypto.SHA256.New()

	hasher.Reset()
	binary.Write(hasher, binary.LittleEndian, epochTime)
	result := hasher.Sum(nil)

	hasher.Reset()
	hasher.Write(result)
	hasher.Write(*p.secret)
	master := hasher.Sum(nil)

	n := int64(0)
	for i := range p.current {
		p.current[i], n = nextPort(&master, n, p)
	}
}

func nextPort(master *[]byte, n int64, p *PortList) (int, int64) {
	hasher := crypto.SHA256.New()

	port := p.hi
	for port >= p.hi {
		hasher.Reset()
		hasher.Write(*master)
		binary.Write(hasher, binary.LittleEndian, n)
		finalHash := hasher.Sum(nil)

		var portTmp int64
		binary.Read(bytes.NewReader(finalHash), binary.LittleEndian, &portTmp)
		port = int(portTmp)

		port &= p.mask
		port += p.lo

		n += 1
	}

	return port, n
}
