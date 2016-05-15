package testutil

import (
	"errors"
	"strconv"
)

func AreByteSlicesEqual(expected []byte, got []byte) error {
	if len(expected) != len(got) {
		return errors.New("Slices have different lengths")
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != got[i] {
			return errors.New("At " + strconv.Itoa(i) + " byte: expected " + strconv.Itoa(int(expected[i])) + ", got " + strconv.Itoa(int(got[i])))
		}
	}

	return nil
}
