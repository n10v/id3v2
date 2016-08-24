// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/bogem/id3v2/util"
)

const (
	mp3Name        = "testdata/test.mp3"
	frontCoverName = "testdata/front_cover.jpg"
	backCoverName  = "testdata/back_cover.jpg"
	framesSize     = 222479
	tagSize        = tagHeaderSize + framesSize
	musicSize      = 4557971
)

func TestSetTags(t *testing.T) {
	tag, err := Open(mp3Name)
	if tag == nil || err != nil {
		t.Error("Error while opening mp3 file: ", err)
	}
	tag.SetTitle("Title")
	tag.SetArtist("Artist")
	tag.SetAlbum("Album")
	tag.SetYear("2016")
	tag.SetGenre("Genre")

	// Setting front cover
	frontCover, err := os.Open(frontCoverName)
	if err != nil {
		t.Error("Error while opening front cover file")
	}
	defer frontCover.Close()

	pic := PictureFrame{
		Encoding:    ENUTF8,
		MimeType:    "image/jpeg",
		PictureType: PTFrontCover,
		Description: "Front cover",
		Picture:     frontCover,
	}
	tag.AddFrame(tag.ID("Attached picture"), pic)

	// Setting back cover
	backCover, err := os.Open(backCoverName)
	if err != nil {
		t.Error("Error while opening back cover file")
	}
	defer backCover.Close()

	pic = PictureFrame{
		Encoding:    ENUTF8,
		MimeType:    "image/jpeg",
		PictureType: PTBackCover,
		Description: "Back cover",
		Picture:     backCover,
	}
	tag.AddFrame(tag.ID("Attached picture"), pic)

	// Setting USLTs
	uslt := UnsynchronisedLyricsFrame{
		Encoding:          ENUTF8,
		Language:          "eng",
		ContentDescriptor: "Content descriptor",
		Lyrics:            "bogem/id3v2",
	}
	tag.AddFrame(tag.ID("Unsynchronised lyrics/text transcription"), uslt)

	uslt = UnsynchronisedLyricsFrame{
		Encoding:          ENUTF8,
		Language:          "ger",
		ContentDescriptor: "Inhaltsdeskriptor",
		Lyrics:            "Einigkeit und Recht und Freiheit",
	}
	tag.AddFrame(tag.ID("Unsynchronised lyrics/text transcription"), uslt)

	// Setting comments
	comm := CommentFrame{
		Encoding:    ENUTF8,
		Language:    "eng",
		Description: "Short description",
		Text:        "The actual text",
	}
	tag.AddFrame(tag.ID("Comments"), comm)

	comm = CommentFrame{
		Encoding:    ENUTF8,
		Language:    "ger",
		Description: "Kurze Beschreibung",
		Text:        "Der eigentliche Text",
	}
	tag.AddFrame(tag.ID("Comments"), comm)

	if err = tag.Flush(); err != nil {
		t.Error("Error while closing a tag: ", err)
	}

}

func TestCorrectnessOfSettingTag(t *testing.T) {
	mp3, err := os.Open(mp3Name)
	if err != nil {
		t.Error("Error while opening mp3 file: ", err)
	}
	defer mp3.Close()

	tagHeader := make([]byte, tagHeaderSize)
	n, err := mp3.Read(tagHeader)
	if n != tagHeaderSize {
		t.Errorf("Expected length of header %v, got %v", tagHeaderSize, n)
	}
	if err != nil {
		t.Error("Error while reading a tag header: ", err)
	}

	size := util.ParseSize(tagHeader[6:10])

	if framesSize != size {
		t.Errorf("Expected size of frames: %v, got: %v", framesSize, size)
	}
}

// Check integrity at the beginning of mp3's music part
func TestIntegrityOfMusicAtTheBeginning(t *testing.T) {
	mp3, err := os.Open(mp3Name)
	if err != nil {
		t.Error("Error while opening mp3 file: ", err)
	}
	defer mp3.Close()

	rd := bufio.NewReader(mp3)
	n, err := rd.Discard(tagSize)
	if n != tagSize {
		t.Errorf("Expected length of discarded bytes %v, got %v", tagSize, n)
	}
	if err != nil {
		t.Error("Error while reading mp3 file: ", err)
	}

	expected := []byte{255, 251, 144, 68, 0, 0, 0}
	got := make([]byte, len(expected))
	n, err = rd.Read(got)
	if n != len(expected) {
		t.Errorf("Expected length of read bytes %v, got %v", len(expected), n)
	}
	if err != nil {
		t.Error("Error while reading mp3 file: ", err)
	}

	if !bytes.Equal(expected, got) {
		t.Fail()
	}
}

// Check integrity at the end of mp3's music part
func TestIntegrityOfMusicAtTheEnd(t *testing.T) {
	mp3, err := os.Open(mp3Name)
	if err != nil {
		t.Error("Error while opening mp3 file: ", err)
	}
	defer mp3.Close()

	rd := bufio.NewReader(mp3)
	expected := []byte{0, 0, 0, 0, 0, 0, 255}
	toDiscard := tagSize + musicSize - len(expected)
	n, err := rd.Discard(toDiscard)
	if n != toDiscard {
		t.Errorf("Expected length of discarded bytes %v, got %v", toDiscard, n)
	}
	if err != nil {
		t.Error("Error while discarding: ", err)
	}

	got := make([]byte, len(expected))
	n, err = rd.Read(got)
	if n != len(expected) {
		t.Errorf("Expected length of read bytes %v, got %v", len(expected), n)
	}
	if err != nil {
		t.Error("Error while reading mp3 file: ", err)
	}

	if !bytes.Equal(expected, got) {
		t.Fail()
	}
}
