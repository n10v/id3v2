// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"io"

	"github.com/bogem/id3v2/bwpool"
	"github.com/bogem/id3v2/rdpool"
	"github.com/bogem/id3v2/util"
)

// TextFrame is used to work with all text frames
// (all T*** frames like TIT2 (title), TALB (album) and so on).
//
// Example of reading text frames:
//
//	tf := tag.GetTextFrame(tag.CommonID("Mood"))
//	fmt.Println(tf.Text)
//
// Example of adding text frames to tag:
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

func (tf TextFrame) Size() int {
	return 1 + len(tf.Text)
}

func (tf TextFrame) WriteTo(w io.Writer) (n int64, err error) {
	var i int
	bw := bwpool.Get(w)
	defer bwpool.Put(bw)

	err = bw.WriteByte(tf.Encoding.Key)
	if err != nil {
		return
	}
	n += 1

	i, err = bw.WriteString(tf.Text)
	if err != nil {
		return
	}
	n += int64(i)

	err = bw.Flush()
	return
}

func parseTextFrame(rd io.Reader) (Framer, error) {
	tfRd := rdpool.Get(rd)
	defer rdpool.Put(tfRd)

	encoding, err := tfRd.ReadByte()
	if err != nil {
		return nil, err
	}

	text, err := tfRd.String()
	if err != nil {
		return nil, err
	}

	tf := TextFrame{
		Encoding: Encodings[encoding],
		Text:     text,
	}

	return tf, nil
}
