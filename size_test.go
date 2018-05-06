// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"testing"
)

var (
	sizeUint  uint = 15351
	sizeBytes      = []byte{0, 0, 0x77, 0x77}
)

func TestWriteBytesSize(t *testing.T) {
	t.Parallel()

	buf := new(bytes.Buffer)
	bw := newBufWriter(buf)

	bw.WriteBytesSize(sizeUint)
	if err := bw.Flush(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), sizeBytes) {
		t.Errorf("Expected: %v, got: %v", sizeBytes, buf.Bytes())
	}
}

func TestParseSize(t *testing.T) {
	t.Parallel()

	size, err := parseSize(sizeBytes)
	if err != nil {
		t.Error(err)
	}
	if size != int64(sizeUint) {
		t.Errorf("Expected: %v, got: %v", sizeUint, size)
	}
}
