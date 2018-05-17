package id3v2

import (
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

// xencodingWrapper is a struct that stores decoder and encoder for
// appropriate x/text/encoding. It's used to reduce allocations
// through creating decoder and encoder only one time and storing it.
type xencodingWrapper struct {
	decoder *xencoding.Decoder
	encoder *xencoding.Encoder
}

func newXEncodingWrapper(e xencoding.Encoding) xencodingWrapper {
	return xencodingWrapper{
		decoder: e.NewDecoder(),
		encoder: e.NewEncoder(),
	}
}

func (e *xencodingWrapper) Decoder() *xencoding.Decoder {
	return e.decoder
}

func (e *xencodingWrapper) Encoder() *xencoding.Encoder {
	return e.encoder
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

	xencodingISO        = newXEncodingWrapper(charmap.ISO8859_1)
	xencodingUTF16BEBOM = newXEncodingWrapper(unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM))
	xencodingUTF16LEBOM = newXEncodingWrapper(unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM))
	xencodingUTF16BE    = newXEncodingWrapper(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM))
	xencodingUTF8       = newXEncodingWrapper(unicode.UTF8)
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
	encoded, _ := toXEncoding.Encoder().String(src)
	return len(encoded)
}

// decodeText decodes src from "from" encoding to UTF-8.
func decodeText(src []byte, from Encoding) string {
	if from.Equals(EncodingUTF8) {
		return string(src)
	}

	fromXEncoding := resolveXEncoding(src, from)
	result, err := fromXEncoding.Decoder().Bytes(src)
	if err != nil {
		return string(src)
	}
	return string(result)
}

// encodeWriteText encodes src from UTF-8 to "to" encoding and writes to bw.
func encodeWriteText(bw *bufWriter, src string, to Encoding) error {
	if to.Equals(EncodingUTF8) {
		bw.WriteString(src)
		return nil
	}

	toXEncoding := resolveXEncoding(nil, to)
	encoded, err := toXEncoding.Encoder().String(src)
	if err != nil {
		return err
	}
	bw.WriteString(encoded)
	return nil
}

func resolveXEncoding(src []byte, encoding Encoding) xencodingWrapper {
	switch encoding.Key {
	case 0:
		return xencodingISO
	case 1:
		if len(src) > 2 && src[0] == 0xFF && src[1] == 0xFE {
			return xencodingUTF16LEBOM
		}
		return xencodingUTF16BEBOM
	case 2:
		return xencodingUTF16BE
	}

	return xencodingUTF8
}
