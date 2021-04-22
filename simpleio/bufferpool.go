package simpleio

import (
	"bytes"
	//"sync"
	"fmt"
)

func PrintString(buf *bytes.Buffer) {
	fmt.Printf("%s\n", buf.String())
}

type SizedBufferPool struct {
	c chan *bytes.Buffer
	a int
}

func NewSizedBufferPool(size int, alloc int) (bp *SizedBufferPool) {
	return &SizedBufferPool{
		c: make(chan *bytes.Buffer, size),
		a: alloc,
	}
}

// Get gets a Buffer from the SizedBufferPool, or creates a new one if none are
// available in the pool. Buffers have a pre-allocated capacity.
func (bp *SizedBufferPool) Get() (b *bytes.Buffer) {
	select {
	case b = <-bp.c:
		// reuse existing buffer
	default:
		// create new buffer
		b = bytes.NewBuffer(make([]byte, 0, bp.a))
	}
	return
}

// Put returns the given Buffer to the SizedBufferPool.
func (bp *SizedBufferPool) Put(b *bytes.Buffer) {
	b.Reset()

	// Release buffers over our maximum capacity and re-create a pre-sized
	// buffer to replace it.
	if cap(b.Bytes()) > bp.a {
		b = bytes.NewBuffer(make([]byte, 0, bp.a))
	}

	select {
	case bp.c <- b:
	default: // Discard the buffer if the pool is full.
	}
}
