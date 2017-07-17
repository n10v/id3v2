// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"bufio"
	"bytes"
	"testing"
)

var (
	sizeInt   = 15351
	sizeBytes = []byte{0, 0, 0x77, 0x77}
)

func TestWriteBytesSize(t *testing.T) {
	t.Parallel()

	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)

	if err := WriteBytesSize(bw, sizeInt); err != nil {
		t.Error(err)
	}
	bw.Flush()
	if !bytes.Equal(buf.Bytes(), sizeBytes) {
		t.Errorf("Expected: %v, got: %v", sizeBytes, buf.Bytes())
	}
}

func TestParseSize(t *testing.T) {
	t.Parallel()

	size, err := ParseSize(sizeBytes)
	if err != nil {
		t.Error(err)
	}
	if size != int64(sizeInt) {
		t.Errorf("Expected: %v, got: %v", sizeInt, size)
	}
}

func BenchmarkWriteBytesSize(b *testing.B) {
	bw := bufio.NewWriter(new(bytes.Buffer))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := WriteBytesSize(bw, 268435454); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := ParseSize([]byte{127, 127, 127, 127}); err != nil {
			b.Fatal(err)
		}
	}
}
