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

// PictureFrame structure is used for picture frames (APIC).
//
// Example of setting a new picture frame to existing tag:
//	frontCover, err := os.Open("artwork.jpg")
//	if err != nil {
//		log.Fatal("Error while opening front cover file")
//  }
//	defer frontCover.Close()
//
//	pic := id3v2.PictureFrame{
//		Encoding:    id3v2.ENUTF8,
//		MimeType:    "image/jpeg",
//		Description: "Front cover",
//		Picture:     frontCover,
//		PictureType: id3v2.PTFrontCover,
//	}
//	tag.AddAttachedPicture(pic)
//
// Available picture types you can see in constants.
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
