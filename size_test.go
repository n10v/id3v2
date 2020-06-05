// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"testing"
)

var (
	synchSafeSizeUint  uint = 15351
	synchSafeSizeBytes      = []byte{0, 0, 119, 119}

	synchUnsafeSizeUint  uint = 65535
	synchUnsafeSizeBytes      = []byte{0, 0, 255, 255}
)

func testWriteSize(sizeUint uint, sizeBytes []byte, synchSafe bool, t *testing.T) {
	t.Parallel()

	buf := new(bytes.Buffer)
	bw := newBufWriter(buf)

	bw.WriteBytesSize(sizeUint, synchSafe)
	if err := bw.Flush(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), sizeBytes) {
		t.Errorf("Expected: %v, got: %v", sizeBytes, buf.Bytes())
	}
}

func testParseSize(sizeUint uint, sizeBytes []byte, synchSafe bool, t *testing.T) {
	t.Parallel()

	size, err := parseSize(sizeBytes, synchSafe)
	if err != nil {
		t.Error(err)
	}
	if size != int64(sizeUint) {
		t.Errorf("Expected: %v, got: %v", sizeUint, size)
	}
}

func TestWriteSynchSafeSize(t *testing.T) {
	testWriteSize(synchSafeSizeUint, synchSafeSizeBytes, true, t)
}

func TestParseSynchSafeSize(t *testing.T) {
	testParseSize(synchSafeSizeUint, synchSafeSizeBytes, true, t)
}

func TestWriteSynchUnsafeSize(t *testing.T) {
	testWriteSize(synchUnsafeSizeUint, synchUnsafeSizeBytes, false, t)
}

func TestParseSynchUnsafeSize(t *testing.T) {
	testParseSize(synchUnsafeSizeUint, synchUnsafeSizeBytes, false, t)
}

func TestParseSynchUnsafeSizeUsingSynchSafeFlag(t *testing.T) {
	t.Parallel()

	_, err := parseSize(synchUnsafeSizeBytes, true)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err != ErrInvalidSizeFormat {
		t.Fatalf("Expected ErrInvalidSizeFormat, got %v", err)
	}

}
