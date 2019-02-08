package id3v2

import (
	"encoding/binary"
	"io"
	"time"
)

const (
	nanosInMillis = 1000000
	IgnoredOffset = 0xFFFFFFFF
)

// ChapterFrame is used to work with CHAP frames
// according to spec from http://id3.org/id3v2-chapters-1.0
// This implementation only supports single TIT2 subframe (Title field).
// All other subframes are ignored.
// If StartOffset or EndOffset == id3v2.IgnoredOffset, then it should be ignored
// and StartTime or EndTime should be utilized
type ChapterFrame struct {
	ElementID   string
	StartTime   time.Duration
	EndTime     time.Duration
	StartOffset uint32
	EndOffset   uint32
	Title       string
}

func (cf ChapterFrame) Size() int {
	titleFrame := TextFrame{
		Encoding: EncodingUTF8,
		Text:     cf.Title,
	}
	return encodedSize(cf.ElementID, EncodingISO) +
		1 + // trailing zero after ElementID
		4 + 4 + 4 + 4 + // (Start, End) (Time, Offset)
		frameHeaderSize + // Title frame header size
		titleFrame.Size()
}

func (cf ChapterFrame) WriteTo(w io.Writer) (n int64, err error) {
	return useBufWriter(w, func(bw *bufWriter) {
		bw.EncodeAndWriteText(cf.ElementID, EncodingISO)
		bw.WriteByte(0)
		// nanoseconds => milliseconds
		binary.Write(bw, binary.BigEndian, int32(cf.StartTime/nanosInMillis))
		binary.Write(bw, binary.BigEndian, int32(cf.EndTime/nanosInMillis))

		binary.Write(bw, binary.BigEndian, cf.StartOffset)
		binary.Write(bw, binary.BigEndian, cf.EndOffset)

		titleFrame := TextFrame{
			Encoding: EncodingUTF8,
			Text:     cf.Title,
		}
		writeFrame(bw, "TIT2", titleFrame, true)
	})
}

func parseChapterFrame(br *bufReader) (Framer, error) {
	ElementID := br.ReadText(EncodingISO)
	chapterTime := make([]int32, 2)
	for i := range chapterTime {
		if err := binary.Read(br, binary.BigEndian, &chapterTime[i]); err != nil {
			return nil, err
		}
	}
	chapterOffset := make([]uint32, 2)
	for i := range chapterOffset {
		if err := binary.Read(br, binary.BigEndian, &chapterOffset[i]); err != nil {
			return nil, err
		}
	}
	var title string

	// borrowed from parse.go
	buf := getByteSlice(32 * 1024)
	defer putByteSlice(buf)
	for {
		// no way to determine whether this should be true or not
		// this is likely should be fixed
		header, err := parseFrameHeader(buf, br, true)
		if err == io.EOF || err == errBlankFrame || err == ErrInvalidSizeFormat {
			break
		}
		if err != nil {
			return nil, err
		}
		id, bodySize := header.ID, header.BodySize
		if id == "TIT2" {
			bodyRd := getLimitedReader(br, bodySize)
			br2 := newBufReader(bodyRd)
			frame, err := parseTextFrame(br2)
			if err != nil {
				putLimitedReader(bodyRd)
				return nil, err
			}
			title = frame.(TextFrame).Text

			putLimitedReader(bodyRd)
			break
		}
	}

	cf := ChapterFrame{
		ElementID: string(ElementID),
		// StartTime is given in milliseconds, so we should convert it to nanoseconds
		// for time.Duration
		StartTime:   time.Duration(chapterTime[0] * nanosInMillis),
		EndTime:     time.Duration(chapterTime[1] * nanosInMillis),
		StartOffset: chapterOffset[0],
		EndOffset:   chapterOffset[1],
		Title:       title,
	}
	return cf, nil
}
