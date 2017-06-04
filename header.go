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
var ErrSmallHeaderSize = errors.New("size of tag header is less than expected")

type tagHeader struct {
	FramesSize int64
	Version    byte
}

// parseHeader parses tag header in rd.
// If there is no tag in rd, it returns errNoTag.
// If rd is smaller than expected, it returns ErrSmallHeaderSize.
func parseHeader(rd io.Reader) (tagHeader, error) {
	var header tagHeader

	data := make([]byte, tagHeaderSize)
	n, err := rd.Read(data)
	if err != nil {
		return header, err
	}
	if n < tagHeaderSize {
		return header, ErrSmallHeaderSize
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

// writeTagHeader writes tag header with framesSize and version to bw.
// buf is needed for util.WriteBytesSize.
func writeTagHeader(bw *bufio.Writer, buf []byte, framesSize int, version byte) error {
	if len(buf) < util.ID3SizeLen {
		return errors.New("writeTagHeader: buf size is less than id3v2 size length")
	}
	size := buf[:util.ID3SizeLen]

	if err := util.WriteBytesSize(size, framesSize); err != nil {
		return err
	}

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
	if _, err := bw.Write(size); err != nil {
		return err
	}

	return nil
}
