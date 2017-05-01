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
	ID       string
	BodySize int64
}

func parseTag(file *os.File, opts Options) (*Tag, error) {
	if file == nil {
		return nil, errors.New("file is nil")
	}

	header, err := parseHeader(file)
	if err == errNoTag {
		return newTag(file, 0, 4), nil
	}
	if err != nil {
		return nil, errors.New("error by parsing tag header: " + err.Error())
	}
	if header.Version < 3 {
		err = errors.New(fmt.Sprint("unsupported version of ID3 tag: ", header.Version))
		return nil, err
	}

	t := newTag(file, tagHeaderSize+header.FramesSize, header.Version)
	if opts.Parse {
		err = t.parseAllFrames(opts)
	}

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

func (t *Tag) parseAllFrames(opts Options) error {
	// Initial position of read - beginning of first frame.
	if _, err := t.file.Seek(tagHeaderSize, os.SEEK_SET); err != nil {
		return err
	}

	// Size of frames in tag = size of whole tag - size of tag header.
	framesSize := t.originalSize - tagHeaderSize

	// Convert descriptions, specified by user in opts.ParseFrames, to IDs.
	// Use map for speed.
	parseIDs := make(map[string]bool, len(opts.ParseFrames))
	for _, description := range opts.ParseFrames {
		parseIDs[t.CommonID(description)] = true
	}

	for framesSize > 0 {
		// Parse frame header.
		header, err := parseFrameHeader(t.file)
		if err == io.EOF || err == errBlankFrame || err == util.ErrInvalidSizeFormat {
			break
		}
		if err != nil {
			return err
		}
		id := header.ID
		bodySize := header.BodySize

		// Substitute the size of the whole frame from framesSize.
		framesSize -= frameHeaderSize + bodySize

		// If user set opts.ParseFrames, take it into consideration.
		if len(parseIDs) > 0 {
			if !parseIDs[id] {
				_, err = t.file.Seek(bodySize, os.SEEK_CUR)
				continue
			}
		}

		// Limit t.file by header.BodySize.
		bodyRd := lrpool.Get(t.file, bodySize)
		defer lrpool.Put(bodyRd)

		// Parse frame body.
		frame, err := parseFrameBody(id, bodyRd)
		if err != nil && err != io.EOF {
			return err
		}

		// Add frame to tag.
		t.AddFrame(id, frame)

		if err == io.EOF {
			break
		}
	}

	return nil
}

func parseFrameHeader(rd io.Reader) (frameHeader, error) {
	var header frameHeader

	// Limit rd by frameHeaderSize.
	bodyRd := lrpool.Get(rd, frameHeaderSize)
	defer lrpool.Put(bodyRd)

	fhBuf := bbpool.Get()
	defer bbpool.Put(fhBuf)

	_, err := fhBuf.ReadFrom(bodyRd)
	if err != nil {
		return header, err
	}
	data := fhBuf.Bytes()

	id := string(data[:4])
	bodySize, err := util.ParseSize(data[4:8])
	if err != nil {
		return header, err
	}

	if id == "" || bodySize == 0 {
		return header, errBlankFrame
	}

	header.ID = id
	header.BodySize = bodySize
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
