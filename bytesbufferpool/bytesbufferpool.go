package bytesbufferpool

import (
	"bytes"
	"sync"
)

var bytesBufferPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func Get() *bytes.Buffer {
	b := bytesBufferPool.Get().(*bytes.Buffer)
	b.Reset()
	return b
}

func Put(b *bytes.Buffer) {
	bytesBufferPool.Put(b)
}
