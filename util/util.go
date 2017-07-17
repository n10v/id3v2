// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

const (
	// ID3SizeLen is length of ID3v2 size format (4 * 0b0xxxxxxx).
	ID3SizeLen = 4

	maxSize  = 268435455 // == 0b00001111 11111111 11111111 11111111
	sizeBase = 7
)

var (
	ErrSizeOverflow      = errors.New("size of tag/frame is greater than allowed in id3 tag")
	ErrInvalidSizeFormat = errors.New("invalid format of tag's/frame's size")
)

// WriteBytesSize writes size to bw in form of ID3v2 size format (4 * 0b0xxxxxxx).
//
// If size is greater than allowed (256MB), then it returns ErrSizeOverflow.
func WriteBytesSize(bw *bufio.Writer, size int) error {
	if size > maxSize {
		return ErrSizeOverflow
	}

	// First 4 bits of size are always "0", because size should be smaller
	// as maxSize. So skip them.
	size <<= 4

	// Let's explain the algorithm on example.
	// E.g. size is 32-bit integer and after the skip of first 4 bits
	// its value is "10100111 01110101 01010010 11110000".
	// In loop we should write every first 7 bits to bw.
	mask := 254 << (3 * 8) // 11111110 000000000 000000000 000000000
	for i := 0; i < ID3SizeLen; i++ {
		// To take first 7 bits we should do `size&mask`.
		firstBits := size & mask
		// firstBits is "10100110 00000000 00000000 00000000" now.
		// firstBits has int type, but we should have a byte.
		// To have a byte we should move first 7 bits to the end of firstBits,
		// because by converting int to byte only last 7 bits are taken.
		firstBits >>= (3*8 + 1)
		// firstBits is "00000000 00000000 00000000 01010011" now.
		bSize := byte(firstBits)
		// Now in bSize we have only "01010011". We can write it to bw.
		if err := bw.WriteByte(bSize); err != nil {
			return err
		}
		// Do the same with next 7 bits.
		size <<= sizeBase
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
