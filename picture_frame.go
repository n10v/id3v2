// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"io"
	"io/ioutil"

	"github.com/bogem/id3v2/bwpool"
	"github.com/bogem/id3v2/rdpool"
	"github.com/bogem/id3v2/util"
)

// PictureFrame structure is used for picture frames (APIC).
// The information about how to add picture frame to tag you can
// see in the docs to tag.AddAttachedPicture function.
//
// Example of reading picture frames:
//
//	pictures := tag.GetFrames(tag.CommonID("Attached picture"))
//	if pictures != nil {
//		for _, f := range pictures {
//			pic, ok := f.(id3v2.PictureFrame)
//			if !ok {
//				log.Fatal("Couldn't assert picture frame")
//			}
//
//			// Do something with picture frame.
//			// For example, print the description:
//			fmt.Println(pic.Description)
//		}
//	}
//
// Example of adding picture frames to tag:
//
//	artwork, err := ioutil.ReadFile("artwork.jpg")
//	if err != nil {
//		log.Fatal("Error while reading artwork file", err)
//	}
//
//	pic := id3v2.PictureFrame{
//		Encoding:    id3v2.ENUTF8,
//		MimeType:    "image/jpeg",
//		PictureType: id3v2.PTFrontCover,
//		Description: "Front cover",
//		Picture:     artwork,
//	}
//	tag.AddAttachedPicture(pic)
//
// Available picture types you can see in constants.
type PictureFrame struct {
	Encoding    util.Encoding
	MimeType    string
	PictureType byte
	Description string
	Picture     []byte
}

func (pf PictureFrame) Size() int {
	return 1 + len(pf.MimeType) + 1 + 1 + len(pf.Description) +
		len(pf.Encoding.TerminationBytes) + len(pf.Picture)
}

func (pf PictureFrame) WriteTo(w io.Writer) (n int64, err error) {
	var i int
	bw := bwpool.Get(w)
	defer bwpool.Put(bw)

	err = bw.WriteByte(pf.Encoding.Key)
	if err != nil {
		return
	}
	n += 1

	i, err = bw.WriteString(pf.MimeType)
	if err != nil {
		return
	}
	n += int64(i)

	err = bw.WriteByte(0)
	if err != nil {
		return
	}
	n += 1

	err = bw.WriteByte(pf.PictureType)
	if err != nil {
		return
	}
	n += 1

	i, err = bw.WriteString(pf.Description)
	if err != nil {
		return
	}
	n += int64(i)

	i, err = bw.Write(pf.Encoding.TerminationBytes)
	if err != nil {
		return
	}
	n += int64(i)

	i, err = bw.Write(pf.Picture)
	if err != nil {
		return
	}
	n += int64(i)

	err = bw.Flush()
	return
}

func parsePictureFrame(rd io.Reader) (Framer, error) {
	bufRd := rdpool.Get(rd)
	defer rdpool.Put(bufRd)

	encodingByte, err := bufRd.ReadByte()
	if err != nil {
		return nil, err
	}
	encoding := Encodings[encodingByte]

	mimeType, err := bufRd.ReadTillDelim(0)
	if err != nil {
		return nil, err
	}
	if _, err = bufRd.Discard(1); err != nil {
		return nil, err
	}

	pictureType, err := bufRd.ReadByte()
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

	picture, err := ioutil.ReadAll(bufRd)
	if err != nil {
		return nil, err
	}

	pf := PictureFrame{
		Encoding:    encoding,
		MimeType:    string(mimeType),
		PictureType: pictureType,
		Description: string(description),
		Picture:     picture,
	}

	return pf, nil
}
