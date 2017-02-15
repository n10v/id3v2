// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	tag, err := Open(mp3Name)
	if tag == nil || err != nil {
		t.Error("Error while opening mp3 file:", err)
	}
	defer tag.Close()

	if tag.Count() != countOfFrames {
		t.Errorf("There are %v frames in tag, but should be %v", tag.Count(), countOfFrames)
	}

	if err := testTwoStrings(tag.Artist(), "Artist"); err != nil {
		t.Error(err)
	}
	if err := testTwoStrings(tag.Title(), "Title"); err != nil {
		t.Error(err)
	}
	if err := testTwoStrings(tag.Album(), "Album"); err != nil {
		t.Error(err)
	}
	if err := testTwoStrings(tag.Year(), "2016"); err != nil {
		t.Error(err)
	}
	if err := testTwoStrings(tag.Genre(), "Genre"); err != nil {
		t.Error(err)
	}

	// Check picture frames
	picFrames := tag.GetFrames(tag.CommonID("Attached picture"))
	if len(picFrames) != 2 {
		t.Errorf("Expected picture frames: %v, got %v", 2, len(picFrames))
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

	if err := testPictureFrames(parsedFrontCover, frontCover); err != nil {
		t.Error(err)
	}
	if err := testPictureFrames(parsedBackCover, backCover); err != nil {
		t.Error(err)
	}

	// Check USLT frames
	usltFrames := tag.GetFrames(tag.CommonID("Unsynchronised lyrics/text transcription"))
	if len(picFrames) != 2 {
		t.Errorf("Expected USLT frames: %v, got %v", 2, len(usltFrames))
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

	if err := testUSLTFrames(parsedEngUSLF, engUSLF); err != nil {
		t.Error(err)
	}
	if err := testUSLTFrames(parsedGerUSLF, gerUSLF); err != nil {
		t.Error(err)
	}

	// Check comment frames
	commFrames := tag.GetFrames(tag.CommonID("Comments"))
	if len(commFrames) != 2 {
		t.Errorf("Expected comment frames: %v, got: %v", 2, len(commFrames))
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

	if err := testCommentFrames(parsedEngComm, engComm); err != nil {
		t.Error(err)
	}
	if err := testCommentFrames(parsedGerComm, gerComm); err != nil {
		t.Error(err)
	}

	// Check unknown frame
	parsedUnknownFramer := tag.GetLastFrame(unknownFrameID)
	if parsedUnknownFramer == nil {
		t.Fatal("Parsed unknown frame is nil")
	}
	parsedUnknownFrame := parsedUnknownFramer.(UnknownFrame)
	if err := testUnknownFrames(parsedUnknownFrame, unknownFrame); err != nil {
		t.Error(err)
	}

	if err = tag.Close(); err != nil {
		t.Error("Error while closing a tag:", err)
	}

}

func testPictureFrames(actual, expected PictureFrame) error {
	if err := testTwoBytes(actual.Encoding.Key, expected.Encoding.Key); err != nil {
		return err
	}
	if err := testTwoStrings(actual.MimeType, expected.MimeType); err != nil {
		return err
	}
	if err := testTwoBytes(actual.PictureType, expected.PictureType); err != nil {
		return err
	}
	if err := testTwoStrings(actual.Description, expected.Description); err != nil {
		return err
	}

	if !bytes.Equal(actual.Picture, expected.Picture) {
		return errors.New("Pictures of picture frames' are different")
	}

	return nil
}

func testUSLTFrames(actual, expected UnsynchronisedLyricsFrame) error {
	if err := testTwoBytes(actual.Encoding.Key, expected.Encoding.Key); err != nil {
		return err
	}
	if err := testTwoStrings(actual.Language, expected.Language); err != nil {
		return err
	}
	if err := testTwoStrings(actual.ContentDescriptor, expected.ContentDescriptor); err != nil {
		return err
	}
	if err := testTwoStrings(actual.Lyrics, expected.Lyrics); err != nil {
		return err
	}

	return nil
}

func testCommentFrames(actual, expected CommentFrame) error {
	if err := testTwoBytes(actual.Encoding.Key, expected.Encoding.Key); err != nil {
		return err
	}
	if err := testTwoStrings(actual.Language, expected.Language); err != nil {
		return err
	}
	if err := testTwoStrings(actual.Description, expected.Description); err != nil {
		return err
	}
	if err := testTwoStrings(actual.Text, expected.Text); err != nil {
		return err
	}

	return nil
}

func testUnknownFrames(actual, expected UnknownFrame) error {
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

func testTwoStrings(actual, expected string) error {
	if actual != expected {
		return fmt.Errorf("Expected %v, got %v", expected, actual)
	}
	return nil
}

func testTwoBytes(actual, expected byte) error {
	if actual != expected {
		return fmt.Errorf("Expected %v, got %v", expected, actual)
	}
	return nil
}
