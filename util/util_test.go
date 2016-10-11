// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"bytes"
	"testing"
)

var (
	sizeInt   int64 = 15351
	sizeBytes       = []byte{0, 0, 0x77, 0x77}
)

func TestParseSize(t *testing.T) {
	size, err := ParseSize(sizeBytes)
	if err != nil {
		t.Error(err)
	}
	if size != sizeInt {
		t.Errorf("Expected: %v, got: %v", sizeInt, size)
	}
}

func TestFormSize(t *testing.T) {
	size, err := FormSize(sizeInt)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(sizeBytes, size) {
		t.Errorf("Expected: %v, got: %v", sizeBytes, size)
	}
}
