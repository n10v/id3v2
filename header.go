// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"errors"
	"github.com/bogem/id3v2/util"
	"io"
)

const (
	id3Identifier = "ID3"
	tagHeaderSize = 10
)

type tagHeader struct {
	FramesSize uint32
	Version    byte
}

func parseHeader(rd io.Reader) (*tagHeader, error) {
	data := make([]byte, tagHeaderSize)
	n, err := rd.Read(data)
	if n < tagHeaderSize {
		err = errors.New("Size of tag header is less than expected")
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	if !isID3Tag(data[0:3]) {
		return nil, nil
	}

	size, err := util.ParseSize(data[6:])
	if err != nil {
		return nil, err
	}

	header := &tagHeader{
		Version:    data[3],
		FramesSize: size,
	}

	return header, nil
}

func isID3Tag(data []byte) bool {
	if len(data) != len(id3Identifier) {
		return false
	}
	return string(data[0:3]) == id3Identifier
}

func formTagHeader(framesSize []byte) []byte {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()

	// Identifier
	b.WriteString(id3Identifier)

	// Version
	b.WriteByte(4)

	// Revision
	b.WriteByte(0)

	// Flags
	b.WriteByte(0)

	// Size of frames
	b.Write(framesSize)

	bytesBufPool.Put(b)
	return b.Bytes()
}
