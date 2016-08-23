package id3v2

import (
	"errors"
	"io"
	"os"

	"github.com/bogem/id3v2/util"
)

const frameHeaderSize = 10

type frameHeader struct {
	ID        string
	FrameSize uint32
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
		return newTag(file, 0), nil
	}
	if header.Version < 3 {
		err = errors.New("Unsupported version of ID3 tag")
		return nil, err
	}

	t := newTag(file, tagHeaderSize+header.FramesSize)
	err = t.findAllFrames()

	return t, nil
}

func newTag(file *os.File, originalSize uint32) *Tag {
	return &Tag{
		commonIDs: V24CommonIDs,

		file:         file,
		originalSize: originalSize,
	}
}

func (t *Tag) findAllFrames() error {
	if t.framesCoords == nil {
		t.framesCoords = make(map[string][]frameCoordinates)
	}

	pos := uint32(tagHeaderSize) // initial position of read - end of tag header (beginning of first frame)
	tagSize := t.originalSize
	f := t.file

	for pos < tagSize {
		if _, err := f.Seek(int64(pos), os.SEEK_SET); err != nil {
			return err
		}

		header, err := parseFrameHeader(f)
		if err != nil {
			return err
		}
		pos += frameHeaderSize

		fc := frameCoordinates{
			Len: int64(header.FrameSize),
			Pos: int64(pos),
		}
		fcs := t.framesCoords[header.ID]
		fcs = append(fcs, fc)
		t.framesCoords[header.ID] = fcs

		pos += header.FrameSize
	}

	return nil
}

func parseFrameHeader(rd io.Reader) (*frameHeader, error) {
	byteHeader := make([]byte, frameHeaderSize)
	n, err := rd.Read(byteHeader)
	if err != nil {
		return nil, err
	}
	if n < frameHeaderSize {
		err = errors.New("Size of frame header is less than expected")
	}

	header := &frameHeader{
		ID:        string(byteHeader[:4]),
		FrameSize: util.ParseSize(byteHeader[4:8]),
	}

	return header, nil

}

func (t Tag) findParseFunc(id string) func(io.Reader) (Framer, error) {
	if id[0] == 'T' {
		return ParseTextFrame
	}
	switch id {
	case t.commonIDs["Attached picture"]:
		return ParsePictureFrame
	case t.commonIDs["Comments"]:
		return ParseCommentFrame
	case t.commonIDs["Unsynchronised lyrics/text transcription"]:
		return ParseUnsynchronisedLyricsFrame
	}
	return nil
}

func readFrame(parseFunc func(io.Reader) (Framer, error), file *os.File, fc frameCoordinates) Framer {
	file.Seek(fc.Pos, os.SEEK_SET)
	rd := &io.LimitedReader{R: file, N: fc.Len}
	fr, err := parseFunc(rd)
	if err != nil {
		panic(err)
	}
	return fr
}

func (t *Tag) parseAllFramesCoords() {
	for id := range t.framesCoords {
		t.parseFramesCoordsWithID(id)
	}
}

func (t *Tag) parseFramesCoordsWithID(id string) {
	fcs, exists := t.framesCoords[id]
	if !exists {
		return
	}

	parseFunc := t.findParseFunc(id)
	if parseFunc != nil {
		for _, fc := range fcs {
			fr := readFrame(parseFunc, t.file, fc)
			t.AddFrame(id, fr)
		}
	}
	// Delete frames with id from t.framesCoords,
	// because they are just being parsed
	delete(t.framesCoords, id)
}
