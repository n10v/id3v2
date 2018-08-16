package id3v2

import (
	"io"
)

// UserDefinedTextFrame is used to work with TXXX frames.
// There can be many UserDefinedTextFrames but the Desciption fields need to be unique.
type UserDefinedTextFrame struct {
	Encoding    Encoding
	Description string
	Value       string
}

func (uf UserDefinedTextFrame) Size() int {
	return 1 + encodedSize(uf.Description, uf.Encoding) + len(uf.Encoding.TerminationBytes) + encodedSize(uf.Value, uf.Encoding)
}

func (uf UserDefinedTextFrame) WriteTo(w io.Writer) (n int64, err error) {
	return useBufWriter(w, func(bw *bufWriter) {
		bw.WriteByte(uf.Encoding.Key)
		bw.EncodeAndWriteText(uf.Description, uf.Encoding)
		bw.Write(uf.Encoding.TerminationBytes)
		bw.EncodeAndWriteText(uf.Value, uf.Encoding)
	})
}

func parseUserDefinedTextFrame(br *bufReader) (Framer, error) {
	encodingKey := br.ReadByte()
	encoding := getEncoding(encodingKey)

	description := br.ReadTillDelims(encoding.TerminationBytes)
	br.Discard(len(encoding.TerminationBytes))

	if br.Err() != nil {
		return nil, br.Err()
	}

	value := getBytesBuffer()
	defer putBytesBuffer(value)

	if _, err := value.ReadFrom(br); err != nil {
		return nil, err
	}

	uf := UserDefinedTextFrame{
		Encoding:    encoding,
		Description: decodeText(description, encoding),
		Value:       decodeText(value.Bytes(), encoding),
	}

	return uf, nil
}
