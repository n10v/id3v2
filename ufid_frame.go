package id3v2

import (
	"io"
)

type UFIDFrame struct {
	OwnerIdentifier string
	Identifier      []byte
}

func (uf UFIDFrame) Size() int {
	return encodedSize(uf.OwnerIdentifier, EncodingISO) + len(EncodingISO.TerminationBytes) + len(uf.Identifier)
}

func (uf UFIDFrame) WriteTo(w io.Writer) (n int64, err error) {
	return useBufWriter(w, func(bw *bufWriter) {
		bw.WriteString(uf.OwnerIdentifier)
		bw.Write(EncodingISO.TerminationBytes)
		bw.Write(uf.Identifier)
	})
}

func parseUFIDFrame(br *bufReader) (Framer, error) {
	owner := br.ReadTillDelims(EncodingISO.TerminationBytes)
	br.Discard(len(EncodingISO.TerminationBytes))

	if br.Err() != nil {
		return nil, br.Err()
	}

	ident := br.ReadAll()

	uf := UFIDFrame{
		OwnerIdentifier: decodeText(owner, EncodingISO),
		Identifier:      ident,
	}

	return uf, nil
}
