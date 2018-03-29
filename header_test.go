// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"testing"
)

var (
	th = tagHeader{
		FramesSize: 15351,
		Version:    4,
	}
	thb = []byte{73, 68, 51, 4, 0, 0, 0, 0, 0x77, 0x77}
)

// TestParseHeader checks if parseHeader works right with correct tag header.
func TestParseHeader(t *testing.T) {
	t.Parallel()

	parsed, err := parseHeader(bytes.NewReader(thb))
	if err != nil {
		t.Error(err)
	}
	if parsed != th {
		t.Fatalf("Expected: %v, got: %v", th, parsed)
	}
}

// TestWriteTagHeader checks if writeTagHeader works right with correct tag header.
func TestWriteTagHeader(t *testing.T) {
	t.Parallel()

	buf := new(bytes.Buffer)
	bw := newBufWriter(buf)

	writeTagHeader(bw, 15351, 4)
	if err := bw.Flush(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(thb, buf.Bytes()) {
		t.Fatalf("Expected %v, got %v", thb, buf.Bytes())
	}
}

// TestSmallTagHeader checks if parseHeader returns an error
// when size of reader is smaller than tagHeaderSize.
func TestSmallTagHeader(t *testing.T) {
	t.Parallel()

	_, err := parseHeader(bytes.NewReader([]byte{0, 0, 0}))
	if err != ErrSmallHeaderSize {
		t.Fatalf("Expected err contains %q, got %q", "less than expected", err)
	}
}

// TestIsNotID3 checks if parseHeader returns correct error
// when there is no "ID3" literal at the beginning.
func TestIsNotID3(t *testing.T) {
	t.Parallel()

	_, err := parseHeader(bytes.NewReader([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}))
	if err != errNoTag {
		t.Fatalf("Expected: %q, got: %q", errNoTag, err)
	}
}
