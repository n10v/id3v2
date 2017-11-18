// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import "io"

// TextFrame is used to work with all text frames
// (all T*** frames like TIT2 (title), TALB (album) and so on).
type TextFrame struct {
	Encoding Encoding
	Text     string
}

func (tf TextFrame) Size() int {
	return 1 + encodedSize(tf.Text, tf.Encoding)
}

func (tf TextFrame) WriteTo(w io.Writer) (n int64, err error) {
	bw, ok := resolveBufioWriter(w)
	if !ok {
		defer putBufioWriter(bw)
	}

	var nn int

	bw.WriteByte(tf.Encoding.Key)
	n += 1

	nn, err = encodeWriteText(bw, tf.Text, tf.Encoding)
	n += int64(nn)
	if err != nil {
		return
	}

	return n, bw.Flush()
}

func parseTextFrame(rd io.Reader) (Framer, error) {
	bufRd := getUtilReader(rd)
	defer putUtilReader(bufRd)

	encodingKey, err := bufRd.ReadByte()
	if err != nil {
		return nil, err
	}
	encoding := getEncoding(encodingKey)

	buf := getBytesBuffer()
	defer putBytesBuffer(buf)

	if _, err := buf.ReadFrom(bufRd); err != nil {
		return nil, err
	}

	tf := TextFrame{
		Encoding: encoding,
		Text:     decodeText(buf.Bytes(), encoding),
	}

	return tf, nil
}
