// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import "io"

// CommentFrame is used to work with COMM frames.
// The information about how to add comment frame to tag you can
// see in the docs to tag.AddCommentFrame function.
//
// You must choose a three-letter language code from
// ISO 639-2 code list: https://www.loc.gov/standards/iso639-2/php/code_list.php
type CommentFrame struct {
	Encoding    Encoding
	Language    string
	Description string
	Text        string
}

func (cf CommentFrame) Size() int {
	return 1 + len(cf.Language) + encodedSize(cf.Description, cf.Encoding) +
		+len(cf.Encoding.TerminationBytes) + encodedSize(cf.Text, cf.Encoding)
}

func (cf CommentFrame) WriteTo(w io.Writer) (n int64, err error) {
	if len(cf.Language) != 3 {
		return n, ErrInvalidLanguageLength
	}

	bw := getBufioWriter(w)
	defer putBufioWriter(bw)

	bw.WriteByte(cf.Encoding.Key)
	bw.WriteString(cf.Language)
	_, err = encodeWriteText(bw, cf.Description, cf.Encoding)
	if err != nil {
		return
	}
	bw.Write(cf.Encoding.TerminationBytes)
	_, err = encodeWriteText(bw, cf.Text, cf.Encoding)
	if err != nil {
		return
	}

	return int64(bw.Buffered()), bw.Flush()
}

func parseCommentFrame(rd io.Reader) (Framer, error) {
	bufRd := getUtilReader(rd)
	defer putUtilReader(bufRd)

	encodingKey, err := bufRd.ReadByte()
	if err != nil {
		return nil, err
	}
	encoding := getEncoding(encodingKey)

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

	text := getBytesBuffer()
	defer putBytesBuffer(text)

	if _, err := text.ReadFrom(bufRd); err != nil {
		return nil, err
	}

	cf := CommentFrame{
		Encoding:    encoding,
		Language:    string(language),
		Description: decodeText(description, encoding),
		Text:        decodeText(text.Bytes(), encoding),
	}

	return cf, nil
}
