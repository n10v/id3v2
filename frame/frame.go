package frame

import (
	"bytes"
	"github.com/bogem/id3v2/util"
	"sync"
)

var bytesBufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

type Framer interface {
	Form() []byte
	ID() string
	SetID(string)
}

type FrameHeader struct {
	ID        string
	FrameSize uint32
}

func FormFrameHeader(f Framer, size uint32) ([]byte, error) {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()
	b.WriteString(f.ID())
	if size, err := util.FormSize(size); err != nil {
		return nil, err
	} else {
		b.Write(size)
	}
	b.WriteByte(0)
	b.WriteByte(0)
	bytesBufPool.Put(b)
	return b.Bytes(), nil
}
