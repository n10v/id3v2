// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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
		err := AreByteSlicesEqual(testData.A, testData.B)
		if (err != nil) && (testData.Expected != false) {
			t.Error(err)
		}
	}
}
