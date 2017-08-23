package id3v2

import (
	"bufio"
	"bytes"
	"testing"
)

func TestDecodeText(t *testing.T) {
	testCases := []struct {
		src  []byte
		from Encoding
		utf8 string
	}{
		{[]byte{0x48, 0xE9, 0x6C, 0x6C, 0xF6}, EncodingISO, "Héllö"},
		{[]byte{0xFF, 0xFE, 0x48, 0x00, 0xE9, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0xF6, 0x00}, EncodingUTF16, "Héllö"}, // UTF-16LE with BOM
		{[]byte{0x00, 0x48, 0x00, 0xE9, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0xF6}, EncodingUTF16BE, "Héllö"},
	}

	for _, tc := range testCases {
		got := decodeText(tc.src, tc.from)
		if got != tc.utf8 {
			t.Errorf("Expected %q from %v encoding, got %q", tc.utf8, tc.from, got)
		}
	}
}

func TestEncodeWriteText(t *testing.T) {
	testCases := []struct {
		src      string
		to       Encoding
		expected []byte
		size     int
	}{
		{"Héllö", EncodingISO, []byte{0x48, 0xE9, 0x6C, 0x6C, 0xF6}, 5},
		{"Héllö", EncodingUTF16, []byte{0xFE, 0xFF, 0x00, 0x48, 0x00, 0xE9, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0xF6}, 12},
		{"Héllö", EncodingUTF16BE, []byte{0x00, 0x48, 0x00, 0xE9, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0xF6}, 10},
	}

	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)

	for _, tc := range testCases {
		buf.Reset()
		n, err := encodeWriteText(bw, tc.src, tc.to)
		if err != nil {
			t.Errorf("Error by encoding and writing text: %v", err)
		}

		bw.Flush()
		got := buf.Bytes()
		if !bytes.Equal(got, tc.expected) {
			t.Errorf("Expected %q from %q encoding, got %q", tc.expected, tc.to, got)
		}
		if n != tc.size {
			t.Errorf("Expected %v size, got %v", tc.size, n)
		}
	}
}
