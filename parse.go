// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"errors"
	"io"
	"strconv"

	"github.com/bogem/id3v2/bspool"
	"github.com/bogem/id3v2/lrpool"
	"github.com/bogem/id3v2/util"
)

const frameHeaderSize = 10

var errBlankFrame = errors.New("id or size of frame are blank")

type frameHeader struct {
	ID       string
	BodySize int64
}

// parse finds ID3v2 tag in rd and parses it to tag considering opts.
// If rd is smaller than expected, it returns ErrSmallHeaderSize.
func (tag *Tag) parse(rd io.Reader, opts Options) error {
	if rd == nil {
		return errors.New("rd is nil")
	}

	header, err := parseHeader(rd)
	if err == errNoTag || err == io.EOF {
		tag.init(rd, 0, 4)
		return nil
	}
	if err != nil {
		return errors.New("error by parsing tag header: " + err.Error())
	}
	if header.Version < 3 {
		err = errors.New("unsupported version of ID3 tag: " + strconv.Itoa(int(header.Version)))
		return err
	}

	tag.init(rd, tagHeaderSize+header.FramesSize, header.Version)
	if !opts.Parse {
		return nil
	}
	return tag.parseFrames(opts)
}

// init initializes tag by deleting all frames in it, setting reader,
// originialSize and version.
// init doesn't parse frames, it only set the fields.
func (tag *Tag) init(rd io.Reader, originalSize int64, version byte) {
	tag.DeleteAllFrames()
	tag.reader = rd
	tag.originalSize = originalSize
	tag.version = version
}

func (tag *Tag) parseFrames(opts Options) error {
	// Size of frames in tag = size of whole tag - size of tag header.
	framesSize := tag.originalSize - tagHeaderSize

	// Convert descriptions, specified by user in opts.ParseFrames, to IDs.
	// Use map for speed.
	parseIDs := make(map[string]bool, len(opts.ParseFrames))
	for _, description := range opts.ParseFrames {
		parseIDs[tag.CommonID(description)] = true
	}

	buf := bspool.Get(32 * 1024)
	defer bspool.Put(buf)
	for framesSize > 0 {
		// Parse frame header.
		header, err := parseFrameHeader(buf, tag.reader)
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

		// Limit tag.reader by header.BodySize.
		bodyRd := lrpool.Get(tag.reader, bodySize)
		defer lrpool.Put(bodyRd)

		// If user set opts.ParseFrames, take it into consideration.
		if len(parseIDs) > 0 && !parseIDs[id] {
			if err := skipReaderBuf(bodyRd, buf); err != nil {
				return err
			}
			continue
		}

		// Parse frame body.
		frame, err := parseFrameBody(id, bodyRd)
		if err != nil && err != io.EOF {
			return err
		}

		// Add frame to tag.
		tag.AddFrame(id, frame)

		if err == io.EOF {
			break
		}
	}

	return nil
}

func parseFrameHeader(buf []byte, rd io.Reader) (frameHeader, error) {
	var header frameHeader

	if len(buf) < frameHeaderSize {
		return header, errors.New("parseFrameHeader: buf is smaller than frame header size")
	}

	fhBuf := buf[:frameHeaderSize]
	if _, err := rd.Read(fhBuf); err != nil {
		return header, err
	}

	id := string(fhBuf[:4])
	bodySize, err := util.ParseSize(fhBuf[4:8])
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

// skipReaderBuf just reads rd until io.EOF.
func skipReaderBuf(rd io.Reader, buf []byte) error {
	for {
		_, err := rd.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
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
