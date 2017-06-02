// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"bytes"
	"errors"
	"io"
)

const (
	// ID3SizeLen is length of ID3v2 size specification format (4 * 0b0xxxxxxx).
	ID3SizeLen = 4

	maxSize  = 268435455 // == 0b11111... (28 digits)
	sizeBase = 7
)

var (
	ErrSizeOverflow      = errors.New("size of tag/frame is greater than allowed in id3 tag")
	ErrDstIsSmall        = errors.New("len(dst) is small. It should be >= 4")
	ErrInvalidSizeFormat = errors.New("invalid format of tag's/frame's size")
)

// WriteBytesSize writes size to last 4 bytes of dst in form
// of ID3v2 size specification format (4 * 0b0xxxxxxx).
//
// If len(dst) is smaller than 4, it returns ErrDstIsSmall.
// If size is greater than allowed (256MB), then it returns ErrSizeOverflow.
func WriteBytesSize(dst []byte, size int) error {
	if len(dst) < ID3SizeLen {
		return ErrDstIsSmall
	}
	if size > maxSize {
		return ErrSizeOverflow
	}

	mask := 1<<sizeBase - 1 // == 0b01111111
	for i := 0; i < ID3SizeLen; i++ {
		dst[len(dst)-1-i] = byte(size & mask)
		size >>= sizeBase
	}

	return nil
}

// ParseSize parses data in form of ID3v2 size specification format (4 * 0b0xxxxxxx)
// and returns parsed int64 number.
//
// If length of data is greater than 4 or if there is invalid size format (e.g.
// one byte in data is like 0b1xxxxxxx), then it returns ErrInvalidSizeFormat.
func ParseSize(data []byte) (int64, error) {
	if len(data) > ID3SizeLen {
		return 0, ErrInvalidSizeFormat
	}

	var size int64
	for _, b := range data {
		if b&128 > 0 { // 128 = 0b1000_0000
			return 0, ErrInvalidSizeFormat
		}

		size = (size << sizeBase) | int64(b)
	}

	return size, nil
}

// ReadAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF.
// Because ReadAll is defined to read from src until EOF,
// it does not treat an EOF from Read as an error to be reported.
func ReadAll(rd io.Reader) ([]byte, error) {
	if lr, ok := rd.(*io.LimitedReader); ok {
		buf := make([]byte, lr.N)
		_, err := lr.R.Read(buf)
		return buf, err
	}

	buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))
	_, err := buf.ReadFrom(rd)
	return buf.Bytes(), err
}
