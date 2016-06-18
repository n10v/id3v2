package util

import (
	"errors"
)

const (
	bytesPerInt = 4
	sizeBase    = 7
)

var (
	byteSize = make([]byte, bytesPerInt) // Made for reusing in FormSize
)

// ParseSize parses byte slice with ID3v2 size (4 * 0b0xxxxxxx) and returns
// uint32 integer.
//
// If length of slice is more than 4 or if there is invalid size format (e.g.
// one byte in slice is like 0b1xxxxxxx), then error occurs.
func ParseSize(data []byte) (uint32, error) {
	var size uint32

	if len(data) > bytesPerInt {
		err := errors.New("Invalid data length (it must be equal or less than 4)")
		return 0, err
	}

	for _, b := range data {
		if b&0x80 > 0 { // 0x80 = 0b1000_0000
			err := errors.New("Invalid size format")
			return 0, err
		}

		size = (size << sizeBase) | uint32(b)
	}

	return size, nil
}

// FormSize transforms uint32 integer to byte slice with
// ID3v2 size (4 * 0b0xxxxxxx).
//
// If size more than allowed (256MB), then error occurs.
func FormSize(n uint32) ([]byte, error) {
	allowedSize := uint32(268435455) // 4 * 0b01111111
	if n > allowedSize {
		return nil, errors.New("Size is more than allowed in id3 tag")
	}

	mask := uint32(1<<sizeBase - 1)

	for i, _ := range byteSize {
		byteSize[len(byteSize)-i-1] = byte(n & mask)
		n >>= sizeBase
	}

	return byteSize, nil
}
