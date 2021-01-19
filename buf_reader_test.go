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

// TestReadTextUTF16WithLeadingEmptyString tests if string encoded in UTF16 with BOM
// with leading empty string with same encoding is read correctly.
//
// E.g. this can happen in comment frame with empty description and encoded in UTF16 with BOM.
//
// See https://github.com/bogem/id3v2/issues/53.
func TestReadTextUTF16WithLeadingEmptyString(t *testing.T) {
	t.Parallel()

	sampleText1 := append(bom, EncodingUTF16.TerminationBytes...)

	utf16C := []byte{0x43, 0x00} // "C" char in UTF-16.
	sampleText2 := append(bom, append(utf16C, EncodingUTF16.TerminationBytes...)...)

	sampleText := append(sampleText1, sampleText2...)

	bufReader := newBufReader(bytes.NewReader(sampleText))

	text := decodeText(bufReader.ReadText(EncodingUTF16), EncodingUTF16)
	if text != "" {
		t.Errorf("Expected empty text, got: %v", text)
	}
	// bufReader should only read sampleText1, so Buffered() should return len of sampleText2.
	if bufReader.buf.Buffered() != len(sampleText2) {
		t.Errorf("Expected buffered: %v, got %v", len(sampleText2), bufReader.buf.Buffered())
	}

	text = decodeText(bufReader.ReadText(EncodingUTF16), EncodingUTF16)
	utf8C := "C"
	if text != utf8C {
		t.Errorf("Expected text: %v, got: %v", utf8C, text)
	}
	// bufReader.buf should be empty, because it should read the whole sampleText.
	if bufReader.buf.Buffered() != 0 {
		t.Errorf("Expected buffered: 0, got %v", bufReader.buf.Buffered())
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
