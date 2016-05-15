package testutil

import (
	"testing"
)

type ByteSlicesTestData struct {
	A        []byte
	B        []byte
	Expected bool
}

func TestByteSlicesEquality(t *testing.T) {
	ByteSlicesTests := []ByteSlicesTestData{
		{[]byte{1, 2}, []byte{1, 2, 3}, false},
		{[]byte{1, 2}, []byte{1, 3}, false},
		{[]byte{1, 2}, []byte{1, 2}, true},
	}

	for _, testData := range ByteSlicesTests {
		got := AreByteSlicesEqual(testData.A, testData.B)
		if got != testData.Expected {
			t.Errorf("For %v and %v expected %v, got %v", testData.A, testData.B, testData.Expected, got)
		}

	}
}
