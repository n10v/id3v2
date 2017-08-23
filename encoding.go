package id3v2

import (
	"bufio"

	xencoding "golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

// Encoding is a struct for encodings.
type Encoding struct {
	Name             string
	Key              byte
	TerminationBytes []byte
}

func (e Encoding) Equals(other Encoding) bool {
	return e.Key == other.Key
}

func (e Encoding) String() string {
	return e.Name
}

// Available encodings.
var (
	// EncodingISO is ISO-8859-1 encoding.
	EncodingISO = Encoding{
		Name:             "ISO-8859-1",
		Key:              0,
		TerminationBytes: []byte{0},
	}

	// EncodingUTF16 is UTF-16 encoded Unicode with BOM.
	EncodingUTF16 = Encoding{
		Name:             "UTF-16 encoded Unicode with BOM",
		Key:              1,
		TerminationBytes: []byte{0, 0},
	}

	// EncodingUTF16BE is UTF-16BE encoded Unicode without BOM.
	EncodingUTF16BE = Encoding{
		Name:             "UTF-16BE encoded Unicode without BOM",
		Key:              2,
		TerminationBytes: []byte{0, 0},
	}

	// EncodingUTF8 is UTF-8 encoded Unicode.
	EncodingUTF8 = Encoding{
		Name:             "UTF-8 encoded Unicode",
		Key:              3,
		TerminationBytes: []byte{0},
	}

	encodings = []Encoding{EncodingISO, EncodingUTF16, EncodingUTF16BE, EncodingUTF8}
)

// getEncoding returns Encoding in accordance with ID3v2 key.
func getEncoding(key byte) Encoding {
	if key > 3 {
		return EncodingUTF8
	}
	return encodings[key]
}

// encodedSize counts length of UTF-8 src if it's encoded to enc.
func encodedSize(src string, enc Encoding) int {
	if enc.Equals(EncodingUTF8) {
		return len(src)
	}

	toXEncoding := resolveXEncoding(nil, enc)
	encoded, _ := toXEncoding.NewEncoder().String(src)
	return len(encoded)
}

// decodeText decodes src from "from" encoding to UTF-8.
func decodeText(src []byte, from Encoding) string {
	if from.Equals(EncodingUTF8) {
		return string(src)
	}

	fromXEncoding := resolveXEncoding(src, from)
	dst, err := fromXEncoding.NewDecoder().Bytes(src)
	if err != nil {
		return string(src)
	}

	return string(dst)
}

// encodeWriteText encodes src from UTF-8 to "to" encoding and writes to bw.
func encodeWriteText(bw *bufio.Writer, src string, to Encoding) (n int, err error) {
	if to.Equals(EncodingUTF8) {
		return bw.WriteString(src)
	}

	toXEncoding := resolveXEncoding(nil, to)
	encoded, err := toXEncoding.NewEncoder().String(src)
	if err != nil {
		return 0, err
	}
	return bw.WriteString(encoded)
}

// resolveXEncoding returns golang.org/x/text/encoding encoding
// from src and encoding.
func resolveXEncoding(src []byte, encoding Encoding) xencoding.Encoding {
	switch encoding.Key {
	case 0: // ISO-8859-1
		return charmap.ISO8859_1
	case 1: // UTF-16 With BOM
		if len(src) < 2 {
			return unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
		} else if src[0] == 0xFE && src[1] == 0xFF {
			return unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
		} else if src[0] == 0xFF && src[1] == 0xFE {
			return unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
		}
		return unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
	case 2: // UTF-16BE without BOM
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case 3: // UTF-8
		return unicode.UTF8
	}

	return unicode.UTF8
}
