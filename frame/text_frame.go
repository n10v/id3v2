// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package frame

import (
	"github.com/bogem/id3v2/bytesbufferpool"
	"github.com/bogem/id3v2/util"
)

// TextFramer is used to work with all text frames
// (all T*** frames like TIT2, TALB and so on).
type TextFramer interface {
	Framer

	Encoding() util.Encoding
	SetEncoding(util.Encoding)

	Text() string
	SetText(string)
}

type TextFrame struct {
	encoding util.Encoding
	text     string
}

func (tf TextFrame) Bytes() ([]byte, error) {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.WriteByte(tf.encoding.Key)
	b.WriteString(tf.text)

	return b.Bytes(), nil
}

func (tf TextFrame) Encoding() util.Encoding {
	return tf.encoding
}

func (tf *TextFrame) SetEncoding(e util.Encoding) {
	tf.encoding = e
}

func (tf TextFrame) Text() string {
	return tf.text
}

func (tf *TextFrame) SetText(text string) {
	tf.text = text
}
