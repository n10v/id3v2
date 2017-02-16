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

// CommentFrame is used to work with COMM frames.
// The information about how to add comment frame to tag you can
// see in the docs to tag.AddCommentFrame function.
type CommentFrame struct {
	Encoding    util.Encoding
	Language    string
	Description string
	Text        string
}

func (cf CommentFrame) Size() int {
	return 1 + len(cf.Language) + len(cf.Description) +
		+len(cf.Encoding.TerminationBytes) + len(cf.Text)
}

func (cf CommentFrame) WriteTo(w io.Writer) (n int64, err error) {
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
	n += int64(i)

	i, err = bw.WriteString(cf.Description)
	if err != nil {
		return
	}
	n += int64(i)

	i, err = bw.Write(cf.Encoding.TerminationBytes)
	if err != nil {
		return
	}
	n += int64(i)

	i, err = bw.WriteString(cf.Text)
	if err != nil {
		return
	}
	n += int64(i)

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

	language, err := bufRd.Next(3)
	if err != nil {
		return nil, err
	}

	description, err := bufRd.ReadTillDelims(encoding.TerminationBytes)
	if err != nil {
		return nil, err
	}
	if _, err = bufRd.Discard(len(encoding.TerminationBytes)); err != nil {
		return nil, err
	}

	text, err := bufRd.String()
	if err != nil {
		return nil, err
	}

	cf := CommentFrame{
		Encoding:    encoding,
		Language:    string(language),
		Description: string(description),
		Text:        text,
	}

	return cf, nil
}
