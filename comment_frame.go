// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"io"

	"github.com/bogem/id3v2/bwpool"
	"github.com/bogem/id3v2/rdpool"
	"github.com/bogem/id3v2/util"
)

// CommentFrame is used to work with COMM frames.
//
// Example of setting a new comment frame to existing tag:
//
//	comment := id3v2.CommentFrame{
//		Encoding:   id3v2.ENUTF8,
//		Language:   "eng",
//		Desciption: "My opinion",
//		Text:       "Very good song",
//	}
//	tag.AddCommentFrame(comment)
//
// You should choose a language code from
// ISO 639-2 code list: https://www.loc.gov/standards/iso639-2/php/code_list.php
type CommentFrame struct {
	Encoding    util.Encoding
	Language    string
	Description string
	Text        string
}

func (cf CommentFrame) Body() []byte {
	b := new(bytes.Buffer)

	b.WriteByte(cf.Encoding.Key)
	b.WriteString(cf.Language)
	b.WriteString(cf.Description)
	b.Write(cf.Encoding.TerminationBytes)
	b.WriteString(cf.Text)

	return b.Bytes()
}

func (cf CommentFrame) Size() int {
	return 1 + len(cf.Language) + len(cf.Description) +
		+len(cf.Encoding.TerminationBytes) + len(cf.Text)
}

func (cf CommentFrame) WriteTo(w io.Writer) (n int, err error) {
	var i int
	bw := bwpool.Get(w)
	defer bwpool.Put(bw)

	err = bw.WriteByte(cf.Encoding.Key)
	if err != nil {
		return
	}
	n += 1

	i, err = bw.WriteString(cf.Language)
	if err != nil {
		return
	}
	n += i

	i, err = bw.WriteString(cf.Description)
	if err != nil {
		return
	}
	n += i

	i, err = bw.Write(cf.Encoding.TerminationBytes)
	if err != nil {
		return
	}
	n += i

	i, err = bw.WriteString(cf.Text)
	if err != nil {
		return
	}
	n += i

	err = bw.Flush()
	return
}

func parseCommentFrame(rd io.Reader) (Framer, error) {
	bufRd := rdpool.Get(rd)
	defer rdpool.Put(bufRd)

	encodingByte, err := bufRd.ReadByte()
	if err != nil {
		return nil, err
	}
	encoding := Encodings[encodingByte]

	language, err := bufRd.ReadSeveralBytes(3)
	if err != nil {
		return nil, err
	}

	description, err := bufRd.ReadTillDelims(encoding.TerminationBytes)
	if err != nil {
		return nil, err
	}

	text, err := bufRd.ReadAll()
	if err != nil {
		return nil, err
	}

	cf := CommentFrame{
		Encoding:    encoding,
		Language:    string(language),
		Description: string(description),
		Text:        string(text),
	}

	return cf, nil
}
