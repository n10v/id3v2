// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"errors"
	"io"

	"github.com/bogem/id3v2/util"
)

const (
	id3Identifier = "ID3"
	tagHeaderSize = 10
)

var errNoTag = errors.New("there is no tag in file")

type tagHeader struct {
	FramesSize int64
	Version    byte
}

func parseHeader(rd io.Reader) (tagHeader, error) {
	var header tagHeader

	data := make([]byte, tagHeaderSize)
	n, err := rd.Read(data)
	if err != nil {
		return header, err
	}
	if n < tagHeaderSize {
		err = errors.New("size of tag header is less than expected")
		return header, err
	}

	if !isID3Tag(data[0:3]) {
		return header, errNoTag
	}

	size, err := util.ParseSize(data[6:])
	if err != nil {
		return header, err
	}

	header.Version = data[3]
	header.FramesSize = size
	return header, nil
}

func isID3Tag(data []byte) bool {
	if len(data) != len(id3Identifier) {
		return false
	}
	return string(data[0:3]) == id3Identifier
}

func writeTagHeader(bw *bufio.Writer, framesSize []byte, version byte) error {
	// Identifier
	if _, err := bw.WriteString(id3Identifier); err != nil {
		return err
	}

	// Version
	if err := bw.WriteByte(version); err != nil {
		return err
	}

	// Revision
	if err := bw.WriteByte(0); err != nil {
		return err
	}

	// Flags
	if err := bw.WriteByte(0); err != nil {
		return err
	}

	// Size of frames
	if _, err := bw.Write(framesSize); err != nil {
		return err
	}

	return nil
}
