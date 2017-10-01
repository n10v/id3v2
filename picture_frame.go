// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import "io"

// PictureFrame structure is used for picture frames (APIC).
// The information about how to add picture frame to tag you can
// see in the docs to tag.AddAttachedPicture function.
//
// Available picture types you can see in constants.
type PictureFrame struct {
	Encoding    Encoding
	MimeType    string
	PictureType byte
	Description string
	Picture     []byte
}

func (pf PictureFrame) Size() int {
	return 1 + len(pf.MimeType) + 1 + 1 + encodedSize(pf.Description, pf.Encoding) +
		len(pf.Encoding.TerminationBytes) + len(pf.Picture)
}

func (pf PictureFrame) WriteTo(w io.Writer) (n int64, err error) {
	bw := getBufioWriter(w)
	defer putBufioWriter(bw)

	bw.WriteByte(pf.Encoding.Key)
	bw.WriteString(pf.MimeType)
	bw.WriteByte(0)
	bw.WriteByte(pf.PictureType)
	_, err = encodeWriteText(bw, pf.Description, pf.Encoding)
	if err != nil {
		return
	}
	bw.Write(pf.Encoding.TerminationBytes)
	bw.Write(pf.Picture)

	return int64(bw.Buffered()), bw.Flush()
}

func parsePictureFrame(rd io.Reader) (Framer, error) {
	bufRd := getUtilReader(rd)
	defer putUtilReader(bufRd)

	encodingKey, err := bufRd.ReadByte()
	if err != nil {
		return nil, err
	}
	encoding := getEncoding(encodingKey)

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

	picture, err := readAll(bufRd)
	if err != nil {
		return nil, err
	}

	pf := PictureFrame{
		Encoding:    encoding,
		MimeType:    string(mimeType),
		PictureType: pictureType,
		Description: decodeText(description, encoding),
		Picture:     picture,
	}

	return pf, nil
}
