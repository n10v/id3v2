package id3v2

import "io"

type UFIDFrame struct {
	OwnerIdentifier string
	Identifier      []byte
}

func (ufid UFIDFrame) Size() int {
	return encodedSize(ufid.OwnerIdentifier, EncodingISO) + len(EncodingISO.TerminationBytes) + len(ufid.Identifier)
}

func (ufid UFIDFrame) WriteTo(w io.Writer) (n int64, err error) {
	return useBufWriter(w, func(bw *bufWriter) {
		bw.WriteString(ufid.OwnerIdentifier)
		bw.Write(EncodingISO.TerminationBytes)
		bw.Write(ufid.Identifier)
	})
}

func parseUFIDFrame(br *bufReader) (Framer, error) {
	owner := br.ReadTillDelims(EncodingISO.TerminationBytes)
	br.Discard(len(EncodingISO.TerminationBytes))

	if br.Err() != nil {
		return nil, br.Err()
	}

	ident := br.ReadAll()

	ufid := UFIDFrame{
		OwnerIdentifier: decodeText(owner, EncodingISO),
		Identifier:      ident,
	}

	return ufid, nil
}
