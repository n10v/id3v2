// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/bogem/id3v2/bytesbufferpool"
	"github.com/bogem/id3v2/util"
)

// Tag stores all frames of opened file.
type Tag struct {
	framesCoords map[string][]frameCoordinates
	frames       map[string]Framer
	sequences    map[string]sequencer
	ids          map[string]string

	file         *os.File
	originalSize int64
}

func (t *Tag) AddFrame(id string, f Framer) {
	if addFunc := t.findSpecificAddFunction(id); addFunc != nil {
		addFunc(f)
	} else {
		t.frames[id] = f
	}
}

// ID returns ID3v2.3 or ID3v2.4 (in appropriate to version of Tag) frame ID
// from given description.
// For example, ID("Mood") will return "TMOO".
func (t Tag) ID(description string) string {
	return t.ids[description]
}

func (t *Tag) findSpecificAddFunction(id string) func(Framer) {
	switch id {
	case t.ids["Attached picture"]:
		return t.addAttachedPicture
	case t.ids["Comments"]:
		return t.addCommentFrame
	case t.ids["Unsynchronised lyrics/text transcription"]:
		return t.addUnsynchronisedLyricsFrame
	}
	return nil
}

func (t *Tag) addAttachedPicture(f Framer) {
	id := t.ids["Attached picture"]
	t.checkExistenceOfSequence(id, newPictureSequence)
	t.addFrameToSequence(id, f)
}

func (t *Tag) addCommentFrame(f Framer) {
	id := t.ids["Comments"]
	t.checkExistenceOfSequence(id, newCommentSequence)
	t.addFrameToSequence(id, f)
}

func (t *Tag) addUnsynchronisedLyricsFrame(f Framer) {
	id := t.ids["Unsynchronised lyrics/text transcription"]
	t.checkExistenceOfSequence(id, newUSLFSequence)
	t.addFrameToSequence(id, f)
}

func (t *Tag) checkExistenceOfSequence(id string, newSequence func() sequencer) {
	if t.sequences[id] == nil {
		t.sequences[id] = newSequence()
	}
}

func (t *Tag) addFrameToSequence(id string, f Framer) {
	t.sequences[id].AddFrame(f)
}

// AllFrames returns all frames in tag, that can be parsed.
func (t *Tag) AllFrames() map[string][]Framer {
	frames := make(map[string][]Framer)

	t.parseAllFramesCoords()

	// Add frames from t.frames
	for id, frame := range t.frames {
		frames[id] = append(frames[id], frame)
	}

	// Add frames from t.sequences
	for id, sequence := range t.sequences {
		frames[id] = append(frames[id], sequence.Frames()...)
	}

	return frames
}

// GetLastFrame returns last frame from slice, which is returned from GetFrames function.
func (t *Tag) GetLastFrame(id string) Framer {
	fs := t.GetFrames(id)
	if len(fs) == 0 || fs == nil {
		return nil
	}
	return fs[len(fs)-1]
}

// GetFrames returns frames with corresponding id.
func (t *Tag) GetFrames(id string) []Framer {
	// If frames with id didn't parsed yet, parse them
	if _, exists := t.framesCoords[id]; exists {
		t.parseFramesCoordsWithID(id)
	}

	if f, exists := t.frames[id]; exists {
		return []Framer{f}
	}

	if s, exists := t.sequences[id]; exists {
		return s.Frames()
	}

	return nil
}

// GetFrames returns text frame with corresponding id.
func (t Tag) GetTextFrame(id string) TextFrame {
	f := t.GetLastFrame(id)
	if f == nil {
		return TextFrame{}
	}
	tf := f.(TextFrame)
	return tf
}

func (t Tag) Title() string {
	f := t.GetTextFrame(t.ids["Title/Songname/Content description"])
	return f.Text
}

func (t *Tag) SetTitle(title string) {
	t.AddFrame(t.ids["Title/Songname/Content description"], TextFrame{Encoding: ENUTF8, Text: title})
}

func (t Tag) Artist() string {
	f := t.GetTextFrame(t.ids["Lead artist/Lead performer/Soloist/Performing group"])
	return f.Text
}

func (t *Tag) SetArtist(artist string) {
	t.AddFrame(t.ids["Lead artist/Lead performer/Soloist/Performing group"], TextFrame{Encoding: ENUTF8, Text: artist})
}

func (t Tag) Album() string {
	f := t.GetTextFrame(t.ids["Album/Movie/Show title"])
	return f.Text
}

func (t *Tag) SetAlbum(album string) {
	t.AddFrame(t.ids["Album/Movie/Show title"], TextFrame{Encoding: ENUTF8, Text: album})
}

func (t Tag) Year() string {
	f := t.GetTextFrame(t.ids["Recording time"])
	return f.Text
}

func (t *Tag) SetYear(year string) {
	t.AddFrame(t.ids["Recording time"], TextFrame{Encoding: ENUTF8, Text: year})
}

func (t Tag) Genre() string {
	f := t.GetTextFrame(t.ids["Content type"])
	return f.Text
}

func (t *Tag) SetGenre(genre string) {
	t.AddFrame(t.ids["Content type"], TextFrame{Encoding: ENUTF8, Text: genre})
}

// Flush writes tag to the file.
func (t Tag) Flush() error {
	// Forming new frames
	frames := t.formAllFrames()

	// Forming size of new frames
	framesSize := util.FormSize(int64(len(frames)))

	// Creating a temp file for mp3 file, which will contain new tag
	newFile, err := ioutil.TempFile("", "")
	defer newFile.Close()
	if err != nil {
		return err
	}

	// Writing to new file new tag header
	if _, err = newFile.Write(formTagHeader(framesSize)); err != nil {
		return err
	}

	// Writing to new file new frames
	if _, err = newFile.Write(frames); err != nil {
		return err
	}

	// Seeking to a music part of mp3
	originalFile := t.file
	defer originalFile.Close()
	if _, err = originalFile.Seek(t.originalSize, os.SEEK_SET); err != nil {
		return err
	}

	// Writing to new file the music part
	if _, err = io.Copy(newFile, originalFile); err != nil {
		return err
	}

	// Getting original file mode
	originalFileStat, err := originalFile.Stat()
	if err != nil {
		return err
	}
	originalFileMode := originalFileStat.Mode()

	// Setting new file mode
	if err = newFile.Chmod(originalFileMode); err != nil {
		return err
	}

	// Replacing original file with new file
	if err = os.Rename(newFile.Name(), originalFile.Name()); err != nil {
		return err
	}

	return nil
}

func (t Tag) formAllFrames() []byte {
	frames := bytesbufferpool.Get()
	defer bytesbufferpool.Put(frames)

	t.writeFrames(frames)
	t.writeSequences(frames)

	return frames.Bytes()
}

func (t Tag) writeFrames(w io.Writer) {
	for id, f := range t.frames {
		w.Write(formFrame(id, f))
	}
}

func (t Tag) writeSequences(w io.Writer) {
	for id, s := range t.sequences {
		for _, f := range s.Frames() {
			w.Write(formFrame(id, f))
		}
	}
}

func formFrame(id string, frame Framer) []byte {
	if id == "" {
		panic("there is blank ID in frames")
	}

	frameBuffer := bytesbufferpool.Get()
	defer bytesbufferpool.Put(frameBuffer)

	frameBody := frame.Body()
	writeFrameHeader(frameBuffer, id, int64(len(frameBody)))
	frameBuffer.Write(frameBody)

	return frameBuffer.Bytes()
}

func writeFrameHeader(buf *bytes.Buffer, id string, frameSize int64) {
	buf.WriteString(id)
	buf.Write(util.FormSize(frameSize))
	buf.Write([]byte{0, 0})
}
