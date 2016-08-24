// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"io"

	"github.com/bogem/id3v2/bytesbufferpool"
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
//	tag.AddFrame(tag.ID("Comments"), comment)
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
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.WriteByte(cf.Encoding.Key)
	if cf.Language == "" {
		panic("language isn't set up in comment frame with description " + cf.Description)
	}
	b.WriteString(cf.Language)
	b.WriteString(cf.Description)
	b.Write(cf.Encoding.TerminationBytes)
	b.WriteString(cf.Text)

	return b.Bytes()
}

func ParseCommentFrame(rd io.Reader) (Framer, error) {
	bufRd := util.NewReader(rd)

	encodingByte, err := bufRd.ReadByte()
	if err != nil {
		return nil, err
	}
	encoding := Encodings[encodingByte]

	language, err := bufRd.ReadBytes(3)
	if err != nil {
		return nil, err
	}

	description, err := bufRd.ReadTillAndWithDelims(encoding.TerminationBytes)
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
