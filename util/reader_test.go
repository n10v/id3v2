package util

import (
	"bytes"
	"fmt"
	"testing"
)

var (
	bs    = []byte{0, 11, 22, 33, 44, 55, 66, 77}
	bsBuf = bytes.NewReader(bs)
	n     = len(bs)
)

func TestReadBytes(t *testing.T) {
	bsReader := NewReader(bsBuf)

	read, err := bsReader.ReadBytes(n)
	if err != nil {
		t.Error(err)
	}
	if n != len(read) {
		t.Errorf("Expecting to read: %v, got: %v", n, read)
	}
	if !bytes.Equal(bs, read) {
		t.Error("Expecting: %v, got: %v", bs, read)
	}

	fmt.Println(bsReader.buf.Buffered())
	if bsReader.buf.Buffered() != 0 {
		t.Errorf("Expecting buffered: %v, got: %v", 0, bsReader.buf.Buffered())
	}
}
