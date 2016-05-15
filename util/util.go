package util

import (
	"errors"
)

const (
	BytesPerInt = 4
	SizeBase    = 7

	NativeEncoding = 3 // UTF-8
)

func ParseSize(data []byte) (size uint32, err error) {
	if len(data) > BytesPerInt {
		err = errors.New("Invalid data length (it must be equal or less than 4)")
		return
	}

	for _, b := range data {
		if b&0x80 > 0 { // 0x80 = 0b1000_0000
			err = errors.New("Invalid size format")
			return
		}

		size = (size << SizeBase) | uint32(b)
	}

	return
}

func FormSize(n uint32) ([]byte, error) {
	if n > 268435455 { //4 * 0b01111111
		return nil, errors.New("Size is more than allowed in id3 tag")
	}

	mask := uint32(1<<SizeBase - 1)
	bytes := make([]byte, BytesPerInt)

	for i, _ := range bytes {
		bytes[len(bytes)-i-1] = byte(n & mask)
		n >>= SizeBase
	}

	return bytes, nil
}
