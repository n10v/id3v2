package frame

import (
	"bytes"
	"github.com/bogem/id3v2/util"
)

const (
	FrameHeaderSize = 10
	EncodingSize    = 1
)

type Framer interface {
	Form() ([]byte, error)
	ID() string
	SetID(string)
	Size() uint32
}

type FrameHeader struct {
	ID        string
	FrameSize uint32
}

func FormFrameHeader(f Framer) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString(f.ID())
	if size, err := util.FormSize(f.Size()); err != nil {
		return nil, err
	} else {
		b.Write(size)
	}
	b.WriteByte(0)
	b.WriteByte(0)
	return b.Bytes(), nil
}
