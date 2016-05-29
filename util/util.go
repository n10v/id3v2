package util

import (
	"errors"
)

const (
	BytesPerInt = 4
	SizeBase    = 7

	NativeEncoding = 3 // UTF-8
)

var (
	byteSize = make([]byte, BytesPerInt) // Made for reusing in FormSize
)

func ParseSize(data []byte) (uint32, error) {
	var size uint32

	if len(data) > BytesPerInt {
		err := errors.New("Invalid data length (it must be equal or less than 4)")
		return 0, err
	}

	for _, b := range data {
		if b&0x80 > 0 { // 0x80 = 0b1000_0000
			err := errors.New("Invalid size format")
			return 0, err
		}

		size = (size << SizeBase) | uint32(b)
	}

	return size, nil
}

func FormSize(n uint32) ([]byte, error) {
	allowedSize := uint32(268435455) // 4 * 0b01111111
	if n > allowedSize {
		return nil, errors.New("Size is more than allowed in id3 tag")
	}

	mask := uint32(1<<SizeBase - 1)

	for i, _ := range byteSize {
		byteSize[len(byteSize)-i-1] = byte(n & mask)
		n >>= SizeBase
	}

	return byteSize, nil
}
