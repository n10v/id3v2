// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"bytes"
	"testing"
)

var (
	sizeInt   uint32 = 15351
	sizeBytes        = []byte{0, 0, 0x77, 0x77}
)

func TestParseSize(t *testing.T) {
	size := ParseSize(sizeBytes)
	if size != sizeInt {
		t.Errorf("Expected: %v, got: %v", sizeInt, size)
	}
}

func TestFormSize(t *testing.T) {
	size := FormSize(sizeInt)
	if !bytes.Equal(sizeBytes, size) {
		t.Errorf("Expected: %v, got: %v", sizeBytes, size)
	}
}
