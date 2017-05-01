// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package id3v2 is the ID3 parsing and writing library for Go.
//
// Example of usage:
//
//	package main
//
//	import (
//		"fmt"
//		"log"
//
//		"github.com/bogem/id3v2"
//	)
//
//	func main() {
//		// Open file and parse tag in it.
//		tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
//		if err != nil {
//			log.Fatal("Error while opening mp3 file: ", err)
//		}
//		defer tag.Close()
//
//		// Read frames.
//		fmt.Println(tag.Artist())
//		fmt.Println(tag.Title())
//
//		// Set simple text frames.
//		tag.SetArtist("Artist")
//		tag.SetTitle("Title")
//
//		// Set comment frame.
//		comment := id3v2.CommentFrame{
//			Encoding:    id3v2.ENUTF8,
//			Language:    "eng",
//			Description: "My opinion",
//			Text:        "Very good song",
//		}
//		tag.AddCommentFrame(comment)
//
//		// Write tag to file.
//		if err = tag.Save(); err != nil {
//			log.Fatal("Error while saving a tag: ", err)
//		}
//	}
package id3v2

import (
	"os"

	"github.com/bogem/id3v2/util"
)

// Available picture types for picture frame.
const (
	PTOther = iota
	PTFileIcon
	PTOtherFileIcon
	PTFrontCover
	PTBackCover
	PTLeafletPage
	PTMedia
	PTLeadArtistSoloist
	PTArtistPerformer
	PTConductor
	PTBandOrchestra
	PTComposer
	PTLyricistTextWriter
	PTRecordingLocation
	PTDuringRecording
	PTDuringPerformance
	PTMovieScreenCapture
	PTBrightColouredFish
	PTIllustration
	PTBandArtistLogotype
	PTPublisherStudioLogotype
)

// Available encodings.
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

	Encodings = []util.Encoding{ENISO, ENUTF16, ENUTF16BE, ENUTF8}
)

// Open opens file with name and passes it to OpenFile.
func Open(name string, opts Options) (*Tag, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return OpenFile(file, opts)
}

// OpenFile parses opened file and finds tag in it considering opts.
// If there is no tag in file, OpenFile will create new one with version ID3v2.4.
func OpenFile(file *os.File, opts Options) (*Tag, error) {
	return parseTag(file, opts)
}
