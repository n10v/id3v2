// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

const frameHeaderSize = 10

var ErrUnsupportedVersion = errors.New("unsupported version of ID3 tag")
var errBlankFrame = errors.New("id or size of frame are blank")

// ErrBodyOverflow is returned when a frame has greater size than the remaining tag size
var ErrBodyOverflow = errors.New("frame went over tag area")

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
		return fmt.Errorf("error by parsing tag header: %v", err)
	}
	if header.Version < 3 {
		return ErrUnsupportedVersion
	}

	tag.init(rd, tagHeaderSize+header.FramesSize, header.Version)
	if !opts.Parse {
		return nil
	}
	return tag.parseFrames(opts)
}

func (tag *Tag) init(rd io.Reader, originalSize int64, version byte) {
	tag.DeleteAllFrames()

	tag.reader = rd
	tag.originalSize = originalSize
	tag.version = version
	tag.setDefaultEncodingBasedOnVersion(version)
}

func (tag *Tag) parseFrames(opts Options) error {
	framesSize := tag.originalSize - tagHeaderSize

	parseableIDs := tag.makeIDsFromDescriptions(opts.ParseFrames)
	isParseFramesProvided := len(opts.ParseFrames) > 0

	synchSafe := tag.Version() == 4

	br := getBufReader(nil)
	defer putBufReader(br)

	buf := getByteSlice(32 * 1024)
	defer putByteSlice(buf)

	for framesSize > 0 {
		header, err := parseFrameHeader(buf, tag.reader, synchSafe)
		if err == io.EOF || err == errBlankFrame || err == ErrInvalidSizeFormat {
			break
		}
		if err != nil {
			return err
		}
		id, bodySize := header.ID, header.BodySize
		fmt.Println(id, framesSize, bodySize)

		framesSize -= frameHeaderSize + bodySize
		if framesSize < 0 {
			return ErrBodyOverflow
		}

		if isParseFramesProvided && !parseableIDs[id] {
			if _, err := io.CopyN(io.Discard, tag.reader, bodySize); err != nil {
				return err
			}
			continue
		}

		body := make([]byte, bodySize)
		if _, err := tag.reader.Read(body); err != nil {
			return err
		}

		frame, err := parseFrameBody(id, newBufReader(bytes.NewReader(body)), tag.version)
		if err != nil && err != io.EOF {
			return err
		}

		tag.AddFrame(id, frame)

		if isParseFramesProvided && !mustFrameBeInSequence(id) {
			delete(parseableIDs, id)

			// If it was last ID in parseIDs, we don't need to parse
			// other frames, so end the parsing.
			if len(parseableIDs) == 0 {
				break
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

func (tag *Tag) makeIDsFromDescriptions(parseFrames []string) map[string]bool {
	ids := make(map[string]bool, len(parseFrames))

	for _, description := range parseFrames {
		ids[tag.CommonID(description)] = true
	}

	return ids
}

func parseFrameHeader(buf []byte, rd io.Reader, synchSafe bool) (frameHeader, error) {
	var header frameHeader

	if len(buf) < frameHeaderSize {
		return header, errors.New("parseFrameHeader: buf is smaller than frame header size")
	}

	fhBuf := buf[:frameHeaderSize]
	if _, err := rd.Read(fhBuf); err != nil {
		return header, err
	}

	id := string(fhBuf[:4])
	bodySize, err := parseSize(fhBuf[4:8], synchSafe)
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

func parseFrameBody(id string, br *bufReader, version byte) (Framer, error) {
	if id[0] == 'T' && id != "TXXX" {
		return parseTextFrame(br.ReadAll())
	}

	if parseFunc, exists := parsers[id]; exists {
		return parseFunc(br, version)
	}

	return parseUnknownFrame(br)
}
