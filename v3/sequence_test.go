// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import "testing"

func TestSequenceCommentFramesUniqueness(t *testing.T) {
	t.Parallel()

	s := getSequence()
	defer putSequence(s)

	s.AddFrame(CommentFrame{Language: "A", Description: "A"})
	testSequenceCount(t, s, 1)
	testFrameUniqueIdentifier(t, s.Frames()[0], "AA")

	s.AddFrame(CommentFrame{Language: "B", Description: "B"})
	testSequenceCount(t, s, 2)
	testFrameUniqueIdentifier(t, s.Frames()[0], "AA")
	testFrameUniqueIdentifier(t, s.Frames()[1], "BB")

	s.AddFrame(CommentFrame{Language: "B", Description: "B"})
	testSequenceCount(t, s, 2)
	testFrameUniqueIdentifier(t, s.Frames()[0], "AA")
	testFrameUniqueIdentifier(t, s.Frames()[1], "BB")
}

func TestSequencePictureFramesUniqueness(t *testing.T) {
	t.Parallel()

	s := getSequence()
	defer putSequence(s)

	s.AddFrame(PictureFrame{Description: "A", PictureType: 0x00})
	testSequenceCount(t, s, 1)
	testFrameUniqueIdentifier(t, s.Frames()[0], "00A")

	// Test against https://github.com/bogem/id3v2/issues/65 regression.
	s.AddFrame(PictureFrame{Description: "A", PictureType: 0x01})
	testSequenceCount(t, s, 2)
	testFrameUniqueIdentifier(t, s.Frames()[0], "00A")
	testFrameUniqueIdentifier(t, s.Frames()[1], "01A")

	s.AddFrame(PictureFrame{Description: "B", PictureType: 0x00})
	testSequenceCount(t, s, 3)
	testFrameUniqueIdentifier(t, s.Frames()[0], "00A")
	testFrameUniqueIdentifier(t, s.Frames()[1], "01A")
	testFrameUniqueIdentifier(t, s.Frames()[2], "00B")

	s.AddFrame(PictureFrame{Description: "B", PictureType: 0x00})
	testSequenceCount(t, s, 3)
	testFrameUniqueIdentifier(t, s.Frames()[0], "00A")
	testFrameUniqueIdentifier(t, s.Frames()[1], "01A")
	testFrameUniqueIdentifier(t, s.Frames()[2], "00B")
}

func TestSequenceUSLFsUniqueness(t *testing.T) {
	t.Parallel()

	s := getSequence()
	defer putSequence(s)

	s.AddFrame(UnsynchronisedLyricsFrame{Language: "A", ContentDescriptor: "A"})
	testSequenceCount(t, s, 1)
	testFrameUniqueIdentifier(t, s.Frames()[0], "AA")

	s.AddFrame(UnsynchronisedLyricsFrame{Language: "B", ContentDescriptor: "B"})
	testSequenceCount(t, s, 2)
	testFrameUniqueIdentifier(t, s.Frames()[0], "AA")
	testFrameUniqueIdentifier(t, s.Frames()[1], "BB")

	s.AddFrame(UnsynchronisedLyricsFrame{Language: "B", ContentDescriptor: "B"})
	testSequenceCount(t, s, 2)
	testFrameUniqueIdentifier(t, s.Frames()[0], "AA")
	testFrameUniqueIdentifier(t, s.Frames()[1], "BB")
}

func TestSequenceUDTFsUniqueness(t *testing.T) {
	t.Parallel()

	s := getSequence()
	defer putSequence(s)

	s.AddFrame(UserDefinedTextFrame{Description: "A"})
	testSequenceA(t, s)

	s.AddFrame(UserDefinedTextFrame{Description: "B", Value: "B"})
	testSequenceAB(t, s)

	s.AddFrame(UserDefinedTextFrame{Description: "B", Value: "C"})
	testSequenceAB(t, s)

	// If one more frame added with same unique identifier, it should rewrite the old one.
	// See https://github.com/bogem/id3v2/issues/42
	valueOfLastFrame := s.Frames()[1].(UserDefinedTextFrame).Value
	if valueOfLastFrame != "C" {
		t.Fatalf("Expected value of UserDefinedTextFrame %q, got %q", "C", valueOfLastFrame)
	}
}

func testSequenceA(t *testing.T, s *sequence) {
	testSequenceCount(t, s, 1)
	testFrameUniqueIdentifier(t, s.Frames()[0], "A")
}

func testSequenceAB(t *testing.T, s *sequence) {
	testSequenceCount(t, s, 2)
	testFrameUniqueIdentifier(t, s.Frames()[0], "A")
	testFrameUniqueIdentifier(t, s.Frames()[1], "B")
}

func testSequenceCount(t *testing.T, s *sequence, expected int) {
	got := s.Count()
	if got != expected {
		t.Errorf("Expected %v frames, got %v", expected, got)
	}
}

func testFrameUniqueIdentifier(t *testing.T, f Framer, expected string) {
	got := f.UniqueIdentifier()
	if got != expected {
		t.Errorf("Expected frame with unique identifier %v, got %v", expected, got)
	}
}
