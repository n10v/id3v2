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
func FormSize(n uint32) []byte {
	allowedSize := uint32(268435455) // 4 * 0b01111111
	if n > allowedSize {
		panic("Size is more than allowed in id3 tag")
	}

	mask := uint32(1<<sizeBase - 1)

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
func ParseSize(data []byte) uint32 {
	var size uint32

	if len(data) > bytesPerInt {
		panic("Invalid data length (it must be equal or less than 4)")
	}

	for _, b := range data {
		if b&0x80 > 0 { // 0x80 = 0b1000_0000
			panic("Invalid size format")
		}

		size = (size << sizeBase) | uint32(b)
	}

	return size
}
