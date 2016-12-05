// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"bytes"
	"io"

	"github.com/bogem/id3v2/rdpool"
	"github.com/bogem/id3v2/util"
)

// UnsynchronisedLyricsFrame is used to work with USLT frames.
//
// Example of setting a new unsynchronised lyrics/text frame to existing tag:
//
//	uslt := id3v2.UnsynchronisedLyricsFrame{
//		Encoding:          id3v2.ENUTF8,
//		Language:          "ger",
//		ContentDescriptor: "Deutsche Nationalhymne",
//		Lyrics:            "Einigkeit und Recht und Freiheit...",
//	}
//	tag.AddUnsynchronisedLyricsFrame(uslt)
type UnsynchronisedLyricsFrame struct {
	Encoding          util.Encoding
	Language          string
	ContentDescriptor string
	Lyrics            string
}

func (uslf UnsynchronisedLyricsFrame) Body() []byte {
	b := new(bytes.Buffer)

	b.WriteByte(uslf.Encoding.Key)
	b.WriteString(uslf.Language)
	b.WriteString(uslf.ContentDescriptor)
	b.Write(uslf.Encoding.TerminationBytes)
	b.WriteString(uslf.Lyrics)

	return b.Bytes()
}

func (uslf UnsynchronisedLyricsFrame) Size() int {
	return 1 + len(uslf.Language) + len(uslf.ContentDescriptor) +
		+len(uslf.Encoding.TerminationBytes) + len(uslf.Lyrics)
}

func (uslf UnsynchronisedLyricsFrame) WriteTo(w io.Writer) (n int, err error) {
	var i int
	bw := bufio.NewWriter(w)

	err = bw.WriteByte(uslf.Encoding.Key)
	if err != nil {
		return
	}
	n += 1

	i, err = bw.WriteString(uslf.Language)
	if err != nil {
		return
	}
	n += i

	i, err = bw.WriteString(uslf.ContentDescriptor)
	if err != nil {
		return
	}
	n += i

	i, err = bw.Write(uslf.Encoding.TerminationBytes)
	if err != nil {
		return
	}
	n += i

	i, err = bw.WriteString(uslf.Lyrics)
	if err != nil {
		return
	}
	n += i

	err = bw.Flush()
	return
}

func parseUnsynchronisedLyricsFrame(rd io.Reader) (Framer, error) {
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

	contentDescriptor, err := bufRd.ReadTillDelims(encoding.TerminationBytes)
	if err != nil {
		return nil, err
	}

	lyrics, err := bufRd.ReadAll()
	if err != nil {
		return nil, err
	}

	uslf := UnsynchronisedLyricsFrame{
		Encoding:          encoding,
		Language:          string(language),
		ContentDescriptor: string(contentDescriptor),
		Lyrics:            string(lyrics),
	}

	return uslf, nil
}
