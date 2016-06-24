// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package frame

import (
	"bytes"
	"errors"

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
	contentDescriptor bytes.Buffer
	lyrics            bytes.Buffer
}

func (uslf UnsynchronisedLyricsFrame) Bytes() ([]byte, error) {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()

	b.WriteByte(uslf.encoding.Key)
	if uslf.language == "" {
		return nil, errors.New("Language isn't set up in USLT frame with description " + uslf.ContentDescriptor())
	}
	b.WriteString(uslf.language)
	b.WriteString(uslf.ContentDescriptor())
	b.Write(uslf.encoding.TerminationBytes)
	b.WriteString(uslf.Lyrics())

	bytesBufPool.Put(b)
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
	return uslf.contentDescriptor.String()
}

func (uslf *UnsynchronisedLyricsFrame) SetContentDescriptor(cd string) {
	uslf.contentDescriptor.Reset()
	uslf.contentDescriptor.WriteString(cd)
}

func (uslf UnsynchronisedLyricsFrame) Lyrics() string {
	return uslf.lyrics.String()
}

func (uslf *UnsynchronisedLyricsFrame) SetLyrics(lyrics string) {
	uslf.lyrics.Reset()
	uslf.lyrics.WriteString(lyrics)
}
