// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"github.com/bogem/id3v2/bytesbufferpool"
	"github.com/bogem/id3v2/util"
)

// TextFrame is used to work with all text frames
// (all T*** frames like TIT2 (title), TALB (album) and so on).
//
// Example of setting a new text frame to existing tag:
//  textFrame := id3v2.TextFrame{
//    Encoding: id3v2.ENUTF8,
//    Text:     "Happy",
//  }
//  id := "TMOO" // Mood frame ID
//  tag.AddFrame(id, textFrame)
type TextFrame struct {
	Encoding util.Encoding
	Text     string
}

func (tf TextFrame) Bytes() []byte {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.WriteByte(tf.Encoding.Key)
	b.WriteString(tf.Text)

	return b.Bytes()
}
