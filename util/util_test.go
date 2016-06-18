// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"github.com/bogem/id3v2/testutil"
	"testing"
)

var (
	sizeInt   uint32 = 15351
	sizeBytes        = []byte{0, 0, 0x77, 0x77}
)

func TestParseSize(t *testing.T) {
	size, err := ParseSize(sizeBytes)
	if size != sizeInt {
		t.Errorf("Expected: %v, got: %v", sizeInt, size)
	}
	if err != nil {
		t.Fail()
	}
}

func TestFormSize(t *testing.T) {
	size, err := FormSize(sizeInt)
	if err != nil {
		t.Fail()
	}
	if err := testutil.AreByteSlicesEqual(sizeBytes, size); err != nil {
		t.Errorf("Expected: %v, got: %v", sizeBytes, size)
	}
}
