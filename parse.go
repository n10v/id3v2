// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"errors"
	"io"
	"os"

	"github.com/bogem/id3v2/util"
)

const frameHeaderSize = 10

var blankFrameError = errors.New("id or size of frame are blank")

type frameHeader struct {
	ID        string
	FrameSize int64
}

func parseTag(file *os.File) (*Tag, error) {
	if file == nil {
		err := errors.New("Invalid file: file is nil")
		return nil, err
	}
	header, err := parseHeader(file)
	if err != nil {
		err = errors.New("Trying to parse tag header: " + err.Error())
		return nil, err
	}
	if header == nil {
		return newTag(file, 0, 4), nil
	}
	if header.Version < 3 {
		err = errors.New("Unsupported version of ID3 tag")
		return nil, err
	}

	t := newTag(file, tagHeaderSize+header.FramesSize, header.Version)
	err = t.parseAllFrames()

	return t, err
}

func newTag(file *os.File, originalSize int64, version byte) *Tag {
	t := &Tag{
		frames:    make(map[string]Framer),
		sequences: make(map[string]sequencer),

		file:         file,
		originalSize: originalSize,
		version:      version,
	}

	if version == 3 {
		t.commonIDs = V23CommonIDs
	} else {
		t.commonIDs = V24CommonIDs
	}

	return t
}

func (t *Tag) parseAllFrames() error {
	// Initial position of read - beginning of first frame
	if _, err := t.file.Seek(tagHeaderSize, os.SEEK_SET); err != nil {
		return err
	}

	size := t.originalSize - tagHeaderSize // Size of all frames = Size of tag - tag header
	fileReader := io.LimitReader(t.file, size)

	for {
		id, frame, err := t.parseFrame(fileReader)
		if err == io.EOF || err == blankFrameError || err == util.InvalidSizeFormat {
			break
		}
		if err != nil {
			return err
		}

		t.AddFrame(id, frame)
	}

	return nil
}

var frameBody = new(io.LimitedReader)

func (t Tag) parseFrame(rd io.Reader) (id string, frame Framer, err error) {
	header, err := parseFrameHeader(rd)
	if err != nil {
		return "", nil, err
	}
	id = header.ID

	parseFunc := t.findParseFunc(id)

	frameBody.R = rd
	frameBody.N = header.FrameSize

	frame, err = parseFunc(frameBody)
	return id, frame, err
}

var fhBuf = make([]byte, frameHeaderSize)

func parseFrameHeader(rd io.Reader) (*frameHeader, error) {
	_, err := rd.Read(fhBuf)
	if err != nil {
		return nil, err
	}

	id := string(fhBuf[:4])
	frameSize, err := util.ParseSize(fhBuf[4:8])
	if err != nil {
		return nil, err
	}

	if id == "" || frameSize == 0 {
		return nil, blankFrameError
	}

	header := &frameHeader{
		ID:        id,
		FrameSize: frameSize,
	}

	return header, nil

}

func (t Tag) findParseFunc(id string) func(io.Reader) (Framer, error) {
	if id[0] == 'T' {
		return parseTextFrame
	}

	switch id {
	case t.CommonID("Attached picture"):
		return parsePictureFrame
	case t.CommonID("Comments"):
		return parseCommentFrame
	case t.CommonID("Unsynchronised lyrics/text transcription"):
		return parseUnsynchronisedLyricsFrame
	}
	return parseUnknownFrame
}
