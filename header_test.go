// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"bytes"
	"testing"
)

var (
	th = tagHeader{
		FramesSize: 0,
		Version:    4,
	}
	thb = []byte{73, 68, 51, 4, 0, 0, 0, 0, 0, 0}
)

func TestParseHeader(t *testing.T) {
	parsed, err := parseHeader(bytes.NewReader(thb))
	if err != nil {
		t.Error(err)
	}
	if parsed != th {
		t.Fail()
	}
}

func TestWriteTagHeader(t *testing.T) {
	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)
	if err := writeTagHeader(bw, []byte{0, 0, 0, 0}, 4); err != nil {
		t.Fatal(err)
	}
	if err := bw.Flush(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(thb, buf.Bytes()) {
		t.Fatalf("Expected %v, got %v", thb, buf.Bytes())
	}
}
