// portList.go
package main

import (
	"bytes"
	"crypto"
	_ "crypto/sha256"
	"encoding/binary"
	"reflect"
	"time"
)

type PortList struct {
	newList        []int
	current        []int
	history        [][]int
	secret         []byte
	sequenceLength int
	mask           int
	hi, lo         int
	testSequence   []int
}

func newPortList(secret []byte, sequenceLength, historyLength, hi, lo int) *PortList {
	p := &PortList{
		secret:         secret,
		sequenceLength: sequenceLength,
		hi:             hi,
		lo:             lo,
	}
	p.history = make([][]int, historyLength)
	p.testSequence = make([]int, sequenceLength)

	mask := 1
	for mask <= (hi - lo) {
		mask = mask << 1
	}
	mask -= 1
	p.mask = mask

	return p
}

func interval(t time.Time) int64 {
	return t.Unix() / int64(refreshInterval.Seconds())
}

func (p *PortList) update(intervalNum int64) {
	p.newList = make([]int, p.sequenceLength)

	hasher := crypto.SHA256.New()

	hasher.Reset()
	_ = binary.Write(hasher, binary.LittleEndian, intervalNum)
	result := hasher.Sum(nil)

	hasher.Reset()
	hasher.Write(result)
	hasher.Write(p.secret)
	master := hasher.Sum(nil)

	n := int64(0)
	for i := range p.newList {
		p.newList[i], n = nextPort(master, n, p)
	}

	copy(p.history[1:], p.history)
	p.history[0] = p.newList
	p.current = p.newList
	p.newList = nil
}

func nextPort(master []byte, n int64, p *PortList) (int, int64) {
	hasher := crypto.SHA256.New()

	port := p.hi
	for shouldRejectPort(port, p) {
		hasher.Reset()
		hasher.Write(master)
		_ = binary.Write(hasher, binary.LittleEndian, n)
		finalHash := hasher.Sum(nil)

		var portTmp int64
		_ = binary.Read(bytes.NewReader(finalHash), binary.LittleEndian, &portTmp)
		port = int(portTmp)

		port &= p.mask
		port += p.lo

		n += 1
	}

	return port, n
}

func shouldRejectPort(port int, p *PortList) bool {
	if port < p.lo || port >= p.hi {
		return true
	}

	for _, oldPort := range p.newList {
		if port == oldPort {
			return true
		}
	}

	for _, list := range p.history {
		for _, oldPort := range list {
			if port == oldPort {
				return true
			}
		}
	}

	return false
}

func (p *PortList) checkFull(port int) bool {
	p.testSequence = append(p.testSequence[1:], port)

	for _, list := range p.history {
		if reflect.DeepEqual(list, p.testSequence) {
			return true
		}
	}

	return false
}
