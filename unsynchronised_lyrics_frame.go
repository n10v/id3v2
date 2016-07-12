// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"github.com/bogem/id3v2/bytesbufferpool"
	"github.com/bogem/id3v2/util"
)

// UnsynchronisedLyricsFrame is used to work with USLT frames.
//
// Example of setting a new unsynchronised lyrics/text frame to existing tag:
//  uslt := id3v2.UnsynchronisedLyricsFrame{
//    Encoding:          id3v2.ENUTF8,
//    Language:          "ger",
//    ContentDescriptor: "Deutsche Nationalhymne",
//    Lyrics:            "Einigkeit und Recht und Freiheit...",
//  }
//  tag.AddUnsynchronisedLyricsFrame(uslt)
type UnsynchronisedLyricsFrame struct {
	Encoding          util.Encoding
	Language          string
	ContentDescriptor string
	Lyrics            string
}

func (uslf UnsynchronisedLyricsFrame) Bytes() []byte {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.WriteByte(uslf.Encoding.Key)
	if uslf.Language == "" {
		panic("language isn't set up in USLT frame with description " + uslf.ContentDescriptor)
	}
	b.WriteString(uslf.Language)
	b.WriteString(uslf.ContentDescriptor)
	b.Write(uslf.Encoding.TerminationBytes)
	b.WriteString(uslf.Lyrics)

	return b.Bytes()
}
