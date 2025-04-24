package cache

import (
	"fmt"
)

type Hash uint64

type Hashable interface {
	Hash() Hash
	Equals(v any) bool
}

func (h Hash) String() string {
	return fmt.Sprintf("%x", uint64(h))
}

func (h Hash) Equals(other any) bool {
	if o, ok := other.(Hash); ok {
		return h == o
	}
	return false
}

func (h Hash) Hash() Hash {
	return h
}

func (h Hash) NumSetBits(startBit, bitStep int) int {
	n := 0
	for i := startBit; i < 64; i += bitStep {
		if (h & (1 << uint(i))) != 0 {
			n++
		}
	}
	return n
}
