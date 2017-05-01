// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import "errors"

const (
	bytesPerInt = 4
	sizeBase    = 7
)

var (
	bSize = make([]byte, bytesPerInt) // Made for reusing in FormSize

	ErrInvalidSizeFormat = errors.New("parsing size: invalid format of tag's/frame's size")
	ErrSizeOverflow      = errors.New("forming size: size of tag/frame is more than allowed in id3 tag")
)

// FormSize transforms int to byte slice with ID3v2 size (4 * 0b0xxxxxxx).
//
// If size more than allowed (256MB), then it returns ErrSizeOverflow.
func FormSize(n int) ([]byte, error) {
	maxN := 268435455 // 0b11111... (28 digits)
	if n > maxN {
		return nil, ErrSizeOverflow
	}

	mask := 1<<sizeBase - 1

	for i := range bSize {
		bSize[len(bSize)-1-i] = byte(n & mask)
		n >>= sizeBase
	}

	return bSize, nil
}

// ParseSize parses byte slice with ID3v2 size (4 * 0b0xxxxxxx) and returns
// parsed int64 number.
//
// If length of slice is more than 4 or if there is invalid size format (e.g.
// one byte in slice is like 0b1xxxxxxx), then it returns ErrInvalidSizeFormat.
func ParseSize(data []byte) (int64, error) {
	if len(data) > bytesPerInt {
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
