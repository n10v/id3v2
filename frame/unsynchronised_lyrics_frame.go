// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package frame

import (
	"errors"

	"github.com/bogem/id3v2/bytesbufferpool"
	"github.com/bogem/id3v2/util"
)

// UnsynchronisedLyricsFramer is used to work with USLT frames.
type UnsynchronisedLyricsFramer interface {
	Framer

	Encoding() util.Encoding
	SetEncoding(util.Encoding)

	Language() string
	SetLanguage(string)

	ContentDescriptor() string
	SetContentDescriptor(string)

	Lyrics() string
	SetLyrics(string)
}

// Just implementation of UnsynchronisedLyricsFramer interface.
type UnsynchronisedLyricsFrame struct {
	encoding          util.Encoding
	language          string
	contentDescriptor string
	lyrics            string
}

func (uslf UnsynchronisedLyricsFrame) Bytes() ([]byte, error) {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.WriteByte(uslf.encoding.Key)
	if uslf.language == "" {
		return nil, errors.New("Language isn't set up in USLT frame with description " + uslf.ContentDescriptor())
	}
	b.WriteString(uslf.language)
	b.WriteString(uslf.contentDescriptor)
	b.Write(uslf.encoding.TerminationBytes)
	b.WriteString(uslf.lyrics)

	return b.Bytes(), nil
}

func (uslf UnsynchronisedLyricsFrame) Encoding() util.Encoding {
	return uslf.encoding
}

func (uslf *UnsynchronisedLyricsFrame) SetEncoding(e util.Encoding) {
	uslf.encoding = e
}

func (uslf UnsynchronisedLyricsFrame) Language() string {
	return uslf.language
}

func (uslf *UnsynchronisedLyricsFrame) SetLanguage(lang string) {
	uslf.language = lang
}

func (uslf UnsynchronisedLyricsFrame) ContentDescriptor() string {
	return uslf.contentDescriptor
}

func (uslf *UnsynchronisedLyricsFrame) SetContentDescriptor(cd string) {
	uslf.contentDescriptor = cd
}

func (uslf UnsynchronisedLyricsFrame) Lyrics() string {
	return uslf.lyrics
}

func (uslf *UnsynchronisedLyricsFrame) SetLyrics(lyrics string) {
	uslf.lyrics = lyrics
}
