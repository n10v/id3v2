// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"io"
	"io/ioutil"

	"github.com/bogem/id3v2/bytesbufferpool"
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
//	tag.AddFrame(tag.ID("Mood"), textFrame)
type TextFrame struct {
	Encoding util.Encoding
	Text     string
}

func (tf TextFrame) Body() []byte {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.WriteByte(tf.Encoding.Key)
	b.WriteString(tf.Text)

	return b.Bytes()
}

func ParseTextFrame(rd io.Reader) (Framer, error) {
	bufRd := bufio.NewReader(rd)

	encoding, err := bufRd.ReadByte()
	if err != nil {
		return nil, err
	}

	text, err := ioutil.ReadAll(bufRd)
	if err != nil {
		return nil, err
	}

	tf := TextFrame{
		Encoding: Encodings[encoding],
		Text:     string(text),
	}

	return tf, nil
}
