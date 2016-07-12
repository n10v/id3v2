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
	"os"

	"github.com/bogem/id3v2/util"
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

// Open opens file with string name and find tag in it.
func Open(name string) (*Tag, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return parseTag(file)
}
