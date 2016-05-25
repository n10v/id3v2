package frame

import (
	"bytes"
	"sync"
)

var bytesBufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

type Framer interface {
	Form() []byte
}
