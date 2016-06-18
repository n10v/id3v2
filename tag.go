// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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
	frames       map[string]frame.Framer
	sequences    map[string]frame.Sequencer
	commonIDs    map[string]string
	file         *os.File
	originalSize uint32
}

func (t *Tag) AddFrame(id string, f frame.Framer) {
	if t.frames == nil {
		t.frames = make(map[string]frame.Framer)
	}
	t.frames[id] = f
}

func (t *Tag) AddAttachedPicture(pf frame.PictureFramer) {
	id := t.commonIDs["Attached picture"]
	t.checkExistenceOfSequence(id, frame.NewPictureSequence)
	t.addFrameToSequence(pf, id)
}

func (t *Tag) AddUnsynchronisedLyricsFrame(uslf frame.UnsynchronisedLyricsFramer) {
	id := t.commonIDs["USLT"]
	t.checkExistenceOfSequence(id, frame.NewUSLFSequence)
	t.addFrameToSequence(uslf, id)
}

func (t *Tag) AddCommentFrame(cf frame.CommentFramer) {
	id := t.commonIDs["Comment"]
	t.checkExistenceOfSequence(id, frame.NewCommentSequence)
	t.addFrameToSequence(cf, id)
}

func (t *Tag) checkExistenceOfSequence(id string, newSequence func() frame.Sequencer) {
	if t.sequences == nil {
		t.sequences = make(map[string]frame.Sequencer)
	}
	if t.sequences[id] == nil {
		t.sequences[id] = newSequence()
	}
}

func (t *Tag) addFrameToSequence(f frame.Framer, id string) {
	t.sequences[id].AddFrame(f)
}

func (t *Tag) SetTitle(title string) {
	t.AddFrame(t.commonIDs["Title"], NewTextFrame(title))
}

func (t *Tag) SetArtist(artist string) {
	t.AddFrame(t.commonIDs["Artist"], NewTextFrame(artist))
}

func (t *Tag) SetAlbum(album string) {
	t.AddFrame(t.commonIDs["Album"], NewTextFrame(album))
}

func (t *Tag) SetYear(year string) {
	t.AddFrame(t.commonIDs["Year"], NewTextFrame(year))
}

func (t *Tag) SetGenre(genre string) {
	t.AddFrame(t.commonIDs["Genre"], NewTextFrame(genre))
}

func newTag(file *os.File, size uint32) *Tag {
	return &Tag{
		commonIDs: frame.V24CommonIDs,

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
	// Creating a temp file for mp3 file, which will contain new tag
	newFile, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	newFileBuf := bufio.NewWriter(newFile)

	// Forming new frames
	frames, err := t.formAllFrames()
	if err != nil {
		return err
	}

	// Forming size of new frames
	framesSize, err := util.FormSize(uint32(len(frames)))
	if err != nil {
		return err
	}

	// Writing to new file new tag header
	if _, err := newFileBuf.Write(formTagHeader(framesSize)); err != nil {
		return err
	}

	// Writing to new file new frames
	if _, err := newFileBuf.Write(frames); err != nil {
		return err
	}

	// Getting a music part of mp3
	// (Discarding an original tag of mp3)
	f := t.file
	defer f.Close()
	if _, err := f.Seek(0, os.SEEK_SET); err != nil {
		return err
	}
	originalFileBuf := bufio.NewReader(f)
	if _, err := originalFileBuf.Discard(int(t.originalSize)); err != nil {
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

func (t Tag) formAllFrames() ([]byte, error) {
	frames := bytesBufPool.Get().(*bytes.Buffer)
	frames.Reset()

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

	bytesBufPool.Put(frames)
	return frames.Bytes(), nil
}

func (t Tag) formFrames() ([]byte, error) {
	frames := bytesBufPool.Get().(*bytes.Buffer)
	frames.Reset()

	for id, f := range t.frames {
		if id == "" {
			return nil, errors.New("Uncorrect ID in frames")
		}
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

func (t Tag) formSequences() ([]byte, error) {
	frames := bytesBufPool.Get().(*bytes.Buffer)
	frames.Reset()

	for id, s := range t.sequences {
		if id == "" {
			return nil, errors.New("Uncorrect ID in sequences")
		}
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
	size, err := util.FormSize(frameSize)
	if err != nil {
		return nil, err
	}
	b.Write(size)
	b.WriteByte(0)
	b.WriteByte(0)

	bytesBufPool.Put(b)
	return b.Bytes(), nil
}
