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
	if *parsed != th {
		t.Fail()
	}
}

func TestFormTagHeader(t *testing.T) {
	formed := formTagHeader([]byte{0, 0, 0, 0})
	if !bytes.Equal(formed, thb) {
		t.Fail()
	}
}
