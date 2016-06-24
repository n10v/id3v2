// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"os"
	"testing"

	"github.com/bogem/id3v2/testutil"
	"github.com/bogem/id3v2/util"
)

const (
	mp3Name        = "testdata/test.mp3"
	frontCoverName = "testdata/front_cover.jpg"
	backCoverName  = "testdata/back_cover.png"
	framesSize     = 62988
	tagSize        = tagHeaderSize + framesSize
	musicSize      = 273310
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
	pic := NewAttachedPicture()
	pic.SetMimeType("image/jpeg")
	pic.SetDescription("Front cover")
	pic.SetPictureType(PTFrontCover)
	if err = pic.SetPictureFromFile(frontCoverName); err != nil {
		t.Error("Error while setting a front cover from file")
	}
	tag.AddAttachedPicture(pic)

	// Setting back cover
	pic = NewAttachedPicture()
	pic.SetMimeType("image/png")
	pic.SetDescription("Back cover")
	pic.SetPictureType(PTBackCover)
	if err = pic.SetPictureFromFile(backCoverName); err != nil {
		t.Error("Error while setting a back cover from file")
	}
	tag.AddAttachedPicture(pic)

	// Setting USLTs
	uslt := NewUnsynchronisedLyricsFrame()
	uslt.SetLanguage("eng")
	uslt.SetContentDescriptor("Content descriptor")
	uslt.SetLyrics("bogem/id3v2")
	tag.AddUnsynchronisedLyricsFrame(uslt)

	uslt = NewUnsynchronisedLyricsFrame()
	uslt.SetLanguage("ger")
	uslt.SetContentDescriptor("Inhaltsdeskriptor")
	uslt.SetLyrics("bogem/id3v2")
	tag.AddUnsynchronisedLyricsFrame(uslt)

	// Setting comments
	cf := NewCommentFrame()
	cf.SetLanguage("eng")
	cf.SetDescription("Short description")
	cf.SetText("The actual text")
	tag.AddCommentFrame(cf)

	cf = NewCommentFrame()
	cf.SetLanguage("ger")
	cf.SetDescription("Kurze Beschreibung")
	cf.SetText("Der eigentliche Text")
	tag.AddCommentFrame(cf)

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

	sizeBytes := tagHeader[6:10]
	size, _ := util.ParseSize(sizeBytes)

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

	expected := []byte{0, 0, 0, 20, 102, 116, 121}
	got := make([]byte, len(expected))
	n, err = rd.Read(got)
	if n != len(expected) {
		t.Errorf("Expected length of read bytes %v, got %v", len(expected), n)
	}
	if err != nil {
		t.Error("Error while reading mp3 file: ", err)
	}

	if err = testutil.AreByteSlicesEqual(expected, got); err != nil {
		t.Error(err)
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
	expected := []byte{3, 162, 192, 0, 3, 224, 203}
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

	if err = testutil.AreByteSlicesEqual(expected, got); err != nil {
		t.Error(err)
	}
}
