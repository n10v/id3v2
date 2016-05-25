package id3v2

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/bogem/id3v2/frame"
	"github.com/bogem/id3v2/util"
	"io/ioutil"
	"os"
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

func (t *Tag) Flush() error {
	// Creating a temp file for mp3 file, which will contain new tag
	newFile, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer os.Remove(newFile.Name())
	newFileBuf := bufio.NewWriter(newFile)

	// Writing to new file new header
	tagHeader, err := FormTagHeader(*t.Header)
	if err != nil {
		return err
	}
	if _, err := newFileBuf.Write(tagHeader); err != nil {
		return err
	}

	// Writing to new file new frames
	// And getting size of new frames
	frames, err := t.FormFrames()
	if err != nil {
		return err
	}
	framesSize, err := newFileBuf.Write(frames)
	if err != nil {
		return err
	}

	// Getting a music part of mp3
	// (Discarding an original tag of mp3)
	f := t.File
	defer f.Close()
	if _, err := f.Seek(0, os.SEEK_SET); err != nil {
		return err
	}
	originalFileBuf := bufio.NewReader(f)
	if _, err := originalFileBuf.Discard(int(t.OriginalSize)); err != nil {
		return err
	}

	// Writing to new file the music part
	if _, err = newFileBuf.ReadFrom(originalFileBuf); err != nil {
		return err
	}

	// Flushing the buffered to new file
	if err = newFileBuf.Flush(); err != nil {
		return err
	}

	// Setting size of frames to new file
	setSize(newFile, uint32(framesSize))

	// Replacing original file with new file
	os.Rename(newFile.Name(), f.Name())
	return nil
}

func setSize(f *os.File, size uint32) (err error) {
	if _, err := f.Seek(0, os.SEEK_SET); err != nil {
		return err
	}

	sizeBytes, err := util.FormSize(size)
	if err != nil {
		return
	}

	if _, err = f.WriteAt(sizeBytes, 6); err != nil {
		return
	}

	return
}

func (t Tag) FormFrames() ([]byte, error) {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()

	for _, f := range t.frames {
		formedFrame := f.Form()
		frameHeader, err := formFrameHeader(f, uint32(len(formedFrame)))
		if err != nil {
			return nil, err
		}
		b.Write(frameHeader)
		b.Write(formedFrame)
	}

	bytesBufPool.Put(b)
	return b.Bytes(), nil
}

func formFrameHeader(f frame.Framer, frameSize uint32) ([]byte, error) {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()

	b.WriteString(f.ID())
	if size, err := util.FormSize(frameSize); err != nil {
		return nil, err
	} else {
		b.Write(size)
	}
	b.WriteByte(0)
	b.WriteByte(0)

	bytesBufPool.Put(b)
	return b.Bytes(), nil
}
