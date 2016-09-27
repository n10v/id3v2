// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

const (
	bytesPerInt = 4
	sizeBase    = 7
)

var (
	byteSize = make([]byte, bytesPerInt) // Made for reusing in FormSize
)

// FormSize transforms uint32 integer to byte slice with
// ID3v2 size (4 * 0b0xxxxxxx).
//
// If size more than allowed (256MB), then panic occurs.
func FormSize(n int64) []byte {
	allowedSize := int64(268435455) // 0b11111... (28 digits)
	if n > allowedSize {
		panic("size of tag/frame is more than allowed in id3 tag")
	}

	mask := int64(1<<sizeBase - 1)

	for i := range byteSize {
		byteSize[len(byteSize)-i-1] = byte(n & mask)
		n >>= sizeBase
	}

	return byteSize
}

// ParseSize parses byte slice with ID3v2 size (4 * 0b0xxxxxxx) and returns
// uint32 integer.
//
// If length of slice is more than 4 or if there is invalid size format (e.g.
// one byte in slice is like 0b1xxxxxxx), then panic occurs.
func ParseSize(data []byte) int64 {
	var size int64

	if len(data) > bytesPerInt {
		panic("invalid length of tag's/frame's size (it must be equal or less than 4)")
	}

	for _, b := range data {
		if b&0x80 > 0 { // 0x80 = 0b1000_0000
			panic("invalid format of tag's/frame's size")
		}

		size = (size << sizeBase) | int64(b)
	}

	return size
}
