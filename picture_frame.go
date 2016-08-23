// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"io"

	"github.com/bogem/id3v2/bytesbufferpool"
	"github.com/bogem/id3v2/util"
)

// PictureFrame structure is used for picture frames (APIC).
//
// Example of setting a new picture frame to existing tag:
//
//	frontCover, err := os.Open("artwork.jpg")
//	if err != nil {
//		log.Fatal("Error while opening front cover file")
//	}
//	defer frontCover.Close()
//
//	pic := id3v2.PictureFrame{
//		Encoding:    id3v2.ENUTF8,
//		MimeType:    "image/jpeg",
//		PictureType: id3v2.PTFrontCover,
//		Description: "Front cover",
//		Picture:     frontCover,
//	}
//	tag.AddAttachedPicture(pic)
//
// Available picture types you can see in constants.
type PictureFrame struct {
	Encoding    util.Encoding
	MimeType    string
	PictureType byte
	Description string
	Picture     io.Reader
}

func (pf PictureFrame) Body() []byte {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.WriteByte(pf.Encoding.Key)
	b.WriteString(pf.MimeType)
	b.WriteByte(0)
	b.WriteByte(pf.PictureType)
	b.WriteString(pf.Description)
	b.Write(pf.Encoding.TerminationBytes)

	if _, err := b.ReadFrom(pf.Picture); err != nil {
		panic("can't read a picture: " + err.Error())
	}

	return b.Bytes()
}

func ParsePictureFrame(rd io.Reader) (Framer, error) {
	bufRd := util.NewReader(rd)

	encodingByte, err := bufRd.ReadByte()
	if err != nil {
		return nil, err
	}
	encoding := Encodings[encodingByte]

	mimeType, err := bufRd.ReadTillAndWithDelim(0)
	if err != nil {
		return nil, err
	}

	pictureType, err := bufRd.ReadByte()
	if err != nil {
		return nil, err
	}

	description, err := bufRd.ReadTillAndWithDelims(encoding.TerminationBytes)
	if err != nil {
		return nil, err
	}

	picture, err := bufRd.ReadAll()
	if err != nil {
		return nil, err
	}

	pf := PictureFrame{
		Encoding:    encoding,
		MimeType:    string(mimeType),
		PictureType: pictureType,
		Description: string(description),
		Picture:     bytes.NewReader(picture),
	}

	return pf, nil
}
