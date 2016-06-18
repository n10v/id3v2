// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package frame

import (
	"bytes"
	"errors"
	"github.com/bogem/id3v2/util"
)

// CommentFramer is used to work with COMM frames.
type CommentFramer interface {
	Framer

	Encoding() util.Encoding
	SetEncoding(util.Encoding)

	Language() string
	SetLanguage(string)

	Description() string
	SetDescription(string)

	Text() string
	SetText(string)
}

// Just implementation of CommentFramer interface.
type CommentFrame struct {
	encoding    util.Encoding
	language    string
	description bytes.Buffer
	text        bytes.Buffer
}

func (cf CommentFrame) Bytes() ([]byte, error) {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()
	defer bytesBufPool.Put(b)

	b.WriteByte(cf.encoding.Key)
	if cf.language == "" {
		return nil, errors.New("Language isn't set up in comment frame with description " + cf.Description())
	}
	b.WriteString(cf.language)
	b.WriteString(cf.Description())
	b.Write(cf.encoding.TerminationBytes)
	b.WriteString(cf.Text())

	return b.Bytes(), nil
}

func (cf CommentFrame) Encoding() util.Encoding {
	return cf.encoding
}

func (cf *CommentFrame) SetEncoding(e util.Encoding) {
	cf.encoding = e
}

func (cf CommentFrame) Language() string {
	return cf.language
}

func (cf *CommentFrame) SetLanguage(lang string) {
	cf.language = lang
}

func (cf CommentFrame) Description() string {
	return cf.description.String()
}

func (cf *CommentFrame) SetDescription(d string) {
	cf.description.Reset()
	cf.description.WriteString(d)
}

func (cf CommentFrame) Text() string {
	return cf.text.String()
}

func (cf *CommentFrame) SetText(text string) {
	cf.text.Reset()
	cf.text.WriteString(text)
}
