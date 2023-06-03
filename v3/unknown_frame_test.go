// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"testing"
)

func TestUnknownFramesUniqueIdentifiers(t *testing.T) {
	uf1, _ := parseUnknownFrame(newBufReader(new(bytes.Buffer)))
	uf2, _ := parseUnknownFrame(newBufReader(new(bytes.Buffer)))

	if uf1.UniqueIdentifier() == uf2.UniqueIdentifier() {
		t.Errorf("Two unknown frame have same unique identifiers, but every unknown frame should have completely unique identifiers.")
	}
}
