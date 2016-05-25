package testutil

import (
	"errors"
	"fmt"
	"strconv"
)

func AreByteSlicesEqual(expected []byte, got []byte) error {
	if len(expected) != len(got) {
		return errors.New("Slices have different lengths")
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != got[i] {
			return errors.New("At " + strconv.Itoa(i) + " byte: expected " + fmt.Sprint(expected[i:]) + ", got " + fmt.Sprint(got[i:]))
		}
	}

	return nil
}
