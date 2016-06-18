package frame

import (
	"bytes"
	"sync"
)

var bytesBufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

// Framer provides a generic interface for frames.
type Framer interface {
	Bytes() ([]byte, error)
}
