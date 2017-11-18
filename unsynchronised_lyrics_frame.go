// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import "io"

// UnsynchronisedLyricsFrame is used to work with USLT frames.
// The information about how to add unsynchronised lyrics/text frame to tag
// you can see in the docs to tag.AddUnsynchronisedLyricsFrame function.
//
// You must choose a three-letter language code from
// ISO 639-2 code list: https://www.loc.gov/standards/iso639-2/php/code_list.php
type UnsynchronisedLyricsFrame struct {
	Encoding          Encoding
	Language          string
	ContentDescriptor string
	Lyrics            string
}

func (uslf UnsynchronisedLyricsFrame) Size() int {
	return 1 + len(uslf.Language) + encodedSize(uslf.ContentDescriptor, uslf.Encoding) +
		+len(uslf.Encoding.TerminationBytes) + encodedSize(uslf.Lyrics, uslf.Encoding)
}

func (uslf UnsynchronisedLyricsFrame) WriteTo(w io.Writer) (n int64, err error) {
	if len(uslf.Language) != 3 {
		return n, ErrInvalidLanguageLength
	}

	bw, ok := resolveBufioWriter(w)
	if !ok {
		defer putBufioWriter(bw)
	}

	var nn int

	bw.WriteByte(uslf.Encoding.Key)
	n += 1

	nn, _ = bw.WriteString(uslf.Language)
	n += int64(nn)

	nn, err = encodeWriteText(bw, uslf.ContentDescriptor, uslf.Encoding)
	n += int64(nn)
	if err != nil {
		return
	}

	nn, _ = bw.Write(uslf.Encoding.TerminationBytes)
	n += int64(nn)

	nn, err = encodeWriteText(bw, uslf.Lyrics, uslf.Encoding)
	n += int64(nn)
	if err != nil {
		return
	}

	return n, bw.Flush()
}

func parseUnsynchronisedLyricsFrame(rd io.Reader) (Framer, error) {
	bufRd := getUtilReader(rd)
	defer putUtilReader(bufRd)

	encodingKey, err := bufRd.ReadByte()
	if err != nil {
		return nil, err
	}
	encoding := getEncoding(encodingKey)

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

	lyrics := getBytesBuffer()
	defer putBytesBuffer(lyrics)

	if _, err := lyrics.ReadFrom(bufRd); err != nil {
		return nil, err
	}

	uslf := UnsynchronisedLyricsFrame{
		Encoding:          encoding,
		Language:          string(language),
		ContentDescriptor: decodeText(contentDescriptor, encoding),
		Lyrics:            decodeText(lyrics.Bytes(), encoding),
	}

	return uslf, nil
}
