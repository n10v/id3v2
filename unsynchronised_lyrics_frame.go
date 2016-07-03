// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"errors"

	"github.com/bogem/id3v2/bytesbufferpool"
	"github.com/bogem/id3v2/util"
)

// UnsynchronisedLyricsFrame is used to work with USLT frames.
type UnsynchronisedLyricsFrame struct {
	ContentDescriptor string
	Encoding          util.Encoding
	Language          string
	Lyrics            string
}

func (uslf UnsynchronisedLyricsFrame) Bytes() ([]byte, error) {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.WriteByte(uslf.Encoding.Key)
	if uslf.Language == "" {
		return nil, errors.New("Language isn't set up in USLT frame with description " + uslf.ContentDescriptor)
	}
	b.WriteString(uslf.Language)
	b.WriteString(uslf.ContentDescriptor)
	b.Write(uslf.Encoding.TerminationBytes)
	b.WriteString(uslf.Lyrics)

	return b.Bytes(), nil
}
