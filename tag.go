package id3v2

import (
	"bytes"
	"errors"
	"github.com/bogem/id3v2/frame"
	"os"
	"sync"
	"sync/atomic"
)

type Tag struct {
	Header    *TagHeader
	frames    map[string]frame.Framer
	commonIDs map[string]string

	File         *os.File
	OriginalSize uint32
}

func (t *Tag) SetAttachedPicture(pf frame.PictureFramer) {
	var f frame.APICSequencer
	id := t.commonIDs["Attached Picture"]

	existingFrame := t.frames[id]
	if existingFrame == nil {
		f = frame.NewAPICSequence()
		f.SetID(id)
	} else {
		f = existingFrame.(frame.APICSequencer)
	}

	f.AddPicture(pf)
	t.frames[id] = f
}

func (t *Tag) SetTextFrame(id string, text string) {
	var f frame.TextFramer

	existingFrame := t.frames[id]
	if existingFrame == nil {
		f = new(frame.TextFrame)
		f.SetID(id)
	} else {
		f = existingFrame.(frame.TextFramer)
	}

	f.SetText(text)
	t.frames[id] = f
}

func (t *Tag) SetTitle(title string) {
	t.SetTextFrame(t.commonIDs["Title"], title)
}

func (t *Tag) SetArtist(artist string) {
	t.SetTextFrame(t.commonIDs["Artist"], artist)
}

func (t *Tag) SetAlbum(album string) {
	t.SetTextFrame(t.commonIDs["Album"], album)
}

func (t *Tag) SetYear(year string) {
	t.SetTextFrame(t.commonIDs["Year"], year)
}

func (t *Tag) SetGenre(genre string) {
	t.SetTextFrame(t.commonIDs["Genre"], genre)
}

func NewTag(file *os.File) *Tag {
	header := &TagHeader{
		Version:    4,
		FramesSize: 0,
	}

	return &Tag{
		Header:    header,
		frames:    make(map[string]frame.Framer),
		commonIDs: frame.V24CommonIDs,

		File:         file,
		OriginalSize: 0,
	}
}

func ParseTag(file *os.File) (*Tag, error) {
	header, err := ParseHeader(file)
	if err != nil {
		err = errors.New("Trying to parse tag header: " + err.Error())
		return nil, err
	}
	if header == nil {
		return NewTag(file), nil
	}
	if header.Version < 3 {
		err = errors.New("Unsupported version of ID3 tag")
		return nil, err
	}

	tag := &Tag{
		Header:    header,
		frames:    make(map[string]frame.Framer),
		commonIDs: frame.V24CommonIDs,

		File:         file,
		OriginalSize: TagHeaderSize + header.FramesSize,
	}

	return tag, nil
}

func (t Tag) CalculateSizeOfAllFrames() uint32 {
	var wg sync.WaitGroup
	wg.Add(len(t.frames))

	var size uint32
	for _, v := range t.frames {
		go func(f frame.Framer) {
			atomic.AddUint32(&size, frame.FrameHeaderSize+f.Size())
			wg.Done()
		}(v)
	}
	wg.Wait()
	return atomic.LoadUint32(&size)
}

func (t Tag) FormFrames() ([]byte, error) {
	var b bytes.Buffer
	for _, frame := range t.frames {
		formedFrame, err := frame.Form()
		if err != nil {
			return nil, err
		}
		b.Write(formedFrame)
	}
	return b.Bytes(), nil
}
