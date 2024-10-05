// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"sync"
	"testing"
)

const (
	mp3Path        = "testdata/test.mp3"
	multiMp3Path   = "testdata/test_multi.mp3"
	frontCoverPath = "testdata/front_cover.jpg"
	backCoverPath  = "testdata/back_cover.jpg"

	framesSize    = 211978
	tagSize       = tagHeaderSize + framesSize
	musicSize     = 3840834
	countOfFrames = 15
)

var (
	frontCover = PictureFrame{
		Encoding:    EncodingUTF8,
		MimeType:    "image/jpeg",
		PictureType: PTFrontCover,
		Description: "Front cover",
		Picture:     mustReadFile(frontCoverPath),
	}
	backCover = PictureFrame{
		Encoding:    EncodingUTF8,
		MimeType:    "image/jpeg",
		PictureType: PTBackCover,
		Description: "Back cover",
		Picture:     mustReadFile(backCoverPath),
	}

	engUSLF = UnsynchronisedLyricsFrame{
		Encoding:          EncodingUTF8,
		Language:          "eng",
		ContentDescriptor: "Content descriptor",
		Lyrics:            "bogem/id3v2",
	}
	gerUSLF = UnsynchronisedLyricsFrame{
		Encoding:          EncodingUTF8,
		Language:          "ger",
		ContentDescriptor: "Inhaltsdeskriptor",
		Lyrics:            "Einigkeit und Recht und Freiheit",
	}

	musicBrainzUDTF = UserDefinedTextFrame{
		Encoding:    EncodingUTF8,
		Description: "MusicBrainz Album Id",
		Value:       "fbd94fb6-2a74-42d0-acbc-81caf8b84984",
	}

	multiUDTF = UserDefinedTextFrame{
		Encoding:    EncodingUTF8,
		Description: "multi",
		Value:       "val1",
		Multi: []string{
			"val1",
			"val2",
		},
	}

	musicBrainzUF = UFIDFrame{
		OwnerIdentifier: "https://musicbrainz.org",
		Identifier:      []byte("fbd94fb6-2a74-42d0-acbc-81caf8b84984"),
	}

	engComm = CommentFrame{
		Encoding:    EncodingUTF8,
		Language:    "eng",
		Description: "Short description",
		Text:        "The actual text",
	}
	gerComm = CommentFrame{
		Encoding:    EncodingUTF8,
		Language:    "ger",
		Description: "Kurze Beschreibung",
		Text:        "Der eigentliche Text",
	}

	popmFrame = PopularimeterFrame{
		Email:   "foo@bar.com",
		Rating:  128,
		Counter: big.NewInt(10000000000000000),
	}

	unknownFrameID = "WPUB"
	unknownFrame   = UnknownFrame{
		Body: []byte("https://soundcloud.com/suicidepart2"),
	}

	// Parse all frames
	parseOpts = Options{Parse: true}
)

func init() {
	if err := resetMP3Tag(); err != nil {
		panic(fmt.Sprintf("Error while reseting mp3 file: %v", err))
	}
}

// resetMP3Tag sets default tag in file located by mp3Path.
func resetMP3Tag() error {
	tag, err := Open(mp3Path, Options{Parse: false})
	if tag == nil || err != nil {
		return err
	}
	defer tag.Close()

	tag.SetTitle("Title")
	tag.SetArtist("Artist")
	tag.SetAlbum("Album")
	tag.SetYear("2016")
	tag.SetGenre("Genre")

	tag.AddAttachedPicture(frontCover)
	tag.AddAttachedPicture(backCover)

	tag.AddUnsynchronisedLyricsFrame(engUSLF)
	tag.AddUnsynchronisedLyricsFrame(gerUSLF)

	tag.AddUserDefinedTextFrame(musicBrainzUDTF)
	tag.AddUFIDFrame(musicBrainzUF)

	tag.AddFrame(tag.CommonID("Popularimeter"), popmFrame)

	tag.AddCommentFrame(engComm)
	tag.AddCommentFrame(gerComm)

	tag.AddFrame(unknownFrameID, unknownFrame)

	return tag.Save()
}

func mustReadFile(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("can't read %q: %v", path, err))
	}
	return contents
}

func TestCountLenSize(t *testing.T) {
	tag, err := Open(mp3Path, parseOpts)
	if tag == nil || err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}
	defer tag.Close()

	// Check count.
	if tag.Count() != countOfFrames {
		t.Errorf("Expected frames: %v, got: %v", countOfFrames, tag.Count())
	}

	// Check len of tag.AllFrames().
	if len(tag.AllFrames()) != 12 {
		t.Errorf("Expected: %v, got: %v", 11, len(tag.AllFrames()))
	}

	// Check saved tag size by reading the 6:10 bytes of mp3 file.
	mp3, err := os.Open(mp3Path)
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

	size, err := parseSize(tagHeader[6:10], true)
	if err != nil {
		t.Error("Error while parsing a tag header size:", err)
	}

	if framesSize != size {
		t.Errorf("Expected size of frames: %v, got: %v", framesSize, size)
	}

	// Check tag.Size().
	if tag.Size() != tagSize {
		t.Errorf("Expected tag.Size(): %v, got: %v", tagSize, tag.Size())
	}
}

// TestIntegrityOfMusicAtTheBeginning checks
// if tag.Save() doesn't truncate or add some extra bytes at the beginning
// of music part.
func TestIntegrityOfMusicAtTheBeginning(t *testing.T) {
	mp3, err := os.Open(mp3Path)
	if err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}
	defer mp3.Close()

	rd := bufio.NewReader(mp3)
	n, err := rd.Discard(tagSize)
	if n != tagSize {
		t.Errorf("Expected length of discarded bytes %v, got %v", tagSize, n)
	}
	if err != nil {
		t.Fatal("Error while reading mp3 file:", err)
	}

	expected := []byte{255, 251, 80, 0, 0, 0, 0}
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
}

// TestIntegrityOfMusicAtTheEnd checks
// if tag.Save() doesn't truncate music part or add some extra bytes at the end
// of music part.
func TestIntegrityOfMusicAtTheEnd(t *testing.T) {
	mp3, err := os.Open(mp3Path)
	if err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}
	defer mp3.Close()

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
}

// TestCheckPermissions checks
// if tag.Save() creates file with the same permissions of original file.
func TestCheckPermissions(t *testing.T) {
	originalFile, err := os.Open(mp3Path)
	if err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}
	originalStat, err := originalFile.Stat()
	if err != nil {
		t.Fatal("Error while getting mp3 file stat:", err)
	}
	originalMode := originalStat.Mode()
	originalFile.Close()

	tag, err := Open(mp3Path, parseOpts)
	if err != nil {
		t.Fatal("Error while parsing a tag:", err)
	}
	if err = tag.Save(); err != nil {
		t.Error("Error while saving a tag:", err)
	}
	tag.Close()

	newFile, err := os.Open(mp3Path)
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

// TestBlankID creates empty tag, adds frame with blank id and checks
// if no tag is written by tag.WriteTo() (it must not write tag to buf
// if there are no or only blank frames).
func TestBlankID(t *testing.T) {
	t.Parallel()

	tag := NewEmptyTag()
	tag.AddFrame("", unknownFrame)

	if tag.Count() > 0 {
		t.Error("There should be no frames in tag, but there are", tag.Count())
	}

	if tag.HasFrames() {
		t.Error("tag.HasFrames() should return false, but it returns true")
	}

	if tag.Size() != 0 {
		t.Error("Size of tag should be 0. Actual tag size:", tag.Size())
	}

	// It should write no frames to buf.
	buf := new(bytes.Buffer)
	if _, err := tag.WriteTo(buf); err != nil {
		t.Fatal("Error while writing a tag:", err)
	}
	if buf.Len() > 0 {
		t.Fatal("tag.WriteTo(buf) should write no frames, but it wrote")
	}
}

// TestInvalidLanguageCommentFrame checks
// if tag.WriteTo() returns the correct error by writing the comment frame with
// incorrect length of language code.
func TestInvalidLanguageCommentFrame(t *testing.T) {
	t.Parallel()

	tag := NewEmptyTag()
	tag.AddCommentFrame(CommentFrame{
		Encoding: EncodingUTF8,
		Language: "en", // should be "eng" according to ISO 639-2
		Text:     "The actual text",
	})

	_, err := tag.WriteTo(ioutil.Discard)
	if err == nil {
		t.Fatal("tag.WriteTo() must return the error about invalid language code")
	}
	if !strings.Contains(err.Error(), "must consist") {
		t.Fatalf("Incorrect error. Expected error contains %q, got %q", "must consist", err)
	}
}

// TestInvalidLanguageUSLF checks
// if tag.WriteTo() returns the correct error by writing the comment frame with
// incorrect length of language code.
func TestInvalidLanguageUSLF(t *testing.T) {
	t.Parallel()

	tag := NewEmptyTag()
	tag.AddUnsynchronisedLyricsFrame(UnsynchronisedLyricsFrame{
		Encoding: EncodingUTF8,
		Language: "en", // should be "eng" according to ISO 639-2
		Lyrics:   "Lyrics",
	})

	_, err := tag.WriteTo(ioutil.Discard)
	if err == nil {
		t.Fatal("tag.WriteTo() must return the error about invalid language code")
	}
	if !strings.Contains(err.Error(), "must consist") {
		t.Fatalf("Incorrect error. Expected error contains %q, got %q", "must consist", err)
	}
}

// TestSaveAndCloseEmptyTag checks
// if tag.Save() and tag.Close() return an error for empty tag.
func TestSaveAndCloseEmptyTag(t *testing.T) {
	t.Parallel()

	tag := NewEmptyTag()
	if err := tag.Save(); err == nil {
		t.Error("By saving empty tag we wait for an error, but it's not returned")
	}
	if err := tag.Close(); err == nil {
		t.Error("By closing empty tag we wait for an error, but it's not returned")
	}
}

// TestEmptyTagWriteTo checks
// if tag.WriteTo() works correctly for empty tag.
func TestEmptyTagWriteTo(t *testing.T) {
	t.Parallel()

	tag := NewEmptyTag()
	tag.SetArtist("Artist")
	tag.SetTitle("Title")

	buf := new(bytes.Buffer)
	if _, err := tag.WriteTo(buf); err != nil {
		t.Fatal("Error while writing to buf:", err)
	}
	if buf.Len() == 0 {
		t.Fatal("buf is empty, but it must have tag")
	}

	parsedTag, err := ParseReader(buf, parseOpts)
	if err != nil {
		t.Fatal("Error while parsing buf:", err)
	}
	if parsedTag.Artist() != "Artist" {
		t.Error("Expected Artist, got", parsedTag.Artist())
	}
	if parsedTag.Title() != "Title" {
		t.Error("Expected Title, got", parsedTag.Title())
	}
}

// TestResetTag checks if tag.Reset() correctly parses mp3.
func TestResetTag(t *testing.T) {
	mp3, err := os.Open(mp3Path)
	if err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}
	defer mp3.Close()

	tag := NewEmptyTag()
	if err := tag.Reset(mp3, parseOpts); err != nil {
		t.Fatal("Error while reseting tag:", err)
	}

	// Check if it parsed.
	if tag.Count() != countOfFrames {
		t.Errorf("Expected frames: %v, got: %v", countOfFrames, tag.Count())
	}
}

// TestConcurrent creates sync.Pool with tags, executes 100 goroutines and
// checks if id3v2 works correctly in concurrent environment (one tag per goroutine).
func TestConcurrent(t *testing.T) {
	tagPool := sync.Pool{New: func() interface{} { return NewEmptyTag() }}

	ec := make(chan error, 100)

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()

			tag := tagPool.Get().(*Tag)
			defer tagPool.Put(tag)

			file, err := os.Open(mp3Path)
			if err != nil {
				ec <- fmt.Errorf("Error while opening mp3: %v", err)
				return
			}
			defer file.Close()

			if err := tag.Reset(file, parseOpts); err != nil {
				ec <- fmt.Errorf("Error while reseting tag to file: %v", err)
				return
			}

			// Just check if it's correctly parsed.
			if tag.Count() != countOfFrames {
				ec <- fmt.Errorf("Expected frames: %v, got: %v", countOfFrames, tag.Count())
				return
			}

			if _, err := tag.WriteTo(ioutil.Discard); err != nil {
				ec <- fmt.Errorf("Error while writing to ioutil.Discard: %v", err)
				return
			}
		}()
	}
	wg.Wait()
	close(ec)

	err := <-ec
	if err != nil {
		t.Fatal(err)
	}
}

// TestEncodedText checks
// if texts of frames encoded with different encodings are correctly written.
func TestEncodedText(t *testing.T) {
	t.Parallel()

	encoded := "Héllö"
	tag := NewEmptyTag()
	tag.AddFrame(tag.CommonID("Title"), TextFrame{
		Encoding: EncodingISO,
		Text:     encoded,
	})
	tag.AddFrame(tag.CommonID("Attached picture"), PictureFrame{
		Encoding:    EncodingUTF16,
		MimeType:    "image/jpeg",
		PictureType: PTFrontCover,
		Description: encoded,
	})
	tag.AddFrame(tag.CommonID("Unsynchronised lyrics/text transcription"), UnsynchronisedLyricsFrame{
		Encoding:          EncodingUTF16BE,
		Language:          "ger",
		ContentDescriptor: encoded,
		Lyrics:            encoded,
	})
	tag.AddFrame(tag.CommonID("Comments"), CommentFrame{
		Encoding:    EncodingUTF8,
		Language:    "eng",
		Description: encoded,
		Text:        encoded,
	})

	buf := new(bytes.Buffer)
	n, err := tag.WriteTo(buf)
	if err != nil {
		t.Fatalf("Error by writing to buf: %v", err)
	}
	if n != int64(tag.Size()) {
		t.Errorf("Expected WriteTo n==%v, got %v", tag.Size(), n)
	}

	tag, err = ParseReader(buf, parseOpts)
	if err != nil {
		t.Fatalf("Error by parsing the tag: %v", err)
	}

	tf := tag.GetLastFrame(tag.CommonID("Title")).(TextFrame)
	if !tf.Encoding.Equals(EncodingISO) && tf.Text != encoded {
		t.Errorf("Expected %q and %q, got %q and %q", EncodingISO, encoded, tf.Encoding, tf.Text)
	}

	pf := tag.GetLastFrame(tag.CommonID("Attached picture")).(PictureFrame)
	if !pf.Encoding.Equals(EncodingUTF16) && pf.Description != encoded {
		t.Errorf("Expected %q and %q, got %q and %q", EncodingISO, encoded, pf.Encoding, pf.Description)
	}

	uslf := tag.GetLastFrame(tag.CommonID("Unsynchronised lyrics/text transcription")).(UnsynchronisedLyricsFrame)
	if !uslf.Encoding.Equals(EncodingUTF16BE) && uslf.ContentDescriptor != encoded && uslf.Lyrics != encoded {
		t.Errorf("Expected %q, %q and %q; got %q, %q and %q", EncodingISO, encoded, encoded, uslf.Encoding, uslf.ContentDescriptor, uslf.Lyrics)
	}

	cf := tag.GetLastFrame(tag.CommonID("Comments")).(CommentFrame)
	if !cf.Encoding.Equals(EncodingUTF8) && cf.Description != encoded && cf.Text != encoded {
		t.Errorf("Expected %q, %q and %q; got %q, %q and %q", EncodingUTF8, encoded, encoded, cf.Encoding, cf.Description, cf.Text)
	}
}

func TestWriteToN(t *testing.T) {
	tag := NewEmptyTag()

	tag.SetTitle("Title")
	tag.AddAttachedPicture(frontCover)
	tag.AddUnsynchronisedLyricsFrame(engUSLF)
	tag.AddCommentFrame(engComm)
	tag.AddFrame(unknownFrameID, unknownFrame)

	buf := new(bytes.Buffer)
	n, err := tag.WriteTo(buf)
	if err != nil {
		t.Fatalf("Error by writing: %v", err)
	}
	if n != int64(tag.Size()) {
		t.Errorf("Expected WriteTo n==%v, got %v", tag.Size(), n)
	}
	if int64(buf.Len()) != n {
		t.Errorf("buf.Len() and n are not equal: %v != %v ", buf.Len(), n)
	}
}
