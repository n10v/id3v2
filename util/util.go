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
	byteSize = make([]byte, bytesPerInt) // Made for reusing in FormSize

	InvalidSizeFormat = errors.New("parsing size: invalid format of tag's/frame's size")
	SizeOverflow      = errors.New("forming size: size of tag/frame is more than allowed in id3 tag")
)

// FormSize transforms int to byte slice with ID3v2 size (4 * 0b0xxxxxxx).
//
// If size more than allowed (256MB), then method returns SizeOverflow.
func FormSize(n int) ([]byte, error) {
	allowedSize := 268435455 // 0b11111... (28 digits)
	if n > allowedSize {
		return nil, SizeOverflow
	}

	mask := 1<<sizeBase - 1

	for i := range byteSize {
		byteSize[len(byteSize)-i-1] = byte(n & mask)
		n >>= sizeBase
	}

	return byteSize, nil
}

// ParseSize parses byte slice with ID3v2 size (4 * 0b0xxxxxxx) and returns
// int64.
//
// If length of slice is more than 4 or if there is invalid size format (e.g.
// one byte in slice is like 0b1xxxxxxx), then method return InvalidSizeFormat.
func ParseSize(data []byte) (int64, error) {
	var size int64

	if len(data) > bytesPerInt {
		return 0, InvalidSizeFormat
	}

	for _, b := range data {
		if b&0x80 > 0 { // 0x80 = 0b1000_0000
			return 0, InvalidSizeFormat
		}

		size = (size << sizeBase) | int64(b)
	}

	return size, nil
}
