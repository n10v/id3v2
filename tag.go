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
	frames    map[string]frame.Framer
	sequences map[string]frame.Sequencer
	commonIDs map[string]string

	File         *os.File
	OriginalSize uint32
}

func (t *Tag) AddAttachedPicture(pf frame.PictureFramer) {
	var f frame.PictureSequencer
	id := t.commonIDs["Attached Picture"]

	existingSequence := t.sequences[id]
	if existingSequence == nil {
		f = frame.NewPictureSequence()
	} else {
		f = existingSequence.(frame.PictureSequencer)
	}

	f.AddPicture(pf)
	t.sequences[id] = f
}

func (t *Tag) AddTextFrame(id string, text string) {
	var f frame.TextFramer

	existingFrame := t.frames[id]
	if existingFrame == nil {
		f = new(frame.TextFrame)
	} else {
		f = existingFrame.(frame.TextFramer)
	}

	f.SetText(text)
	t.frames[id] = f
}

func (t *Tag) AddUnsynchronisedLyricsFrame(uslf frame.UnsynchronisedLyricsFramer) {
	var f frame.USLFSequencer
	id := t.commonIDs["Unsynchronised Lyrics/Text"]

	existingSequence := t.sequences[id]
	if existingSequence == nil {
		f = frame.NewUSLFSequence()
	} else {
		f = existingSequence.(frame.USLFSequencer)
	}

	f.AddUSLF(uslf)
	t.sequences[id] = f
}

func (t *Tag) AddCommentFrame(cf frame.CommentFramer) {
	var f frame.CommentSequencer
	id := t.commonIDs["Comment"]

	existingSequence := t.sequences[id]
	if existingSequence == nil {
		f = frame.NewCommentSequence()
	} else {
		f = existingSequence.(frame.CommentSequencer)
	}

	f.AddComment(cf)
	t.sequences[id] = f
}

func (t *Tag) SetTitle(title string) {
	t.AddTextFrame(t.commonIDs["Title"], title)
}

func (t *Tag) SetArtist(artist string) {
	t.AddTextFrame(t.commonIDs["Artist"], artist)
}

func (t *Tag) SetAlbum(album string) {
	t.AddTextFrame(t.commonIDs["Album"], album)
}

func (t *Tag) SetYear(year string) {
	t.AddTextFrame(t.commonIDs["Year"], year)
}

func (t *Tag) SetGenre(genre string) {
	t.AddTextFrame(t.commonIDs["Genre"], genre)
}

func NewTag(file *os.File) *Tag {
	return &Tag{
		frames:    make(map[string]frame.Framer),
		sequences: make(map[string]frame.Sequencer),
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
		frames:    make(map[string]frame.Framer),
		sequences: make(map[string]frame.Sequencer),
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
	newFileBuf := bufio.NewWriter(newFile)

	// Writing to new file new tag header
	if _, err := newFileBuf.Write(FormTagHeader()); err != nil {
		return err
	}

	// Writing to new file new frames
	// And getting size of them
	frames, err := t.FormAllFrames()
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

	// Flushing the buffered data to new file
	if err = newFileBuf.Flush(); err != nil {
		return err
	}

	// Setting size of frames to new file
	if err = setSize(newFile, uint32(framesSize)); err != nil {
		return err
	}

	// Replacing original file with new file
	if err = os.Rename(newFile.Name(), f.Name()); err != nil {
		return err
	}

	// And closing it
	if err = newFile.Close(); err != nil {
		return err
	}

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

func (t Tag) FormAllFrames() ([]byte, error) {
	frames := bytesBufPool.Get().(*bytes.Buffer)
	frames.Reset()

	if f, err := t.FormFrames(); err != nil {
		return nil, err
	} else {
		frames.Write(f)
	}

	if s, err := t.FormSequences(); err != nil {
		return nil, err
	} else {
		frames.Write(s)
	}

	bytesBufPool.Put(frames)
	return frames.Bytes(), nil
}

func (t Tag) FormFrames() ([]byte, error) {
	frames := bytesBufPool.Get().(*bytes.Buffer)
	frames.Reset()

	for id, f := range t.frames {
		frameBody, err := f.Bytes()
		if err != nil {
			return nil, err
		}
		frameHeader, err := formFrameHeader(id, uint32(len(frameBody)))
		if err != nil {
			return nil, err
		}
		frames.Write(frameHeader)
		frames.Write(frameBody)
	}

	bytesBufPool.Put(frames)
	return frames.Bytes(), nil
}

func (t Tag) FormSequences() ([]byte, error) {
	frames := bytesBufPool.Get().(*bytes.Buffer)
	frames.Reset()

	for id, s := range t.sequences {
		for _, f := range s.Frames() {
			frameBody, err := f.Bytes()
			if err != nil {
				return nil, err
			}
			frameHeader, err := formFrameHeader(id, uint32(len(frameBody)))
			if err != nil {
				return nil, err
			}
			frames.Write(frameHeader)
			frames.Write(frameBody)
		}
	}

	bytesBufPool.Put(frames)
	return frames.Bytes(), nil
}

func formFrameHeader(id string, frameSize uint32) ([]byte, error) {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()

	b.WriteString(id)
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
