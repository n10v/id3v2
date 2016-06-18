// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// ID3v2 writing library for Go.
//
// This library can only set and write tags, but can't read them.
// So if you only want to set tags, it fits you.
// And if there is a tag of version 3 or 4, this library will just
// delete this tag, because it can't parse tags yet.
// If version of the tag is small than 3, this library will return an error.
//
// Example of creating a new tag and writing it in file:
//	package main
//
//	import (
// 		"github.com/bogem/id3v2"
// 		"log"
// 	)
//
// 	func main() {
//		// Open file and find tag in it
// 		tag, err := id3v2.Open("file.mp3")
// 		if err != nil {
// 		 log.Fatal("Error while opening mp3 file: ", err)
// 		}
//
//		// Set tags
// 		tag.SetArtist("Artist")
// 		tag.SetTitle("Title")
//
//		comment := id3v2.NewCommentFrame()
//		comment.SetLanguage("eng")
//		comment.SetDescription("Short description")
//		comment.SetText("The actual text")
//		tag.AddCommentFrame(comment)
//
//		// Write it to file
// 		if err = tag.Flush(); err != nil {
// 			log.Fatal("Error while flushing a tag: ", err)
// 		}
// 	}
package id3v2

import (
	"bytes"
	"github.com/bogem/id3v2/frame"
	"github.com/bogem/id3v2/util"
	"os"
	"sync"
)

// Possible picture types for picture frame.
const (
	PTOther                   = 0
	PTFileIcon                = 1
	PTOtherFileIcon           = 2
	PTFrontCover              = 3
	PTBackCover               = 4
	PTLeafletPage             = 5
	PTMedia                   = 6
	PTLeadArtistSoloist       = 7
	PTArtistPerformer         = 8
	PTConductor               = 9
	PTBandOrchestra           = 10
	PTComposer                = 11
	PTLyricistTextWriter      = 12
	PTRecordingLocation       = 13
	PTDuringRecording         = 14
	PTDuringPerformance       = 15
	PTMovieScreenCapture      = 16
	PTBrightColouredFish      = 17
	PTIllustration            = 18
	PTBandArtistLogotype      = 19
	PTPublisherStudioLogotype = 20
)

// Possible encodings.
var (
	// ISO-8859-1.
	ENISO = util.Encoding{
		Key:              0,
		TerminationBytes: []byte{0},
	}

	// UTF-16 encoded Unicode with BOM.
	ENUTF16 = util.Encoding{
		Key:              1,
		TerminationBytes: []byte{0, 0},
	}

	// UTF-16BE encoded Unicode without BOM.
	ENUTF16BE = util.Encoding{
		Key:              2,
		TerminationBytes: []byte{0, 0},
	}

	// UTF-8 encoded Unicode.
	ENUTF8 = util.Encoding{
		Key:              3,
		TerminationBytes: []byte{0},
	}
)

var bytesBufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

// Open opens file with string name and find tag in it.
func Open(name string) (*Tag, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return parseTag(file)
}

// NewAttachedPicture creates and initializes a new attached picture frame.
//
// Example of setting a new picture frame to existing tag:
// 		pic := id3v2.NewAttachedPicture()
// 		pic.SetMimeType("image/jpeg")
// 		pic.SetDescription("Cover")
// 		pic.SetPictureType(id3v2.PTFrontCover)
// 		if err := pic.SetPictureFromFile("artwork.jpg"); err != nil {
// 			log.Fatal("Error while setting a picture from file: ", err)
// 		}
// 		tag.AddAttachedPicture(pic)
//
// Available picture types you can see in constants.
func NewAttachedPicture() *frame.PictureFrame {
	pf := new(frame.PictureFrame)
	pf.SetEncoding(ENUTF8)
	return pf
}

// NewCommentFrame creates and initializes a new comment frame.
//
// Example of setting a new comment frame to existing tag:
//	comm := id3v2.NewCommentFrame()
//	comm.SetLanguage("eng")
//	comm.SetDescription("Short description")
//	comm.SetText("The actual text")
//	tag.AddCommentFrame(comm)
//
// You should choose a language code from
// ISO 639-2 code list: https://www.loc.gov/standards/iso639-2/php/code_list.php
func NewCommentFrame() *frame.CommentFrame {
	cf := new(frame.CommentFrame)
	cf.SetEncoding(ENUTF8)
	return cf
}

// NewTextFrame creates and initializes a new text frame. E.g. it is used by
// such frames as artist, title and etc.
//
// Example of setting a new text frame to existing tag:
//	textFrame := id3v2.NewTextFrame("Happy")
//	id := "TMOO" // Mood frame ID
//	tag.AddFrame(id, textFrame)
func NewTextFrame(text string) *frame.TextFrame {
	tf := new(frame.TextFrame)
	tf.SetEncoding(ENUTF8)
	tf.SetText(text)
	return tf
}

// NewUnsynchronisedLyricsFrame creates and initializes a new
// unsynchronised lyrics/text frame.
//
// Example of setting a new unsynchronised lyrics/text frame to existing tag:
//	uslt := id3v2.NewUnsynchronisedLyricsFrame()
//	uslt.SetLanguage("eng")
//	uslt.SetContentDescriptor("Content descriptor")
//	uslt.SetLyrics("Lyrics")
//	tag.AddUnsynchronisedLyricsFrame(uslt)
func NewUnsynchronisedLyricsFrame() *frame.UnsynchronisedLyricsFrame {
	uslf := new(frame.UnsynchronisedLyricsFrame)
	uslf.SetEncoding(ENUTF8)
	return uslf
}
