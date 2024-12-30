// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import "io"

type LinkFrame struct {
	Encoding Encoding
	Url      string
}

func (lf LinkFrame) Size() int {
	return 1 + encodedSize(lf.Url, lf.Encoding) + len(lf.Encoding.TerminationBytes)
}

func (lf LinkFrame) UniqueIdentifier() string {
	return "ID"
}

func (lf LinkFrame) WriteTo(w io.Writer) (int64, error) {
	return useBufWriter(w, func(bw *bufWriter) {
		bw.WriteByte(lf.Encoding.Key)
		bw.EncodeAndWriteText(lf.Url, lf.Encoding)

		// https://github.com/bogem/id3v2/pull/52
		// https://github.com/bogem/id3v2/pull/33
		bw.Write(lf.Encoding.TerminationBytes)
	})
}

func parseLinkFrame(br *bufReader) (Framer, error) {
	encoding := getEncoding(br.ReadByte())

	if br.Err() != nil {
		return nil, br.Err()
	}

	buf := getBytesBuffer()
	defer putBytesBuffer(buf)
	if _, err := buf.ReadFrom(br); err != nil {
		return nil, err
	}

	lf := LinkFrame{
		Encoding: encoding,
		Url:      decodeText(buf.Bytes(), encoding),
	}

	return lf, nil
}
