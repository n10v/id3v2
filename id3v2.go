package id3v2

import (
	"bytes"
	"github.com/bogem/id3v2/frame"
	"os"
	"sync"
)

var bytesBufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func Open(name string) (*Tag, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return ParseTag(file)
}

func NewAttachedPicture() *frame.PictureFrame {
	return frame.NewPictureFrame()
}
