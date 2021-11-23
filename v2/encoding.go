package id3v2

import (
	"bytes"
	"io/ioutil"

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

// bom is used in UTF-16 encoded Unicode with BOM.
// See https://en.wikipedia.org/wiki/Byte_order_mark.
var bom = []byte{0xFF, 0xFE}

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

	bw := getBufWriter(ioutil.Discard)
	defer putBufWriter(bw)

	encodeWriteText(bw, src, enc)

	return bw.Written()

}

// decodeText decodes src from "from" encoding to UTF-8.
func decodeText(src []byte, from Encoding) string {
	src = bytes.TrimSuffix(src, from.TerminationBytes) // See https://github.com/bogem/id3v2/issues/41

	if from.Equals(EncodingUTF8) {
		return string(src)
	}

	// If src is just BOM, then it's an empty string.
	if from.Equals(EncodingUTF16) && bytes.Equal(src, bom) {
		return ""
	}

	fromXEncoding := resolveXEncoding(src, from)
	result, err := fromXEncoding.Decoder().Bytes(src)
	if err != nil {
		return string(src)
	}

	// HACK: Delete REPLACEMENT CHARACTER (�) if encoding went wrong.
	// See https://apps.timwhitlock.info/unicode/inspect?s=%EF%BF%BD.
	// See https://en.wikipedia.org/wiki/Byte_order_mark#UTF-8.
	if from.Equals(EncodingUTF16) {
		// bytes.Replace(s, old, new, -1) is the same as bytes.ReplaceAll(s, old, new),
		// but bytes.ReplaceAll is only added in Go 1.12.
		result = bytes.Replace(result, []byte{0xEF, 0xBF, 0xBD}, []byte{}, -1)
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

	if to.Equals(EncodingUTF16) && !bytes.HasSuffix([]byte(encoded), []byte{0}) {
		bw.WriteByte(0)
	}

	return nil
}

func resolveXEncoding(src []byte, encoding Encoding) xencodingWrapper {
	switch encoding.Key {
	case 0:
		return xencodingISO
	case 1:
		if len(src) > 2 && bytes.Equal(src[:2], bom) {
			return xencodingUTF16LEBOM
		}
		return xencodingUTF16BEBOM
	case 2:
		return xencodingUTF16BE
	}

	return xencodingUTF8
}
