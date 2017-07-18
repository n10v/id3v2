// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package id3v2 is the ID3 parsing and writing library for Go.
package id3v2

import (
	"io"
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
	EncodingISO = util.Encoding{
		Key:              0,
		TerminationBytes: []byte{0},
	}

	// UTF-16 encoded Unicode with BOM.
	EncodingUTF16 = util.Encoding{
		Key:              1,
		TerminationBytes: []byte{0, 0},
	}

	// UTF-16BE encoded Unicode without BOM.
	EncodingUTF16BE = util.Encoding{
		Key:              2,
		TerminationBytes: []byte{0, 0},
	}

	// UTF-8 encoded Unicode.
	EncodingUTF8 = util.Encoding{
		Key:              3,
		TerminationBytes: []byte{0},
	}

	Encodings = []util.Encoding{EncodingISO, EncodingUTF16, EncodingUTF16BE, EncodingUTF8}
)

// Open opens file with name and passes it to OpenFile.
// If there is no tag in file, it will create new one with version ID3v2.4.
func Open(name string, opts Options) (*Tag, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return ParseReader(file, opts)
}

// ParseReader parses rd and finds tag in it considering opts.
// If there is no tag in rd, it will create new one with version ID3v2.4.
func ParseReader(rd io.Reader, opts Options) (*Tag, error) {
	tag := NewEmptyTag()
	err := tag.parse(rd, opts)
	return tag, err
}

// NewEmptyTag returns an empty ID3v2.4 tag without any frames and reader.
func NewEmptyTag() *Tag {
	tag := new(Tag)
	tag.init(nil, 0, 4)
	return tag
}
