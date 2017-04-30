// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bogem/id3v2/bbpool"
	"github.com/bogem/id3v2/lrpool"
	"github.com/bogem/id3v2/util"
)

const frameHeaderSize = 10

var errBlankFrame = errors.New("id or size of frame are blank")

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
	if err == errNoTag {
		return newTag(file, 0, 4), nil
	}
	if err != nil {
		err = errors.New("trying to parse tag header: " + err.Error())
		return nil, err
	}
	if header.Version < 3 {
		err = errors.New(fmt.Sprint("unsupported version of ID3 tag: ", header.Version))
		return nil, err
	}

	t := newTag(file, tagHeaderSize+header.FramesSize, header.Version)
	err = t.parseAllFrames()

	return t, err
}

func newTag(file *os.File, originalSize int64, version byte) *Tag {
	return &Tag{
		frames:    make(map[string]Framer),
		sequences: make(map[string]*sequence),

		file:         file,
		originalSize: originalSize,
		version:      version,
	}
}

func (t *Tag) parseAllFrames() error {
	// Initial position of read - beginning of first frame
	if _, err := t.file.Seek(tagHeaderSize, os.SEEK_SET); err != nil {
		return err
	}

	framesSize := t.originalSize - tagHeaderSize
	fileReader := io.LimitReader(t.file, framesSize)

	for {
		id, frame, err := parseFrame(fileReader)
		if err == io.EOF || err == errBlankFrame || err == util.ErrInvalidSizeFormat {
			break
		}
		if err != nil {
			return err
		}

		t.AddFrame(id, frame)
	}

	return nil
}

func parseFrame(rd io.Reader) (id string, frame Framer, err error) {
	header, err := parseFrameHeader(rd)
	if err != nil {
		return "", nil, err
	}
	id = header.ID

	frameRd := lrpool.Get()
	defer lrpool.Put(frameRd)
	frameRd.R = rd
	frameRd.N = header.FrameSize

	frame, err = parseFrameBody(id, frameRd)
	return id, frame, err
}

func parseFrameHeader(rd io.Reader) (frameHeader, error) {
	var header frameHeader

	headerRd := lrpool.Get()
	defer lrpool.Put(headerRd)
	headerRd.R = rd
	headerRd.N = frameHeaderSize

	fhBuf := bbpool.Get()
	defer bbpool.Put(fhBuf)

	_, err := fhBuf.ReadFrom(headerRd)
	if err != nil {
		return header, err
	}
	fhBytes := fhBuf.Bytes()

	id := string(fhBytes[:4])
	frameSize, err := util.ParseSize(fhBytes[4:8])
	if err != nil {
		return header, err
	}

	if id == "" || frameSize == 0 {
		return header, errBlankFrame
	}

	header.ID = id
	header.FrameSize = frameSize
	return header, nil

}

func parseFrameBody(id string, rd io.Reader) (Framer, error) {
	if id[0] == 'T' {
		return parseTextFrame(rd)
	}

	if parseFunc, exists := parsers[id]; exists {
		return parseFunc(rd)
	}

	return parseUnknownFrame(rd)
}
