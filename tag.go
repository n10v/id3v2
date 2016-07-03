// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/bogem/id3v2/bytesbufferpool"
	"github.com/bogem/id3v2/util"
)

type Tag struct {
	frames       map[string]Framer
	sequences    map[string]Sequencer
	commonIDs    map[string]string
	file         *os.File
	originalSize uint32
}

func (t *Tag) AddFrame(id string, f Framer) {
	if t.frames == nil {
		t.frames = make(map[string]Framer)
	}
	t.frames[id] = f
}

func (t *Tag) AddAttachedPicture(pf PictureFrame) {
	id := t.commonIDs["Attached picture"]
	t.checkExistenceOfSequence(id, NewPictureSequence)
	t.addFrameToSequence(pf, id)
}

func (t *Tag) AddUnsynchronisedLyricsFrame(uslf UnsynchronisedLyricsFrame) {
	id := t.commonIDs["USLT"]
	t.checkExistenceOfSequence(id, NewUSLFSequence)
	t.addFrameToSequence(uslf, id)
}

func (t *Tag) AddCommentFrame(cf CommentFrame) {
	id := t.commonIDs["Comment"]
	t.checkExistenceOfSequence(id, NewCommentSequence)
	t.addFrameToSequence(cf, id)
}

func (t *Tag) checkExistenceOfSequence(id string, newSequence func() Sequencer) {
	if t.sequences == nil {
		t.sequences = make(map[string]Sequencer)
	}
	if t.sequences[id] == nil {
		t.sequences[id] = newSequence()
	}
}

func (t *Tag) addFrameToSequence(f Framer, id string) {
	t.sequences[id].AddFrame(f)
}

func (t *Tag) SetTitle(title string) {
	t.AddFrame(t.commonIDs["Title"], TextFrame{Encoding: ENUTF8, Text: title})
}

func (t *Tag) SetArtist(artist string) {
	t.AddFrame(t.commonIDs["Artist"], TextFrame{Encoding: ENUTF8, Text: artist})
}

func (t *Tag) SetAlbum(album string) {
	t.AddFrame(t.commonIDs["Album"], TextFrame{Encoding: ENUTF8, Text: album})
}

func (t *Tag) SetYear(year string) {
	t.AddFrame(t.commonIDs["Year"], TextFrame{Encoding: ENUTF8, Text: year})
}

func (t *Tag) SetGenre(genre string) {
	t.AddFrame(t.commonIDs["Genre"], TextFrame{Encoding: ENUTF8, Text: genre})
}

func newTag(file *os.File, size uint32) *Tag {
	return &Tag{
		commonIDs: V24CommonIDs,

		file:         file,
		originalSize: size,
	}
}

func parseTag(file *os.File) (*Tag, error) {
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

	return newTag(file, tagHeaderSize+header.FramesSize), nil
}

// Flush writes tag to the file.
func (t Tag) Flush() error {
	// Forming new frames
	frames, err := t.formAllFrames()
	if err != nil {
		return err
	}

	// Forming size of new frames
	framesSize := util.FormSize(uint32(len(frames)))

	// Creating a temp file for mp3 file, which will contain new tag
	newFile, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}

	// Writing to new file new tag header
	if _, err := newFile.Write(formTagHeader(framesSize)); err != nil {
		return err
	}

	// Writing to new file new frames
	if _, err := newFile.Write(frames); err != nil {
		return err
	}

	// Seeking to a music part of mp3
	originalFile := t.file
	defer originalFile.Close()
	if _, err := originalFile.Seek(int64(t.originalSize), os.SEEK_SET); err != nil {
		return err
	}

	// Writing to new file the music part
	if _, err = io.Copy(newFile, originalFile); err != nil {
		return err
	}

	// Replacing original file with new file
	if err = os.Rename(newFile.Name(), originalFile.Name()); err != nil {
		return err
	}

	// And closing it
	if err = newFile.Close(); err != nil {
		return err
	}

	return nil
}

func (t Tag) formAllFrames() ([]byte, error) {
	frames := bytesbufferpool.Get()
	defer bytesbufferpool.Put(frames)

	f, err := t.formFrames()
	if err != nil {
		return nil, err
	}
	frames.Write(f)

	s, err := t.formSequences()
	if err != nil {
		return nil, err
	}
	frames.Write(s)

	return frames.Bytes(), nil
}

func (t Tag) formFrames() ([]byte, error) {
	frames := bytesbufferpool.Get()
	defer bytesbufferpool.Put(frames)

	for id, f := range t.frames {
		if id == "" {
			return nil, errors.New("Uncorrect ID in frames")
		}
		frameBody, err := f.Bytes()
		if err != nil {
			return nil, err
		}
		frameHeader := formFrameHeader(id, uint32(len(frameBody)))
		frames.Write(frameHeader)
		frames.Write(frameBody)
	}

	return frames.Bytes(), nil
}

func (t Tag) formSequences() ([]byte, error) {
	frames := bytesbufferpool.Get()
	defer bytesbufferpool.Put(frames)

	for id, s := range t.sequences {
		if id == "" {
			return nil, errors.New("Uncorrect ID in sequences")
		}
		for _, f := range s.Frames() {
			frameBody, err := f.Bytes()
			if err != nil {
				return nil, err
			}
			frameHeader := formFrameHeader(id, uint32(len(frameBody)))
			frames.Write(frameHeader)
			frames.Write(frameBody)
		}
	}

	return frames.Bytes(), nil
}

func formFrameHeader(id string, frameSize uint32) []byte {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.WriteString(id)
	b.Write(util.FormSize(frameSize))
	b.WriteByte(0)
	b.WriteByte(0)

	return b.Bytes()
}
