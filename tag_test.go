// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/bogem/id3v2/util"
)

const (
	mp3Name        = "testdata/test.mp3"
	frontCoverName = "testdata/front_cover.jpg"
	backCoverName  = "testdata/back_cover.jpg"
	framesSize     = 222524
	tagSize        = tagHeaderSize + framesSize
	musicSize      = 4557971
)

var (
	frontCover = PictureFrame{
		Encoding:    ENUTF8,
		MimeType:    "image/jpeg",
		PictureType: PTFrontCover,
		Description: "Front cover",
	}
	backCover = PictureFrame{
		Encoding:    ENUTF8,
		MimeType:    "image/jpeg",
		PictureType: PTBackCover,
		Description: "Back cover",
	}

	engUSLF = UnsynchronisedLyricsFrame{
		Encoding:          ENUTF8,
		Language:          "eng",
		ContentDescriptor: "Content descriptor",
		Lyrics:            "bogem/id3v2",
	}
	gerUSLF = UnsynchronisedLyricsFrame{
		Encoding:          ENUTF8,
		Language:          "ger",
		ContentDescriptor: "Inhaltsdeskriptor",
		Lyrics:            "Einigkeit und Recht und Freiheit",
	}

	engComm = CommentFrame{
		Encoding:    ENUTF8,
		Language:    "eng",
		Description: "Short description",
		Text:        "The actual text",
	}
	gerComm = CommentFrame{
		Encoding:    ENUTF8,
		Language:    "ger",
		Description: "Kurze Beschreibung",
		Text:        "Der eigentliche Text",
	}

	unknownFrameID = "WPUB"
	unknownFrame   = UnknownFrame{
		body: []byte("https://soundcloud.com/suicidepart2"),
	}
)

func TestBlankID(t *testing.T) {
	// Delete all frames in tag and add one blank id
	tag, err := Open(mp3Name)
	if tag == nil || err != nil {
		t.Fatal("Error while opening mp3 file: ", err)
	}

	tag.DeleteAllFrames()
	tag.AddFrame("", frontCover)

	// tag.Save should write no frames to file
	if err = tag.Save(); err != nil {
		t.Error("Error while saving a tag: ", err)
	}

	tag.Close()

	// Parse tag. It should be no frames
	parsedTag, err := Open(mp3Name)
	if parsedTag == nil || err != nil {
		t.Fatal("Error while opening mp3 file: ", err)
	}

	allFrames := parsedTag.AllFrames()
	if len(allFrames) > 0 {
		t.Error("There should be no frames in tag, but there are", len(allFrames))
	}
}

func TestSetTags(t *testing.T) {
	tag, err := Open(mp3Name)
	if tag == nil || err != nil {
		t.Fatal("Error while opening mp3 file: ", err)
	}
	defer tag.Close()

	tag.SetTitle("Title")
	tag.SetArtist("Artist")
	tag.SetAlbum("Album")
	tag.SetYear("2016")
	tag.SetGenre("Genre")

	// Set picture frames
	if err := resetPictureReaders(); err != nil {
		t.Fatal(err)
	}
	tag.AddAttachedPicture(frontCover)
	tag.AddAttachedPicture(backCover)

	// Set USLTs
	tag.AddUnsynchronisedLyricsFrame(engUSLF)
	tag.AddUnsynchronisedLyricsFrame(gerUSLF)

	// Set comments
	tag.AddCommentFrame(engComm)
	tag.AddCommentFrame(gerComm)

	// Set unknown frame
	tag.AddFrame(unknownFrameID, unknownFrame)

	if err = tag.Save(); err != nil {
		t.Error("Error while saving a tag: ", err)
	}
}

func TestCorrectnessOfSettingTag(t *testing.T) {
	mp3, err := os.Open(mp3Name)
	if err != nil {
		t.Fatal("Error while opening mp3 file: ", err)
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

	size, err := util.ParseSize(tagHeader[6:10])
	if err != nil {
		t.Error("Error while parsing a tag header size: ", err)
	}

	if framesSize != size {
		t.Errorf("Expected size of frames: %v, got: %v", framesSize, size)
	}
}

func resetPictureReaders() error {
	frontCoverFile, err := os.Open(frontCoverName)
	if err != nil {
		return errors.New("Error while opening front cover file: " + err.Error())
	}
	frontCover.Picture = frontCoverFile

	backCoverFile, err := os.Open(backCoverName)
	if err != nil {
		return errors.New("Error while opening back cover file: " + err.Error())
	}
	backCover.Picture = backCoverFile

	return nil
}

// Check integrity at the beginning of mp3's music part
func TestIntegrityOfMusicAtTheBeginning(t *testing.T) {
	mp3, err := os.Open(mp3Name)
	if err != nil {
		t.Fatal("Error while opening mp3 file: ", err)
	}
	defer mp3.Close()

	rd := bufio.NewReader(mp3)
	n, err := rd.Discard(tagSize)
	if n != tagSize {
		t.Errorf("Expected length of discarded bytes %v, got %v", tagSize, n)
	}
	if err != nil {
		t.Fatal("Error while reading mp3 file: ", err)
	}

	expected := []byte{255, 251, 144, 68, 0, 0, 0}
	got := make([]byte, len(expected))
	n, err = rd.Read(got)
	if n != len(expected) {
		t.Errorf("Expected length of read bytes %v, got %v", len(expected), n)
	}
	if err != nil {
		t.Fatal("Error while reading mp3 file: ", err)
	}

	if !bytes.Equal(expected, got) {
		t.Fail()
	}
}

// Check integrity at the end of mp3's music part
func TestIntegrityOfMusicAtTheEnd(t *testing.T) {
	mp3, err := os.Open(mp3Name)
	if err != nil {
		t.Fatal("Error while opening mp3 file: ", err)
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
		t.Fatal("Error while discarding: ", err)
	}

	got := make([]byte, len(expected))
	n, err = rd.Read(got)
	if n != len(expected) {
		t.Errorf("Expected length of read bytes %v, got %v", len(expected), n)
	}
	if err != nil {
		t.Fatal("Error while reading mp3 file: ", err)
	}

	if !bytes.Equal(expected, got) {
		t.Fail()
	}
}
