// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"errors"
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
	commonIDs map[string]string

	file         *os.File
	originalSize int64
	version      byte
}

func (t *Tag) AddFrame(id string, f Framer) {
	switch id {
	case t.CommonID("Attached picture"):
		pf := f.(PictureFrame)
		t.AddAttachedPicture(pf)
	case t.CommonID("Comments"):
		cf := f.(CommentFrame)
		t.AddCommentFrame(cf)
	case t.CommonID("Unsynchronised lyrics/text transcription"):
		uslf := f.(UnsynchronisedLyricsFrame)
		t.AddUnsynchronisedLyricsFrame(uslf)
	default:
		t.frames[id] = f
	}
}

func (t *Tag) AddAttachedPicture(pf PictureFrame) {
	id := t.CommonID("Attached picture")
	t.checkExistenceOfSequence(id, newPictureSequence)
	t.addFrameToSequence(id, pf)
}

func (t *Tag) AddCommentFrame(cf CommentFrame) {
	id := t.CommonID("Comments")
	t.checkExistenceOfSequence(id, newCommentSequence)
	t.addFrameToSequence(id, cf)
}

func (t *Tag) AddUnsynchronisedLyricsFrame(uslf UnsynchronisedLyricsFrame) {
	id := t.CommonID("Unsynchronised lyrics/text transcription")
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
// For example, CommonID("Language") will return "TLAN".
func (t Tag) CommonID(description string) string {
	return t.commonIDs[description]
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
//	bpmFramer := tag.GetLastFrame(tag.CommonID("BPM"))
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
//	pictures := tag.GetFrames(tag.CommonID("Attached picture"))
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
	f := t.GetTextFrame(t.CommonID("Title/Songname/Content description"))
	return f.Text
}

func (t *Tag) SetTitle(title string) {
	t.AddFrame(t.CommonID("Title/Songname/Content description"), TextFrame{Encoding: ENUTF8, Text: title})
}

func (t Tag) Artist() string {
	f := t.GetTextFrame(t.CommonID("Lead artist/Lead performer/Soloist/Performing group"))
	return f.Text
}

func (t *Tag) SetArtist(artist string) {
	t.AddFrame(t.CommonID("Lead artist/Lead performer/Soloist/Performing group"), TextFrame{Encoding: ENUTF8, Text: artist})
}

func (t Tag) Album() string {
	f := t.GetTextFrame(t.CommonID("Album/Movie/Show title"))
	return f.Text
}

func (t *Tag) SetAlbum(album string) {
	t.AddFrame(t.CommonID("Album/Movie/Show title"), TextFrame{Encoding: ENUTF8, Text: album})
}

func (t Tag) Year() string {
	f := t.GetTextFrame(t.CommonID("Year"))
	return f.Text
}

func (t *Tag) SetYear(year string) {
	t.AddFrame(t.CommonID("Year"), TextFrame{Encoding: ENUTF8, Text: year})
}

func (t Tag) Genre() string {
	f := t.GetTextFrame(t.CommonID("Content type"))
	return f.Text
}

func (t *Tag) SetGenre(genre string) {
	t.AddFrame(t.CommonID("Content type"), TextFrame{Encoding: ENUTF8, Text: genre})
}

// Save writes tag to the file.
func (t *Tag) Save() error {
	// Create a temp file for mp3 file, which will contain new tag
	newFile, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}

	// Make sure we clean up the temp file if it's still around
	defer os.Remove(newFile.Name())

	// If there is at least one frame, write it
	if len(t.frames) > 0 || len(t.sequences) > 0 {
		// Form new frames
		frames, err := t.formAllFrames()
		if err != nil {
			return err
		}

		// Form size of new frames
		framesSize, err := util.FormSize(len(frames))
		if err != nil {
			return err
		}

		// Write to new file new tag header
		if _, err = newFile.Write(formTagHeader(framesSize, t.version)); err != nil {
			return err
		}

		// Write to new file new frames
		if _, err = newFile.Write(frames); err != nil {
			return err
		}
	}

	// Seek to a music part of original file
	originalFile := t.file
	if _, err = originalFile.Seek(t.originalSize, os.SEEK_SET); err != nil {
		return err
	}

	// Write to new file the music part
	if _, err = io.Copy(newFile, originalFile); err != nil {
		return err
	}

	// Get original file mode
	originalFileStat, err := originalFile.Stat()
	if err != nil {
		return err
	}
	originalFileMode := originalFileStat.Mode()

	// Set original file mode to new file
	if err = os.Chmod(newFile.Name(), originalFileMode); err != nil {
		return err
	}

	// Close files to allow replacing
	newFile.Close()
	originalFile.Close()

	// Replace original file with new file
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

var blankID = errors.New("blank ID")

func (t Tag) formAllFrames() ([]byte, error) {
	framesBuffer := bbpool.Get()
	defer bbpool.Put(framesBuffer)

	for id, frames := range t.AllFrames() {
		for _, f := range frames {
			formedFrame, err := formFrame(id, f)
			if err == blankID {
				continue
			}
			if err != nil {
				return nil, err
			}
			framesBuffer.Write(formedFrame)
		}
	}

	return framesBuffer.Bytes(), nil
}

func formFrame(id string, frame Framer) ([]byte, error) {
	if id == "" {
		return nil, blankID
	}

	frameBuffer := bbpool.Get()
	defer bbpool.Put(frameBuffer)

	frameBody := frame.Body()
	if err := writeFrameHeader(frameBuffer, id, len(frameBody)); err != nil {
		return nil, err
	}
	frameBuffer.Write(frameBody)

	return frameBuffer.Bytes(), nil
}

func writeFrameHeader(buf *bytes.Buffer, id string, frameSize int) error {
	size, err := util.FormSize(frameSize)
	if err != nil {
		return err
	}

	buf.WriteString(id)
	buf.Write(size)
	buf.Write([]byte{0, 0})
	return nil
}
