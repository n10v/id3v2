// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// TestParse compares parsed frames with expected frames.
func TestParse(t *testing.T) {
	if err := resetMP3Tag(); err != nil {
		t.Fatal("Error while reseting mp3 file:", err)
	}

	tag, err := Open(mp3Name, defaultOpts)
	if tag == nil || err != nil {
		t.Error("Error while opening mp3 file:", err)
	}
	defer tag.Close()

	testTextFrames(t, tag)
	testPictureFrames(t, tag)
	testUSLTFrames(t, tag)
	testCommentFrames(t, tag)
	testUnknownFrames(t, tag)
}

func testTextFrames(t *testing.T, tag *Tag) {
	if err := compareTwoStrings(tag.Artist(), "Artist"); err != nil {
		t.Error(err)
	}
	if err := compareTwoStrings(tag.Title(), "Title"); err != nil {
		t.Error(err)
	}
	if err := compareTwoStrings(tag.Album(), "Album"); err != nil {
		t.Error(err)
	}
	if err := compareTwoStrings(tag.Year(), "2016"); err != nil {
		t.Error(err)
	}
	if err := compareTwoStrings(tag.Genre(), "Genre"); err != nil {
		t.Error(err)
	}
}

func compareTwoStrings(actual, expected string) error {
	if actual != expected {
		return fmt.Errorf("Expected %q, got %q", expected, actual)
	}
	return nil
}

func testPictureFrames(t *testing.T, tag *Tag) {
	picFrames := tag.GetFrames(tag.CommonID("Attached picture"))
	if len(picFrames) != 2 {
		t.Fatalf("Expected picture frames: %v, got %v", 2, len(picFrames))
	}

	var parsedFrontCover, parsedBackCover PictureFrame
	for _, f := range picFrames {
		pf, ok := f.(PictureFrame)
		if !ok {
			t.Fatal("Couldn't assert picture frame")
		}
		if pf.PictureType == PTFrontCover {
			parsedFrontCover = pf
		}
		if pf.PictureType == PTBackCover {
			parsedBackCover = pf
		}
	}

	if err := comparePictureFrames(parsedFrontCover, frontCover); err != nil {
		t.Error(err)
	}
	if err := comparePictureFrames(parsedBackCover, backCover); err != nil {
		t.Error(err)
	}
}

func comparePictureFrames(actual, expected PictureFrame) error {
	if err := compareTwoBytes(actual.Encoding.Key, expected.Encoding.Key); err != nil {
		return err
	}
	if err := compareTwoStrings(actual.MimeType, expected.MimeType); err != nil {
		return err
	}
	if err := compareTwoBytes(actual.PictureType, expected.PictureType); err != nil {
		return err
	}
	if err := compareTwoStrings(actual.Description, expected.Description); err != nil {
		return err
	}

	if !bytes.Equal(actual.Picture, expected.Picture) {
		return errors.New("Pictures of picture frames' are different")
	}

	return nil
}

func testUSLTFrames(t *testing.T, tag *Tag) {
	usltFrames := tag.GetFrames(tag.CommonID("Unsynchronised lyrics/text transcription"))
	if len(usltFrames) != 2 {
		t.Fatalf("Expected USLT frames: %v, got %v", 2, len(usltFrames))
	}

	var parsedEngUSLF, parsedGerUSLF UnsynchronisedLyricsFrame
	for _, f := range usltFrames {
		uslf, ok := f.(UnsynchronisedLyricsFrame)
		if !ok {
			t.Fatal("Couldn't assert USLT frame")
		}
		if uslf.Language == "eng" {
			parsedEngUSLF = uslf
		}
		if uslf.Language == "ger" {
			parsedGerUSLF = uslf
		}
	}

	if err := compareUSLTFrames(parsedEngUSLF, engUSLF); err != nil {
		t.Error(err)
	}
	if err := compareUSLTFrames(parsedGerUSLF, gerUSLF); err != nil {
		t.Error(err)
	}
}

func compareUSLTFrames(actual, expected UnsynchronisedLyricsFrame) error {
	if err := compareTwoBytes(actual.Encoding.Key, expected.Encoding.Key); err != nil {
		return err
	}
	if err := compareTwoStrings(actual.Language, expected.Language); err != nil {
		return err
	}
	if err := compareTwoStrings(actual.ContentDescriptor, expected.ContentDescriptor); err != nil {
		return err
	}
	if err := compareTwoStrings(actual.Lyrics, expected.Lyrics); err != nil {
		return err
	}

	return nil
}

func testCommentFrames(t *testing.T, tag *Tag) {
	commFrames := tag.GetFrames(tag.CommonID("Comments"))
	if len(commFrames) != 2 {
		t.Fatalf("Expected comment frames: %v, got: %v", 2, len(commFrames))
	}

	var parsedEngComm, parsedGerComm CommentFrame
	for _, f := range commFrames {
		cf, ok := f.(CommentFrame)
		if !ok {
			t.Fatal("Couldn't assert comment frame")
		}
		if cf.Language == "eng" {
			parsedEngComm = cf
		}
		if cf.Language == "ger" {
			parsedGerComm = cf
		}
	}

	if err := compareCommentFrames(parsedEngComm, engComm); err != nil {
		t.Error(err)
	}
	if err := compareCommentFrames(parsedGerComm, gerComm); err != nil {
		t.Error(err)
	}
}

func compareCommentFrames(actual, expected CommentFrame) error {
	if err := compareTwoBytes(actual.Encoding.Key, expected.Encoding.Key); err != nil {
		return err
	}
	if err := compareTwoStrings(actual.Language, expected.Language); err != nil {
		return err
	}
	if err := compareTwoStrings(actual.Description, expected.Description); err != nil {
		return err
	}
	if err := compareTwoStrings(actual.Text, expected.Text); err != nil {
		return err
	}

	return nil
}

func testUnknownFrames(t *testing.T, tag *Tag) {
	parsedUnknownFramer := tag.GetLastFrame(unknownFrameID)
	if parsedUnknownFramer == nil {
		t.Fatal("Parsed unknown frame is nil")
	}
	parsedUnknownFrame := parsedUnknownFramer.(UnknownFrame)
	if err := compareUnknownFrames(parsedUnknownFrame, unknownFrame); err != nil {
		t.Error(err)
	}
}

func compareUnknownFrames(actual, expected UnknownFrame) error {
	actualBody := new(bytes.Buffer)
	expectedBody := new(bytes.Buffer)
	if _, err := actual.WriteTo(actualBody); err != nil {
		return err
	}
	if _, err := expected.WriteTo(expectedBody); err != nil {
		return err
	}
	if !bytes.Equal(actualBody.Bytes(), expectedBody.Bytes()) {
		return errors.New("Body of unknown frame isn't the same as expected")
	}
	return nil
}

func compareTwoBytes(actual, expected byte) error {
	if actual != expected {
		return fmt.Errorf("Expected %v, got %v", expected, actual)
	}
	return nil
}

// TestParseOptionsParseFalse checks,
// if parseTag() will not parse the tag, if Options{Parse: false} is set.
func TestParseOptionsParseFalse(t *testing.T) {
	tag, err := Open(mp3Name, Options{Parse: false})
	if tag == nil || err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}
	if tag.HasFrames() {
		t.Errorf("tag has %v frames, but should have no frames at all", tag.Count())
	}
}

// TestParseOptionsParseFrames checks,
// if tag.parseAllFrames() will parse only frames, that set in Options.ParseFrames.
func TestParseOptionsParseFrames(t *testing.T) {
	tag, err := Open(mp3Name, Options{Parse: true, ParseFrames: []string{"Artist", "Title"}})
	if tag == nil || err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}
	if tag.Count() > 2 {
		t.Errorf("tag should have only artist and title frames, but it has %v frames", tag.Count())
	}
	if tag.Artist() == "" {
		t.Errorf("tag should have an artist, but it doesn't")
	}
	if tag.Title() == "" {
		t.Errorf("tag should have a title, but it doesn't")
	}
}

// TestParseInvalidFrameSize creates an empty tag, writes tag header,
// valid TIT2 frame and frame with invalid size, then checks
// if valid frame is parsed and there is only this frame in tag.
func TestParseInvalidFrameSize(t *testing.T) {
	t.Parallel()

	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)

	// Write tag header.
	dst := make([]byte, 4)
	if err := writeTagHeader(bw, dst, tagHeaderSize+16, 4); err != nil {
		t.Fatal(err)
	}
	if err := bw.Flush(); err != nil {
		t.Fatal(err)
	}
	// Write valid TIT2 frame.
	buf.Write([]byte{0x54, 0x49, 0x54, 0x32, 00, 00, 00, 06, 00, 00, 03}) // header and encoding
	buf.WriteString("Title")
	// Write invalid frame (size byte can't be greater than 127).
	buf.Write([]byte{0x54, 0x49, 0x54, 0x32, 255, 255, 255, 255, 00, 00})

	tag, err := ParseReader(buf, defaultOpts)
	if tag == nil || err != nil {
		t.Fatal("Error while parsing mp3 file:", err)
	}
	if tag.Title() != "Title" {
		t.Errorf("Expected title: %q, got: %q", "Title", tag.Title())
	}
	if tag.Count() != 1 {
		t.Error("There should be only 1 frame in tag, but there are", tag.Count())
	}
}

// TestParseEmptyReader checks if ParseReader() correctly parses empty readers.
func TestParseEmptyReader(t *testing.T) {
	t.Parallel()

	tag, err := ParseReader(new(bytes.Buffer), Options{Parse: true})
	if err != nil {
		t.Error("Error while parsing empty reader:", err)
	}
	if tag.HasFrames() {
		t.Error("Tag should not have any frames, but it has", tag.Count())
	}
}

// TestParseReaderNil checks
// if ParseReader returns correct error when calling ParseReader(nil, Options{}).
func TestParseReaderNil(t *testing.T) {
	t.Parallel()

	_, err := ParseReader(nil, Options{Parse: true})
	if err == nil {
		t.Fatal("Expected that err is not nil, but err is nil")
	}
	if !strings.Contains(err.Error(), "rd is nil") {
		t.Fatalf("Expected err contains %q, got %q", "rd is nil", err)
	}

}
