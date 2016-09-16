// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/bogem/id3v2/bbpool"
	"github.com/bogem/id3v2/util"
)

// Tag stores all frames of opened file.
type Tag struct {
	frames    map[string]Framer
	sequences map[string]sequencer
	ids       map[string]string

	file         *os.File
	originalSize int64
	version      byte
}

func (t *Tag) AddFrame(id string, f Framer) {
	switch id {
	case t.ID("Attached picture"):
		pf := f.(PictureFrame)
		t.AddAttachedPicture(pf)
	case t.ID("Comments"):
		cf := f.(CommentFrame)
		t.AddCommentFrame(cf)
	case t.ID("Unsynchronised lyrics/text transcription"):
		uslf := f.(UnsynchronisedLyricsFrame)
		t.AddUnsynchronisedLyricsFrame(uslf)
	default:
		t.frames[id] = f
	}
}

func (t *Tag) AddAttachedPicture(pf PictureFrame) {
	id := t.ID("Attached picture")
	t.checkExistenceOfSequence(id, newPictureSequence)
	t.addFrameToSequence(id, pf)
}

func (t *Tag) AddCommentFrame(cf CommentFrame) {
	id := t.ID("Comments")
	t.checkExistenceOfSequence(id, newCommentSequence)
	t.addFrameToSequence(id, cf)
}

func (t *Tag) AddUnsynchronisedLyricsFrame(uslf UnsynchronisedLyricsFrame) {
	id := t.ID("Unsynchronised lyrics/text transcription")
	t.checkExistenceOfSequence(id, newUSLFSequence)
	t.addFrameToSequence(id, uslf)
}

func (t *Tag) checkExistenceOfSequence(id string, newSequence func() sequencer) {
	if t.sequences[id] == nil {
		t.sequences[id] = newSequence()
	}
}

func (t *Tag) addFrameToSequence(id string, f Framer) {
	t.sequences[id].AddFrame(f)
}

// ID returns ID3v2.3 or ID3v2.4 (in appropriate to version of Tag) frame ID
// from given description.
// For example, ID("Language") will return "TLAN".
func (t Tag) ID(description string) string {
	return t.ids[description]
}

// AllFrames returns map, that contains all frames in tag, that could be parsed.
// The key of this map is an ID of frame and value is an array of frames.
func (t *Tag) AllFrames() map[string][]Framer {
	frames := make(map[string][]Framer)

	for id, frame := range t.frames {
		frames[id] = append(frames[id], frame)
	}

	for id, sequence := range t.sequences {
		frames[id] = append(frames[id], sequence.Frames()...)
	}

	return frames
}

// DeleteAllFrames deletes all frames in tag.
func (t *Tag) DeleteAllFrames() {
	t.frames = make(map[string]Framer)
	t.sequences = make(map[string]sequencer)
}

// DeleteFrames deletes frames in tag with given id.
func (t *Tag) DeleteFrames(id string) {
	delete(t.frames, id)
	delete(t.sequences, id)
}

// GetLastFrame returns last frame from slice, that is returned from GetFrames function.
// GetLastFrame is suitable for frames, that can be only one in whole tag.
// For example, for text frames.
//
// Example of usage:
//	bpmFramer := t.GetLastFrame(t.ID("BPM"))
//	if bpmFramer != nil {
//		bpm, ok := bpmFramer.(id3v2.TextFrame)
//		if !ok {
//			log.Fatal("Couldn't assert bpm frame")
//		}
//		fmt.Println(bpm.Text)
//	}
func (t *Tag) GetLastFrame(id string) Framer {
	fs := t.GetFrames(id)
	if len(fs) == 0 || fs == nil {
		return nil
	}
	return fs[len(fs)-1]
}

// GetFrames returns frames with corresponding id.
//
// Example of usage:
//	pictures := tag.GetFrames(tag.ID("Attached picture"))
//	if pictures != nil {
//		for _, f := range pictures {
//			pic, ok := f.(id3v2.PictureFrame)
//			if !ok {
//				log.Fatal("Couldn't assert picture frame")
//			}
//
//			// Do some operations with picture frame:
//			fmt.Println(pic.Description) // For example, print description of picture frame
//			image, err := ioutil.ReadAll(pic.Picture) // Or read a picture from picture frame
//			if err != nil {
//				log.Fatal("Error while reading a picture from picture frame: ", err)
//			}
//		}
//	}
func (t *Tag) GetFrames(id string) []Framer {
	if f, exists := t.frames[id]; exists {
		return []Framer{f}
	} else if s, exists := t.sequences[id]; exists {
		return s.Frames()
	}

	return nil
}

// GetTextFrame returns text frame with corresponding id.
func (t Tag) GetTextFrame(id string) TextFrame {
	f := t.GetLastFrame(id)
	if f == nil {
		return TextFrame{}
	}
	tf := f.(TextFrame)
	return tf
}

func (t Tag) Title() string {
	f := t.GetTextFrame(t.ID("Title/Songname/Content description"))
	return f.Text
}

func (t *Tag) SetTitle(title string) {
	t.AddFrame(t.ID("Title/Songname/Content description"), TextFrame{Encoding: ENUTF8, Text: title})
}

func (t Tag) Artist() string {
	f := t.GetTextFrame(t.ID("Lead artist/Lead performer/Soloist/Performing group"))
	return f.Text
}

func (t *Tag) SetArtist(artist string) {
	t.AddFrame(t.ID("Lead artist/Lead performer/Soloist/Performing group"), TextFrame{Encoding: ENUTF8, Text: artist})
}

func (t Tag) Album() string {
	f := t.GetTextFrame(t.ID("Album/Movie/Show title"))
	return f.Text
}

func (t *Tag) SetAlbum(album string) {
	t.AddFrame(t.ID("Album/Movie/Show title"), TextFrame{Encoding: ENUTF8, Text: album})
}

func (t Tag) Year() string {
	f := t.GetTextFrame(t.ID("Year"))
	return f.Text
}

func (t *Tag) SetYear(year string) {
	t.AddFrame(t.ID("Year"), TextFrame{Encoding: ENUTF8, Text: year})
}

func (t Tag) Genre() string {
	f := t.GetTextFrame(t.ID("Content type"))
	return f.Text
}

func (t *Tag) SetGenre(genre string) {
	t.AddFrame(t.ID("Content type"), TextFrame{Encoding: ENUTF8, Text: genre})
}

// Save writes tag to the file.
func (t *Tag) Save() error {
	// Forming new frames
	frames := t.formAllFrames()

	// Forming size of new frames
	framesSize := util.FormSize(int64(len(frames)))

	// Creating a temp file for mp3 file, which will contain new tag
	newFile, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}

	// Make sure we clean up the temp file if it is still around.
	defer os.Remove(newFile.Name())

	// Writing to new file new tag header
	if _, err = newFile.Write(formTagHeader(framesSize, t.version)); err != nil {
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
	t.file = newFile

	return nil
}

// Close closes the tag's file, rendering it unusable for I/O.
// It returns an error, if any.
func (t *Tag) Close() error {
	return t.file.Close()
}

func (t Tag) formAllFrames() []byte {
	framesBuffer := bbpool.Get()
	defer bbpool.Put(framesBuffer)

	t.writeFrames(framesBuffer)

	return framesBuffer.Bytes()
}

func (t Tag) writeFrames(w io.Writer) {
	for id, frames := range t.AllFrames() {
		for _, f := range frames {
			w.Write(formFrame(id, f))
		}
	}
}

func formFrame(id string, frame Framer) []byte {
	if id == "" {
		panic("there is blank ID in frames")
	}

	frameBuffer := bbpool.Get()
	defer bbpool.Put(frameBuffer)

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
