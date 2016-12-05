// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"bytes"
	"io"

	"github.com/bogem/id3v2/rdpool"
	"github.com/bogem/id3v2/util"
)

// TextFrame is used to work with all text frames
// (all T*** frames like TIT2 (title), TALB (album) and so on).
//
// Example of setting a new text frame to existing tag:
//
//	textFrame := id3v2.TextFrame{
//		Encoding: id3v2.ENUTF8,
//		Text:     "Happy",
//	}
//	tag.AddFrame(tag.CommonID("Mood"), textFrame)
type TextFrame struct {
	Encoding util.Encoding
	Text     string
}

func (tf TextFrame) Body() []byte {
	b := new(bytes.Buffer)

	b.WriteByte(tf.Encoding.Key)
	b.WriteString(tf.Text)

	return b.Bytes()
}

func (tf TextFrame) Size() int {
	return 1 + len(tf.Text)
}

func (tf TextFrame) WriteTo(w io.Writer) (n int, err error) {
	var i int
	bw := bufio.NewWriter(w)

	err = bw.WriteByte(tf.Encoding.Key)
	if err != nil {
		return
	}
	n += 1

	i, err = bw.WriteString(tf.Text)
	if err != nil {
		return
	}
	n += i

	err = bw.Flush()
	return
}

func parseTextFrame(rd io.Reader) (Framer, error) {
	bufRd := rdpool.Get(rd)
	defer rdpool.Put(bufRd)

	encoding, err := bufRd.ReadByte()
	if err != nil {
		return nil, err
	}

	text, err := bufRd.ReadAll()
	if err != nil {
		return nil, err
	}

	tf := TextFrame{
		Encoding: Encodings[encoding],
		Text:     string(text),
	}

	return tf, nil
}
