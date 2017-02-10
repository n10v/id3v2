// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bogem/id3v2/bwpool"
	"github.com/bogem/id3v2/util"
)

// Tag stores all frames of opened file.
type Tag struct {
	frames    map[string]Framer
	sequences map[string]sequencer

	file         *os.File
	originalSize int64
	version      byte
}

// AddFrame adds f to tag with appropriate id. If id is "" or f is nil,
// AddFrame will not add f to tag.
//
// If you want to add attached picture, comment or unsynchronised lyrics/text
// transcription frames, better use AddAttachedPicture, AddCommentFrame
// or AddUnsynchronisedLyricsFrame methods respectively.
func (t *Tag) AddFrame(id string, f Framer) {
	if id == "" || f == nil {
		return
	}

	if _, exists := sequenceConstructors[id]; exists { // check if frame should be added in sequence
		if t.sequences[id] == nil {
			constructor := sequenceConstructors[id]
			t.sequences[id] = constructor()
		}
		t.sequences[id].AddFrame(f)
	} else {
		t.frames[id] = f
	}
}

func (t *Tag) AddAttachedPicture(pf PictureFrame) {
	id := t.CommonID("Attached picture")
	t.AddFrame(id, pf)
}

func (t *Tag) AddCommentFrame(cf CommentFrame) {
	id := t.CommonID("Comments")
	t.AddFrame(id, cf)
}

func (t *Tag) AddUnsynchronisedLyricsFrame(uslf UnsynchronisedLyricsFrame) {
	id := t.CommonID("Unsynchronised lyrics/text transcription")
	t.AddFrame(id, uslf)
}

// CommonID returns ID3v2.3 or ID3v2.4 (in appropriate to version of Tag) frame ID
// from given description.
// For example, CommonID("Language") will return "TLAN".
// All descriptions you can find in file common_ids.go or in id3 documentation (for fourth version: http://id3.org/id3v2.4.0-frames; for third version: http://id3.org/id3v2.3.0#Declared_ID3v2_frames).
func (t Tag) CommonID(description string) string {
	if t.version == 3 {
		return V23CommonIDs[description]
	}
	return V24CommonIDs[description]
}

// AllFrames returns map, that contains all frames in tag, that could be parsed.
// The key of this map is an ID of frame and value is an array of frames.
func (t *Tag) AllFrames() map[string][]Framer {
	frames := make(map[string][]Framer)

	for id, f := range t.frames {
		frames[id] = []Framer{f}
	}
	for id, sequence := range t.sequences {
		frames[id] = sequence.Frames()
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
	// Avoid an allocation of slice in GetFrames,
	// if there is anyway one frame.
	if f, exists := t.frames[id]; exists {
		return f
	}

	fs := t.GetFrames(id)
	if len(fs) == 0 || fs == nil {
		return nil
	}
	return fs[len(fs)-1]
}

// GetFrames returns frames with corresponding id.
// It returns nil if there is no frames with given id.
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

// Count returns the number of frames in tag.
func (t Tag) Count() int {
	n := len(t.frames)
	for _, s := range t.sequences {
		n += s.Count()
	}
	return n
}

// HasAnyFrames checks if there is at least one frame in tag.
// It's much faster than tag.Count() > 0.
func (t Tag) HasAnyFrames() bool {
	return len(t.frames) > 0 || len(t.sequences) > 0
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

// iterateOverAllFrames iterates over every single frame in tag and call
// f for them. It consumps no memory at all, unlike the tag.AllFrames().
// It returns error only if f returns error.
func (t Tag) iterateOverAllFrames(f func(id string, frame Framer) error) error {
	for id, frame := range t.frames {
		if err := f(id, frame); err != nil {
			return err
		}
	}
	for id, sequence := range t.sequences {
		for _, frame := range sequence.Frames() {
			if err := f(id, frame); err != nil {
				return err
			}
		}
	}
	return nil
}

// Size returns the size of all ID3 tag in bytes.
func (t Tag) Size() int {
	if !t.HasAnyFrames() {
		return 0
	}

	var n int
	n += tagHeaderSize // Add the size of tag header
	t.iterateOverAllFrames(func(id string, f Framer) error {
		n += frameHeaderSize + f.Size()
		return nil
	})

	return n
}

// Version returns current ID3v2 version of tag.
func (t Tag) Version() byte {
	return t.version
}

// SetVersion sets given ID3v2 version to tag.
// If version is less than 3 or more than 4, then this method will do nothing.
// If tag has some frames, which are deprecated or changed in given version,
// then to your notice you can delete, change or just stay them.
func (t *Tag) SetVersion(version byte) {
	if version < 3 || version > 4 {
		return
	}
	t.version = version
}

// Save writes tag to the file. If there are no frames in tag, Save will write
// only music part without any ID3v2 information.
func (t *Tag) Save() error {
	// Create a temp file for mp3 file, which will contain new tag
	newFile, err := ioutil.TempFile(filepath.Dir(t.file.Name()), "id3v2-")
	if err != nil {
		return err
	}

	// Make sure we clean up the temp file if it's still around
	defer os.Remove(newFile.Name())

	// If there is at least one frame, write whole tag in new file
	var framesSize int
	if t.HasAnyFrames() {
		tagSize, err := t.Write(newFile)
		if err != nil {
			return err
		}
		framesSize = tagSize - tagHeaderSize
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

	// Set t.originalSize to new frames size
	t.originalSize = int64(framesSize)

	// Get original file mode
	originalFileStat, err := originalFile.Stat()
	if err != nil {
		return err
	}

	// Set original file mode to new file
	if err = os.Chmod(newFile.Name(), originalFileStat.Mode()); err != nil {
		return err
	}

	// Close files to allow replacing
	newFile.Close()
	originalFile.Close()

	// Replace original file with new file
	if err = os.Rename(newFile.Name(), originalFile.Name()); err != nil {
		return err
	}

	// Set t.file to new file with original name
	t.file, err = os.Open(originalFile.Name())
	if err != nil {
		return err
	}

	return nil
}

// Write writes whole tag in w.
func (t Tag) Write(w io.Writer) (n int, err error) {
	// Form size of frames
	framesSize := t.Size() - tagHeaderSize
	byteFramesSize, err := util.FormSize(framesSize)
	if err != nil {
		return 0, err
	}

	// Write tag header
	if _, err = w.Write(formTagHeader(byteFramesSize, t.version)); err != nil {
		return n, err
	}
	n += tagHeaderSize

	// Write frames
	nn, err := t.writeAllFrames(w)
	if err != nil {
		return n, err
	}
	n += nn

	return
}

// writeAllFrames writes all frames to w and returns
// the number of bytes written and error.
func (t Tag) writeAllFrames(w io.Writer) (int, error) {
	bw := bwpool.Get(w)
	defer bwpool.Put(bw)

	err := t.iterateOverAllFrames(func(id string, f Framer) error {
		return writeFrame(bw, id, f)
	})
	if err != nil {
		return 0, err
	}

	return bw.Buffered(), bw.Flush()
}

func writeFrame(bw *bufio.Writer, id string, frame Framer) error {
	if err := writeFrameHeader(bw, id, frame.Size()); err != nil {
		return err
	}

	if _, err := frame.WriteTo(bw); err != nil {
		return err
	}

	return nil
}

func writeFrameHeader(bw *bufio.Writer, id string, frameSize int) error {
	size, err := util.FormSize(frameSize)
	if err != nil {
		return err
	}

	bw.WriteString(id)
	bw.Write(size)
	bw.Write([]byte{0, 0})
	return nil
}

// Close closes the tag's file, rendering it unusable for I/O.
// It returns an error, if any.
func (t *Tag) Close() error {
	return t.file.Close()
}
