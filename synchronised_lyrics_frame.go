// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"io"
	"math/big"
)

// SynchronisedLyricsFrame is used to work with USLT frames.
// The information about how to add unsynchronised lyrics/text frame to tag
// you can see in the docs to tag.AddUnsynchronisedLyricsFrame function.
//
// You must choose a three-letter language code from
// ISO 639-2 code list: https://www.loc.gov/standards/iso639-2/php/code_list.php
type SynchronisedLyricsFrame struct {
	Encoding             Encoding
	Language             string
	TimestampFormat      byte
	ContentType          byte
	ContentDescriptor    string
	SynchronizedTextSpec []SyncedText
}

type TimestampFormat int

const (
	Unknown TimestampFormat = iota
	AbsoluteMpegFrames
	AbsoluteMilliseconds
)

var (
	ContentType = map[int]string{
		0: "Other",
		1: "Lyrics",
		2: "Transcription",
		3: "Movement",
		4: "Events",
		5: "Chord",
		6: "Trivia",
		7: "WebpageUrls",
		8: "ImageUrls",
	}
)

type SyncedText struct {
	Text      string
	Timestamp *big.Int
}

func (sylf SynchronisedLyricsFrame) Size() int {
	// s := binary.Size(sylf.SynchronizedTextSpec)
	var s int
	for _, v := range sylf.SynchronizedTextSpec {
		s += len(v.Text)
		s += len(sylf.Encoding.TerminationBytes)
		s += len(v.counterBytes())
	}

	return 1 + len(sylf.Language) + encodedSize(sylf.ContentDescriptor, sylf.Encoding) +
		+len(sylf.Encoding.TerminationBytes) + s +
		+1 + 1
}

func (sylf SynchronisedLyricsFrame) UniqueIdentifier() string {
	return sylf.Language + sylf.ContentDescriptor
}

func (sy SyncedText) counterBytes() []byte {
	bytes := sy.Timestamp.Bytes()
	// bs := make([]byte, 4)

	// binary.LittleEndian.PutUint64(bs, uint64(sy.Timestamp))

	// Specification requires at least 4 bytes for counter, pad if necessary.
	bytesNeeded := 4 - len(bytes)
	if bytesNeeded > 0 {
		padding := make([]byte, bytesNeeded)
		bytes = append(padding, bytes...)
	}

	return bytes
}

func (sylf SynchronisedLyricsFrame) WriteTo(w io.Writer) (n int64, err error) {
	if len(sylf.Language) != 3 {
		return n, ErrInvalidLanguageLength
	}
	return useBufWriter(w, func(bw *bufWriter) {
		bw.WriteByte(sylf.Encoding.Key)
		bw.WriteString(sylf.Language)
		bw.WriteByte(sylf.TimestampFormat)
		bw.WriteByte(sylf.ContentType)
		bw.EncodeAndWriteText(sylf.ContentDescriptor, sylf.Encoding)
		bw.Write(sylf.Encoding.TerminationBytes)
		for _, v := range sylf.SynchronizedTextSpec {
			bw.EncodeAndWriteText(v.Text, sylf.Encoding)
			bw.Write(sylf.Encoding.TerminationBytes)
			bw.Write(v.counterBytes())
		}
	})
}

func parseSynchronisedLyricsFrame(br *bufReader) (Framer, error) {
	encoding := getEncoding(br.ReadByte())
	language := br.Next(3)
	timestampFormat := br.Next(1)
	contentType := br.Next(1)
	contentDescriptor := br.ReadText(encoding)

	if br.Err() != nil {
		return nil, br.Err()
	}

	lyrics := getBytesBuffer()
	defer putBytesBuffer(lyrics)

	if _, err := lyrics.ReadFrom(br); err != nil {
		return nil, err
	}
	d := decodeText(lyrics.Bytes(), encoding)
	var y []SyncedText
	for {
		idx := bytes.IndexByte([]byte(d), '\x00')
		t := SyncedText{Text: d[:idx]}
		d = d[idx+1:]

		timeStampBigInt, _ := new(big.Int).SetString(d[:2], 10)

		t.Timestamp = timeStampBigInt
		d = d[2:]

		y = append(y, t)

		if len(d) < 3 || bytes.IndexByte([]byte(d), '\x00') < 2 {
			break
		}
	}
	sylf := SynchronisedLyricsFrame{
		Encoding:             encoding,
		Language:             string(language),
		TimestampFormat:      timestampFormat[0],
		ContentType:          contentType[0],
		ContentDescriptor:    decodeText(contentDescriptor, encoding),
		SynchronizedTextSpec: y,
	}

	return sylf, nil
}
