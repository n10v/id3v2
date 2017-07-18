// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/bogem/id3v2/bspool"
	"github.com/bogem/id3v2/bwpool"
	"github.com/bogem/id3v2/util"
)

var ErrNoFile = errors.New("tag was not initialized with file")

// Tag stores all information about opened tag.
type Tag struct {
	frames    map[string]Framer
	sequences map[string]*sequence

	reader io.Reader

	originalSize int64
	version      byte
}

// AddFrame adds f to tag with appropriate id. If id is "" or f is nil,
// AddFrame will not add it to tag.
//
// If you want to add attached picture, comment or unsynchronised lyrics/text
// transcription frames, better use AddAttachedPicture, AddCommentFrame
// or AddUnsynchronisedLyricsFrame methods respectively.
func (tag *Tag) AddFrame(id string, f Framer) {
	if id == "" || f == nil {
		return
	}

	if mustFrameBeInSequence(id) {
		if tag.sequences[id] == nil {
			tag.sequences[id] = getSequence()
		}
		tag.sequences[id].AddFrame(f)
	} else {
		tag.frames[id] = f
	}
}

// AddAttachedPicture adds a picture frame to tag.
func (tag *Tag) AddAttachedPicture(pf PictureFrame) {
	id := tag.CommonID("Attached picture")
	tag.AddFrame(id, pf)
}

// AddCommentFrame adds a comment frame to tag.
func (tag *Tag) AddCommentFrame(cf CommentFrame) {
	id := tag.CommonID("Comments")
	tag.AddFrame(id, cf)
}

// AddUnsynchronisedLyricsFrame adds an unsynchronised lyrics/text frame
// to tag.
func (tag *Tag) AddUnsynchronisedLyricsFrame(uslf UnsynchronisedLyricsFrame) {
	id := tag.CommonID("Unsynchronised lyrics/text transcription")
	tag.AddFrame(id, uslf)
}

// CommonID returns ID3v2.3 or ID3v2.4 (in appropriate to version of Tag) frame ID
// from given description.
// For example, CommonID("Language") will return "TLAN".
// If it can't find the ID with given description, it returns the description.
//
// All descriptions you can find in file common_ids.go
// or in id3 documentation (for fourth version: http://id3.org/id3v2.4.0-frames;
// for third version: http://id3.org/id3v2.3.0#Declared_ID3v2_frames).
func (tag *Tag) CommonID(description string) string {
	var ids map[string]string
	if tag.version == 3 {
		ids = V23CommonIDs
	} else {
		ids = V24CommonIDs
	}
	if id, ok := ids[description]; ok {
		return id
	}
	return description
}

// AllFrames returns map, that contains all frames in tag, that could be parsed.
// The key of this map is an ID of frame and value is an array of frames.
func (tag *Tag) AllFrames() map[string][]Framer {
	frames := make(map[string][]Framer)

	for id, f := range tag.frames {
		frames[id] = []Framer{f}
	}
	for id, sequence := range tag.sequences {
		frames[id] = sequence.Frames()
	}

	return frames
}

// DeleteAllFrames deletes all frames in tag.
func (tag *Tag) DeleteAllFrames() {
	if tag.frames == nil || len(tag.frames) > 0 {
		tag.frames = make(map[string]Framer)
	}
	if tag.sequences == nil || len(tag.sequences) > 0 {
		for _, s := range tag.sequences {
			putSequence(s)
		}
		tag.sequences = make(map[string]*sequence)
	}
}

// DeleteFrames deletes frames in tag with given id.
func (tag *Tag) DeleteFrames(id string) {
	delete(tag.frames, id)
	if s, ok := tag.sequences[id]; ok {
		putSequence(s)
		delete(tag.sequences, id)
	}
}

// Reset deletes all frames in tag and parses rd considering opts.
func (tag *Tag) Reset(rd io.Reader, opts Options) error {
	tag.DeleteAllFrames()
	return tag.parse(rd, opts)
}

// GetFrames returns frames with corresponding id.
// It returns nil if there is no frames with given id.
func (tag *Tag) GetFrames(id string) []Framer {
	if f, exists := tag.frames[id]; exists {
		return []Framer{f}
	} else if s, exists := tag.sequences[id]; exists {
		return s.Frames()
	}
	return nil
}

// GetLastFrame returns last frame from slice, that is returned from GetFrames function.
// GetLastFrame is suitable for frames, that can be only one in whole tag.
// For example, for text frames.
func (tag *Tag) GetLastFrame(id string) Framer {
	// Avoid an allocation of slice in GetFrames,
	// if there is anyway one frame.
	if f, exists := tag.frames[id]; exists {
		return f
	}

	fs := tag.GetFrames(id)
	if len(fs) == 0 {
		return nil
	}
	return fs[len(fs)-1]
}

// GetTextFrame returns text frame with corresponding id.
func (tag *Tag) GetTextFrame(id string) TextFrame {
	f := tag.GetLastFrame(id)
	if f == nil {
		return TextFrame{}
	}
	tf := f.(TextFrame)
	return tf
}

// Count returns the number of frames in tag.
func (tag *Tag) Count() int {
	n := len(tag.frames)
	for _, s := range tag.sequences {
		n += s.Count()
	}
	return n
}

// HasFrames checks if there is at least one frame in tag.
// It's much faster than tag.Count() > 0.
func (tag *Tag) HasFrames() bool {
	return len(tag.frames) > 0 || len(tag.sequences) > 0
}

func (tag *Tag) Title() string {
	f := tag.GetTextFrame(tag.CommonID("Title/Songname/Content description"))
	return f.Text
}

func (tag *Tag) SetTitle(title string) {
	tag.AddFrame(tag.CommonID("Title/Songname/Content description"), TextFrame{Encoding: EncodingUTF8, Text: title})
}

func (tag *Tag) Artist() string {
	f := tag.GetTextFrame(tag.CommonID("Lead artist/Lead performer/Soloist/Performing group"))
	return f.Text
}

func (tag *Tag) SetArtist(artist string) {
	tag.AddFrame(tag.CommonID("Lead artist/Lead performer/Soloist/Performing group"), TextFrame{Encoding: EncodingUTF8, Text: artist})
}

func (tag *Tag) Album() string {
	f := tag.GetTextFrame(tag.CommonID("Album/Movie/Show title"))
	return f.Text
}

func (tag *Tag) SetAlbum(album string) {
	tag.AddFrame(tag.CommonID("Album/Movie/Show title"), TextFrame{Encoding: EncodingUTF8, Text: album})
}

func (tag *Tag) Year() string {
	f := tag.GetTextFrame(tag.CommonID("Year"))
	return f.Text
}

func (tag *Tag) SetYear(year string) {
	tag.AddFrame(tag.CommonID("Year"), TextFrame{Encoding: EncodingUTF8, Text: year})
}

func (tag *Tag) Genre() string {
	f := tag.GetTextFrame(tag.CommonID("Content type"))
	return f.Text
}

func (tag *Tag) SetGenre(genre string) {
	tag.AddFrame(tag.CommonID("Content type"), TextFrame{Encoding: EncodingUTF8, Text: genre})
}

// iterateOverAllFrames iterates over every single frame in tag and calls
// f for them. It consumps no memory at all, unlike the tag.AllFrames().
// It returns error only if f returns error.
func (tag *Tag) iterateOverAllFrames(f func(id string, frame Framer) error) error {
	for id, frame := range tag.frames {
		if err := f(id, frame); err != nil {
			return err
		}
	}
	for id, sequence := range tag.sequences {
		for _, frame := range sequence.Frames() {
			if err := f(id, frame); err != nil {
				return err
			}
		}
	}
	return nil
}

// Size returns the size of all frames in bytes.
func (tag *Tag) Size() int {
	if !tag.HasFrames() {
		return 0
	}

	var n int
	n += tagHeaderSize // Add the size of tag header
	tag.iterateOverAllFrames(func(id string, f Framer) error {
		n += frameHeaderSize + f.Size() // Add the whole frame size
		return nil
	})

	return n
}

// Version returns current ID3v2 version of tag.
func (tag *Tag) Version() byte {
	return tag.version
}

// SetVersion sets given ID3v2 version to tag.
// If version is less than 3 or greater than 4, then this method will do nothing.
// If tag has some frames, which are deprecated or changed in given version,
// then to your notice you can delete, change or just stay them.
func (tag *Tag) SetVersion(version byte) {
	if version < 3 || version > 4 {
		return
	}
	tag.version = version
}

// Save writes tag to the file, if tag was opened with a file.
// If there are no frames in tag, Save will write
// only music part without any ID3v2 information.
// If tag was initiliazed not with file, it returns ErrNoFile.
func (tag *Tag) Save() error {
	file, ok := tag.reader.(*os.File)
	if !ok {
		return ErrNoFile
	}

	// Get original file mode.
	originalFile := file
	originalStat, err := originalFile.Stat()
	if err != nil {
		return err
	}

	// Create a temp file for mp3 file, which will contain new tag.
	name := file.Name() + "-id3v2"
	newFile, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, originalStat.Mode())
	if err != nil {
		return err
	}

	// Make sure we clean up the temp file if it's still around.
	defer os.Remove(newFile.Name())

	// Write tag in new file.
	tagSize, err := tag.WriteTo(newFile)
	if err != nil {
		return err
	}

	// Seek to a music part of original file.
	if _, err = originalFile.Seek(tag.originalSize, os.SEEK_SET); err != nil {
		return err
	}

	// Write to new file the music part.
	buf := bspool.Get(32 * 1024)
	defer bspool.Put(buf)
	if _, err = io.CopyBuffer(newFile, originalFile, buf); err != nil {
		return err
	}

	// Close files to allow replacing.
	newFile.Close()
	originalFile.Close()

	// Replace original file with new file.
	if err = os.Rename(newFile.Name(), originalFile.Name()); err != nil {
		return err
	}

	// Set tag.reader to new file with original name.
	tag.reader, err = os.Open(originalFile.Name())
	if err != nil {
		return err
	}

	// Set tag.originalSize to new frames size.
	if tagSize > tagHeaderSize {
		tag.originalSize = tagSize - tagHeaderSize
	} else {
		tag.originalSize = 0
	}

	return nil
}

// WriteTo writes whole tag in w if there is at least one frame.
// It returns the number of bytes written and error during the write.
// It returns nil as error if the write was successful.
func (tag *Tag) WriteTo(w io.Writer) (n int64, err error) {
	if w == nil {
		return 0, errors.New("w is nil")
	}

	// Count size of frames.
	framesSize := tag.Size() - tagHeaderSize
	if framesSize <= 0 {
		return 0, nil
	}

	// Write tag header.
	bw := bwpool.Get(w)
	defer bwpool.Put(bw)
	if err := writeTagHeader(bw, framesSize, tag.version); err != nil {
		return n, err
	}
	n += tagHeaderSize

	// Write frames.
	err = tag.iterateOverAllFrames(func(id string, f Framer) error {
		nn, err := writeFrame(bw, id, f)
		n += nn
		return err
	})
	if err != nil {
		return n, err
	}

	return n, bw.Flush()
}

func writeFrame(bw *bufio.Writer, id string, frame Framer) (int64, error) {
	if err := writeFrameHeader(bw, id, frame.Size()); err != nil {
		return 0, err
	}

	frameSize, err := frame.WriteTo(bw)
	return frameHeaderSize + frameSize, err
}

func writeFrameHeader(bw *bufio.Writer, id string, frameSize int) error {
	// ID
	if _, err := bw.WriteString(id); err != nil {
		return err
	}

	// Size
	if err := util.WriteBytesSize(bw, frameSize); err != nil {
		return err
	}

	// Flags
	if _, err := bw.Write([]byte{0, 0}); err != nil {
		return err
	}
	return nil
}

// Close closes tag's file, if tag was opened with a file.
// If tag was initiliazed not with file, it returns ErrNoFile.
func (tag *Tag) Close() error {
	file, ok := tag.reader.(*os.File)
	if !ok {
		return ErrNoFile
	}
	return file.Close()
}
