package id3v2

import (
	"bytes"
	"github.com/bogem/id3v2/testutil"
	"testing"
)

var (
	th = TagHeader{
		FramesSize: 77,
		Version:    4,
	}
	thb = []byte{73, 68, 51, 4, 0, 0, 0, 0, 0, 77}
)

func TestParseHeader(t *testing.T) {
	parsed, err := ParseHeader(bytes.NewReader(thb))
	if err != nil {
		t.Error(err)
	}
	if *parsed != th {
		t.Fail()
	}
}

func TestIsID3Tag(t *testing.T) {
	if !isID3Tag([]byte(ID3Identifier)) {
		t.Fail()
	}
}

func TestFormTagHeader(t *testing.T) {
	formed, err := FormTagHeader(th)
	if err != nil {
		t.Fail()
	}

	if err := testutil.AreByteSlicesEqual(formed, thb); err != nil {
		t.Errorf("Expected: %v, got: %v", thb, formed)
	}
}
