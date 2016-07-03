// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"github.com/bogem/id3v2/bytesbufferpool"
	"github.com/bogem/id3v2/util"
)

// TextFrame is used to work with all text frames
// (all T*** frames like TIT2, TALB and so on).
type TextFrame struct {
	Encoding util.Encoding
	Text     string
}

func (tf TextFrame) Bytes() ([]byte, error) {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.WriteByte(tf.Encoding.Key)
	b.WriteString(tf.Text)

	return b.Bytes(), nil
}
