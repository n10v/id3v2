// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/bogem/id3v2/util"
)

const (
	mp3Name        = "testdata/test.mp3"
	frontCoverName = "testdata/front_cover.jpg"
	backCoverName  = "testdata/back_cover.jpg"
	framesSize     = 222524
	tagSize        = tagHeaderSize + framesSize
	musicSize      = 4557843
	countOfFrames  = 12
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

	// Parse all frames
	defaultOpts = Options{
		Parse: true,
	}
)

func init() {
	var err error

	// Set covers' picture.
	frontCover.Picture, err = ioutil.ReadFile(frontCoverName)
	if err != nil {
		panic("Error while reading front cover file: " + err.Error())
	}
	backCover.Picture, err = ioutil.ReadFile(backCoverName)
	if err != nil {
		panic("Error while reading back cover file: " + err.Error())
	}
	if err := resetMP3Tag(); err != nil {
		panic("Error while reseting mp3 file: " + err.Error())
	}
}

// resetMP3Tag sets the default frames to mp3Name.
func resetMP3Tag() error {
	tag, err := Open(mp3Name, defaultOpts)
	if tag == nil || err != nil {
		return err
	}

	tag.SetTitle("Title")
	tag.SetArtist("Artist")
	tag.SetAlbum("Album")
	tag.SetYear("2016")
	tag.SetGenre("Genre")

	// Set picture frames
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
		return err
	}

	if err = tag.Close(); err != nil {
		return err
	}

	return nil
}

func TestCountLenSize(t *testing.T) {
	tag, err := Open(mp3Name, defaultOpts)
	if tag == nil || err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}

	// Check count
	if tag.Count() != countOfFrames {
		t.Errorf("Expected frames: %v, got: %v", countOfFrames, tag.Count())
	}

	// Check len of tag.AllFrames()
	if len(tag.AllFrames()) != 9 {
		t.Errorf("Expected: %v, got: %v", 9, len(tag.AllFrames()))
	}

	// Check saved tag size by reading the 6:10 bytes of mp3 file
	mp3, err := os.Open(mp3Name)
	if err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}
	tagHeader := make([]byte, tagHeaderSize)
	n, err := mp3.Read(tagHeader)
	if n != tagHeaderSize {
		t.Errorf("Expected length of header %v, got %v", tagHeaderSize, n)
	}
	if err != nil {
		t.Error("Error while reading a tag header:", err)
	}

	size, err := util.ParseSize(tagHeader[6:10])
	if err != nil {
		t.Error("Error while parsing a tag header size:", err)
	}

	if framesSize != size {
		t.Errorf("Expected size of frames: %v, got: %v", framesSize, size)
	}

	// Check tag.Size
	tagSize := tagHeaderSize + framesSize
	if tag.Size() != tagSize {
		t.Errorf("Expected tag.Size(): %v, got: %v", tagSize, tag.Size())
	}

	if err = tag.Close(); err != nil {
		t.Error("Error while closing a tag:", err)
	}
}

// TestIntegrityOfMusicAtTheBeginning checks
// if tag.Save doesn't truncate or add some extra bytes at the beginning
// of music part.
func TestIntegrityOfMusicAtTheBeginning(t *testing.T) {
	mp3, err := os.Open(mp3Name)
	if err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}

	rd := bufio.NewReader(mp3)
	n, err := rd.Discard(tagSize)
	if n != tagSize {
		t.Errorf("Expected length of discarded bytes %v, got %v", tagSize, n)
	}
	if err != nil {
		t.Fatal("Error while reading mp3 file:", err)
	}

	expected := []byte{255, 251, 144, 68, 0, 0, 0}
	got := make([]byte, len(expected))
	n, err = rd.Read(got)
	if n != len(expected) {
		t.Errorf("Expected length of read bytes %v, got %v", len(expected), n)
	}
	if err != nil {
		t.Fatal("Error while reading mp3 file:", err)
	}

	if !bytes.Equal(expected, got) {
		t.Fatalf("Expected %v, got %v", expected, got)
	}

	if err = mp3.Close(); err != nil {
		t.Error("Error while closing a tag:", err)
	}
}

// TestIntegrityOfMusicAtTheEnd checks
// if tag.Save doesn't truncate music part or add some extra bytes at the end
// of music part.
func TestIntegrityOfMusicAtTheEnd(t *testing.T) {
	mp3, err := os.Open(mp3Name)
	if err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}

	rd := bufio.NewReader(mp3)
	expected := []byte{85, 85, 85, 85, 85, 85, 85}
	toDiscard := tagSize + musicSize - len(expected)
	n, err := rd.Discard(toDiscard)
	if n != toDiscard {
		t.Errorf("Expected length of discarded bytes %v, got %v", toDiscard, n)
	}
	if err != nil {
		t.Fatal("Error while discarding:", err)
	}

	got := make([]byte, len(expected))
	n, err = rd.Read(got)
	if n != len(expected) {
		t.Errorf("Expected length of read bytes %v, got %v", len(expected), n)
	}
	if err != nil {
		t.Fatal("Error while reading mp3 file:", err)
	}

	if !bytes.Equal(expected, got) {
		t.Fatalf("Expected %v, got %v", expected, got)
	}

	if err = mp3.Close(); err != nil {
		t.Error("Error while closing a tag:", err)
	}
}

// TestCheckPermissions checks
// if tag.Save creates file with the same permissions of original file.
func TestCheckPermissions(t *testing.T) {
	originalFile, err := os.Open(mp3Name)
	if err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}
	originalStat, err := originalFile.Stat()
	if err != nil {
		t.Fatal("Error while getting mp3 file stat:", err)
	}
	originalMode := originalStat.Mode()
	originalFile.Close()

	tag, err := Open(mp3Name, defaultOpts)
	if err != nil {
		t.Fatal("Error while parsing a tag:", err)
	}
	if err = tag.Save(); err != nil {
		t.Error("Error while saving a tag:", err)
	}
	if err = tag.Close(); err != nil {
		t.Error("Error while closing a tag:", err)
	}

	newFile, err := os.Open(mp3Name)
	if err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}
	newStat, err := newFile.Stat()
	if err != nil {
		t.Fatal("Error while getting mp3 file stats:", err)
	}
	newMode := newStat.Mode()

	if originalMode != newMode {
		t.Errorf("Expected permissions: %v, got %v", originalMode, newMode)
	}
}

// TestBlankID deletes all frames in tag, adds frame with blank id and checks
// if no tag is written by tag.Size (tag.WriteTo must not write tag to file
// if there are 0 frames).
func TestBlankID(t *testing.T) {
	tag, err := Open(mp3Name, defaultOpts)
	if tag == nil || err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}

	tag.DeleteAllFrames()
	tag.AddFrame("", unknownFrame)

	if tag.Count() > 0 {
		t.Error("There should be no frames in tag, but there are", tag.Count())
	}

	if tag.HasFrames() {
		t.Error("tag.HasFrames should return false, but it returns true")
	}

	if tag.Size() != 0 {
		t.Error("Size of tag should be 0. Actual tag size:", tag.Size())
	}

	// tag.Save should write no frames to file
	if err = tag.Save(); err != nil {
		t.Error("Error while saving a tag:", err)
	}

	if err = tag.Close(); err != nil {
		t.Error("Error while closing a tag:", err)
	}

	// Parse tag. It should be no frames
	parsedTag, err := Open(mp3Name, defaultOpts)
	if parsedTag == nil || err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}

	if tag.Count() > 0 {
		t.Error("There should be no frames in parsed tag, but there are", tag.Count())
	}

	if tag.HasFrames() {
		t.Error("Parsed tag.HasFrames should return false, but it returns true")
	}

	if tag.Size() != 0 {
		t.Error("Size of parsed tag should be 0. Actual tag size:", tag.Size())
	}
}

// TestInvalidLanguageCommentFrame checks,
// if tag.Save returns the correct error by writing the comment frame with
// incorrect length of language code.
func TestInvalidLanguageCommentFrame(t *testing.T) {
	tag, err := Open(mp3Name, defaultOpts)
	if tag == nil || err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}

	tag.DeleteAllFrames()
	tag.AddCommentFrame(CommentFrame{
		Encoding: ENUTF8,
		Language: "en", // should be "eng" according to ISO 639-2
		Text:     "The actual text",
	})

	err = tag.Save()
	if err == nil {
		t.Fatal("tag.Save must return the error about invalid language code")
	}
	if !strings.Contains(err.Error(), "must consist") {
		t.Fatalf("Incorrect error. Expected error contains %q, got %q", "must consist", err)
	}

}

// TestInvalidLanguageUSLF checks,
// if tag.Save returns the correct error by writing the comment frame with
// incorrect length of language code.
func TestInvalidLanguageUSLF(t *testing.T) {
	tag, err := Open(mp3Name, defaultOpts)
	if tag == nil || err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}

	tag.DeleteAllFrames()
	tag.AddUnsynchronisedLyricsFrame(UnsynchronisedLyricsFrame{
		Encoding: ENUTF8,
		Language: "en", // should be "eng" according to ISO 639-2
		Lyrics:   "Lyrics",
	})

	err = tag.Save()
	if err == nil {
		t.Fatal("tag.Save must return the error about invalid language code")
	}
	if !strings.Contains(err.Error(), "must consist") {
		t.Fatalf("Incorrect error. Expected error contains %q, got %q", "must consist", err)
	}

}
