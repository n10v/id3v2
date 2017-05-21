// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"errors"
	"io"

	"github.com/bogem/id3v2/bwpool"
	"github.com/bogem/id3v2/rdpool"
	"github.com/bogem/id3v2/util"
)

// UnsynchronisedLyricsFrame is used to work with USLT frames.
// The information about how to add unsynchronised lyrics/text frame to tag
// you can see in the docs to tag.AddUnsynchronisedLyricsFrame function.
//
// You must choose a three-letter language code from
// ISO 639-2 code list: https://www.loc.gov/standards/iso639-2/php/code_list.php
type UnsynchronisedLyricsFrame struct {
	Encoding          util.Encoding
	Language          string
	ContentDescriptor string
	Lyrics            string
}

func (uslf UnsynchronisedLyricsFrame) Size() int {
	return 1 + len(uslf.Language) + len(uslf.ContentDescriptor) +
		+len(uslf.Encoding.TerminationBytes) + len(uslf.Lyrics)
}

func (uslf UnsynchronisedLyricsFrame) WriteTo(w io.Writer) (n int64, err error) {
	var i int
	bw := bwpool.Get(w)
	defer bwpool.Put(bw)

	err = bw.WriteByte(uslf.Encoding.Key)
	if err != nil {
		return
	}
	n++

	if len(uslf.Language) != 3 {
		return n, errors.New("language code must consist of three letters according to ISO 639-2")
	}
	i, err = bw.WriteString(uslf.Language)
	if err != nil {
		return
	}
	n += int64(i)

	i, err = bw.WriteString(uslf.ContentDescriptor)
	if err != nil {
		return
	}
	n += int64(i)

	i, err = bw.Write(uslf.Encoding.TerminationBytes)
	if err != nil {
		return
	}
	n += int64(i)

	i, err = bw.WriteString(uslf.Lyrics)
	if err != nil {
		return
	}
	n += int64(i)

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

	language, err := bufRd.Next(3)
	if err != nil {
		return nil, err
	}

	contentDescriptor, err := bufRd.ReadTillDelims(encoding.TerminationBytes)
	if err != nil {
		return nil, err
	}
	if _, err = bufRd.Discard(len(encoding.TerminationBytes)); err != nil {
		return nil, err
	}

	lyrics, err := bufRd.String()
	if err != nil {
		return nil, err
	}

	uslf := UnsynchronisedLyricsFrame{
		Encoding:          encoding,
		Language:          string(language),
		ContentDescriptor: string(contentDescriptor),
		Lyrics:            lyrics,
	}

	return uslf, nil
}
