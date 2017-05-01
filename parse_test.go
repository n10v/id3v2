// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bogem/id3v2/util"
)

// TestParseInvalidFrameSize creates new temp file, writes tag header,
// valid TIT2 frame and frame with invalid size to it, then checks
// if valid frame is parsed and there is only this frame in tag.
func TestParseInvalidFrameSize(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal("Error while opening mp3 file:", err)
	}
	defer os.Remove(file.Name())

	size, _ := util.FormSize(16 + 10)

	// Write tag header
	bw := bufio.NewWriter(file)
	if err := writeTagHeader(bw, size, 4); err != nil {
		t.Fatal(err)
	}
	if err := bw.Flush(); err != nil {
		t.Fatal(err)
	}
	// Write valid TIT2 frame
	file.Write([]byte{0x54, 0x49, 0x54, 0x32, 00, 00, 00, 06, 00, 00, 03, 0x54, 0x69, 0x74, 0x6C, 0x65})
	// Write invalid frame (size byte can't be more than 127)
	file.Write([]byte{0x54, 0x49, 0x54, 0x32, 255, 255, 255, 255, 00, 00})

	file.Seek(0, os.SEEK_SET)

	tag, err := OpenFile(file, defaultOpts)
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

// TestParse compares parsed frames with expected frames.
func TestParse(t *testing.T) {
	var err error

	if err = resetMP3Tag(); err != nil {
		t.Fatal("Error while reseting mp3 file:", err)
	}

	tag, err := Open(mp3Name, defaultOpts)
	if tag == nil || err != nil {
		t.Error("Error while opening mp3 file:", err)
	}

	if err = compareTwoStrings(tag.Artist(), "Artist"); err != nil {
		t.Error(err)
	}
	if err = compareTwoStrings(tag.Title(), "Title"); err != nil {
		t.Error(err)
	}
	if err = compareTwoStrings(tag.Album(), "Album"); err != nil {
		t.Error(err)
	}
	if err = compareTwoStrings(tag.Year(), "2016"); err != nil {
		t.Error(err)
	}
	if err = compareTwoStrings(tag.Genre(), "Genre"); err != nil {
		t.Error(err)
	}
	if err = testPictureFrames(tag); err != nil {
		t.Error(err)
	}
	if err = testUSLTFrames(tag); err != nil {
		t.Error(err)
	}
	if err = testCommentFrames(tag); err != nil {
		t.Error(err)
	}
	if err = testUnknownFrames(tag); err != nil {
		t.Error(err)
	}
}

func testPictureFrames(tag *Tag) error {
	picFrames := tag.GetFrames(tag.CommonID("Attached picture"))
	if len(picFrames) != 2 {
		return fmt.Errorf("Expected picture frames: %v, got %v", 2, len(picFrames))
	}

	var parsedFrontCover, parsedBackCover PictureFrame
	for _, f := range picFrames {
		pf, ok := f.(PictureFrame)
		if !ok {
			return errors.New("Couldn't assert picture frame")
		}
		if pf.PictureType == PTFrontCover {
			parsedFrontCover = pf
		}
		if pf.PictureType == PTBackCover {
			parsedBackCover = pf
		}
	}

	if err := comparePictureFrames(parsedFrontCover, frontCover); err != nil {
		return err
	}
	if err := comparePictureFrames(parsedBackCover, backCover); err != nil {
		return err
	}

	return nil
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

func testUSLTFrames(tag *Tag) error {
	usltFrames := tag.GetFrames(tag.CommonID("Unsynchronised lyrics/text transcription"))
	if len(usltFrames) != 2 {
		return fmt.Errorf("Expected USLT frames: %v, got %v", 2, len(usltFrames))
	}

	var parsedEngUSLF, parsedGerUSLF UnsynchronisedLyricsFrame
	for _, f := range usltFrames {
		uslf, ok := f.(UnsynchronisedLyricsFrame)
		if !ok {
			return errors.New("Couldn't assert USLT frame")
		}
		if uslf.Language == "eng" {
			parsedEngUSLF = uslf
		}
		if uslf.Language == "ger" {
			parsedGerUSLF = uslf
		}
	}

	if err := compareUSLTFrames(parsedEngUSLF, engUSLF); err != nil {
		return err
	}
	if err := compareUSLTFrames(parsedGerUSLF, gerUSLF); err != nil {
		return err
	}

	return nil
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

func testCommentFrames(tag *Tag) error {
	commFrames := tag.GetFrames(tag.CommonID("Comments"))
	if len(commFrames) != 2 {
		return fmt.Errorf("Expected comment frames: %v, got: %v", 2, len(commFrames))
	}

	var parsedEngComm, parsedGerComm CommentFrame
	for _, f := range commFrames {
		cf, ok := f.(CommentFrame)
		if !ok {
			return errors.New("Couldn't assert comment frame")
		}
		if cf.Language == "eng" {
			parsedEngComm = cf
		}
		if cf.Language == "ger" {
			parsedGerComm = cf
		}
	}

	if err := compareCommentFrames(parsedEngComm, engComm); err != nil {
		return err
	}
	if err := compareCommentFrames(parsedGerComm, gerComm); err != nil {
		return err
	}

	return nil
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

func testUnknownFrames(tag *Tag) error {
	parsedUnknownFramer := tag.GetLastFrame(unknownFrameID)
	if parsedUnknownFramer == nil {
		return errors.New("Parsed unknown frame is nil")
	}
	parsedUnknownFrame := parsedUnknownFramer.(UnknownFrame)
	if err := compareUnknownFrames(parsedUnknownFrame, unknownFrame); err != nil {
		return err
	}

	if err := tag.Close(); err != nil {
		return errors.New("Error while closing a tag: " + err.Error())
	}

	return nil
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

func compareTwoStrings(actual, expected string) error {
	if actual != expected {
		return fmt.Errorf("Expected %q, got %q", expected, actual)
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
// if parseTag will not parse the tag, if Options{Parse: false} is set.
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
// if tag.parseAllFrames will parse only frames, that set in Options.ParseFrames.
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
