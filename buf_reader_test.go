// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"io"
	"testing"
)

var bs = []byte{0, 11, 22, 33, 44, 55, 77, 88, 55, 55, 66, 77, 88}

func TestReadTillDelim(t *testing.T) {
	t.Parallel()

	bufReader := newBufReader(bytes.NewReader(bs))

	firstIndexOf55 := bytes.Index(bs, []byte{55})
	if firstIndexOf55 < 0 {
		t.Fatal("Can't find 55 in bs")
	}
	expected := bs[:firstIndexOf55]

	read, err := bufReader.readTillDelim(55)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expected, read) {
		t.Errorf("Expected: %v, got: %v", expected, read)
	}
	if len(bs)-len(expected) != bufReader.buf.Buffered() {
		t.Errorf("Expected buffered: %v, got: %v", len(bs)-len(expected), bufReader.buf.Buffered())
	}
}

func TestReadTillZero(t *testing.T) {
	t.Parallel()

	bufReader := newBufReader(bytes.NewReader(bs))

	read, err := bufReader.readTillDelim(0)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte{}, read) {
		t.Errorf("Expected empty buffer, got %v", read)
	}
	if len(bs) != bufReader.buf.Buffered() {
		t.Errorf("Expected buffered: %v, got: %v", len(bs), bufReader.buf.Buffered())
	}
}

func TestNext(t *testing.T) {
	t.Parallel()

	bufReader := newBufReader(bytes.NewReader(bs))
	n := 5 // Read 5 elements

	read := bufReader.Next(n)
	if bufReader.Err() != nil {
		t.Fatal(bufReader.Err())
	}
	if n != len(read) {
		t.Errorf("Expected: %v, got: %v", n, read)
	}
	if !bytes.Equal(bs[:n], read) {
		t.Errorf("Expected: %v, got: %v", bs[:n], len(read))
	}
	if len(bs)-n != bufReader.buf.Buffered() {
		t.Errorf("Expected buffered: %v, got: %v", len(bs)-n, bufReader.buf.Buffered())
	}
}

func TestReadTillDelimEOF(t *testing.T) {
	t.Parallel()

	bufReader := newBufReader(bytes.NewReader(bs))
	_, err := bufReader.readTillDelim(234)
	if err != io.EOF {
		t.Errorf("Expected io.EOF, got %v", err)
	}
}

func TestReadTillDelims(t *testing.T) {
	t.Parallel()

	bufReader := newBufReader(bytes.NewReader(bs))
	delims := []byte{55, 66}
	expectedLen := 9

	read, err := bufReader.readTillDelims(delims)
	if err != nil {
		t.Fatal(err)
	}
	if expectedLen != len(read) {
		t.Errorf("Expected: %v, got: %v", expectedLen, len(read))
	}
	if !bytes.Equal(bs[:expectedLen], read) {
		t.Errorf("Expected: %v, got: %v", bs[:expectedLen], read)
	}
	if len(bs)-len(read) != bufReader.buf.Buffered() {
		t.Errorf("Expected buffered: %v, got: %v", len(bs)-len(read), bufReader.buf.Buffered())
	}
}
