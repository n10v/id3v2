package id3v2

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	maskOrdered  = byte(1 << 0)
	maskToplevel = byte(1 << 1)
)

var ErrUnexpectedId = errors.New("unexpected ID")

type ChapterTocFrame struct {
	ElementID string
	// This frame is the root of the Table of Contents tree and is not a child of any other "CTOC" frame.
	TopLevel bool
	// This provides a hint as to whether the elements should be played as a continuous ordered sequence or played individually.
	Ordered     bool
	ChapterIds  []string
	Description *TextFrame
}

func (ctf ChapterTocFrame) Size() int {
	size := encodedSize(ctf.ElementID, EncodingISO)
	size += 1 // trailing zero after ElementID
	size += 1 // CTOC Flags
	// The Entry count is the number of entries in the Child Element ID
	// list that follows and must be greater than zero.
	size += 1 // Entrycount

	// entries
	for _, id := range ctf.ChapterIds {
		size += encodedSize(id, EncodingISO)
		size += 1 // trailing zero after ID
	}

	// (optional) descriptive data
	if ctf.Description != nil {
		size += frameHeaderSize // Description frame header size
		size += ctf.Description.Size()
	}

	return size
}

func (ctf ChapterTocFrame) UniqueIdentifier() string {
	return ctf.ElementID
}

func (ctf ChapterTocFrame) WriteTo(w io.Writer) (n int64, err error) {
	return useBufWriter(w, func(bw *bufWriter) {
		bw.EncodeAndWriteText(ctf.ElementID, EncodingISO)
		bw.WriteByte(0)

		ctocFlags := byte(0)
		if ctf.TopLevel {
			ctocFlags |= maskToplevel
		}
		if ctf.Ordered {
			ctocFlags |= maskOrdered
		}

		binary.Write(bw, binary.BigEndian, ctocFlags)

		binary.Write(bw, binary.BigEndian, uint8(len(ctf.ChapterIds)))

		for _, id := range ctf.ChapterIds {
			bw.EncodeAndWriteText(id, EncodingISO)
			bw.WriteByte(0)
		}

		if ctf.Description != nil {
			writeFrame(bw, "TIT2", *ctf.Description, true)
		}
	})
}

func parseChapterTocFrame(br *bufReader, version byte) (Framer, error) {
	elementID := string(br.ReadText(EncodingISO))
	synchSafe := version == 4
	var ctocFlags byte
	if err := binary.Read(br, binary.BigEndian, &ctocFlags); err != nil {
		return nil, err
	}

	var elements uint8
	if err := binary.Read(br, binary.BigEndian, &elements); err != nil {
		return nil, err
	}

	chaptersIDs := make([]string, elements)
	for i := uint8(0); i < elements; i++ {
		chaptersIDs[i] = string(br.ReadText(EncodingISO))
	}

	var description TextFrame

	// borrowed from parse.go
	buf := getByteSlice(32 * 1024)
	defer putByteSlice(buf)

	for {
		header, err := parseFrameHeader(buf, br, synchSafe)
		if err == io.EOF || err == errBlankFrame || err == ErrInvalidSizeFormat {
			break
		}

		if err != nil {
			return nil, err
		}

		if header.ID != "TIT2" {
			return nil, fmt.Errorf("expected: '%s', got: '%s'  : %w", "TIT2", header.ID, ErrUnexpectedId)
		}

		bodyRd := getLimitedReader(br, header.BodySize)
		br := newBufReader(bodyRd)
		frame, err := parseTextFrame(br)
		if err != nil {
			putLimitedReader(bodyRd)
			return nil, err
		}
		description = frame.(TextFrame)

		putLimitedReader(bodyRd)
	}

	tocFrame := ChapterTocFrame{
		ElementID:   elementID,
		TopLevel:    (ctocFlags & maskToplevel) == maskToplevel,
		Ordered:     (ctocFlags & maskOrdered) == maskOrdered,
		ChapterIds:  chaptersIDs,
		Description: &description,
	}

	return tocFrame, nil
}
