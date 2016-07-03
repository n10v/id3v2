// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"errors"
	"io"

	"github.com/bogem/id3v2/bytesbufferpool"
	"github.com/bogem/id3v2/util"
)

// PictureFrame is used to work with APIC frames.
type PictureFrame struct {
	Description string
	Encoding    util.Encoding
	MimeType    string
	Picture     io.Reader
	PictureType byte
}

func (pf PictureFrame) Bytes() ([]byte, error) {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.WriteByte(pf.Encoding.Key)
	b.WriteString(pf.MimeType)
	b.WriteByte(0)
	if pf.PictureType < 0 || pf.PictureType > 20 {
		return nil, errors.New("Incorrect picture type in picture frame with description " + pf.Description)
	}
	b.WriteByte(pf.PictureType)
	b.WriteString(pf.Description)
	b.Write(pf.Encoding.TerminationBytes)

	if _, err := b.ReadFrom(pf.Picture); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
