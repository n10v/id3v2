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
	testSequenceA(t, s)

	s.AddFrame(CommentFrame{Language: "B", Description: "B"})
	testSequenceAB(t, s)

	s.AddFrame(CommentFrame{Language: "B", Description: "B"})
	testSequenceAB(t, s)
}

func TestSequencePictureFramesUniqueness(t *testing.T) {
	t.Parallel()

	s := getSequence()
	defer putSequence(s)

	s.AddFrame(PictureFrame{Description: "A"})
	testSequenceA(t, s)

	s.AddFrame(PictureFrame{Description: "B"})
	testSequenceAB(t, s)

	s.AddFrame(PictureFrame{Description: "B"})
	testSequenceAB(t, s)
}

func TestSequenceUSLFsUniqueness(t *testing.T) {
	t.Parallel()

	s := getSequence()
	defer putSequence(s)

	s.AddFrame(UnsynchronisedLyricsFrame{Language: "A", ContentDescriptor: "A"})
	testSequenceA(t, s)

	s.AddFrame(UnsynchronisedLyricsFrame{Language: "B", ContentDescriptor: "B"})
	testSequenceAB(t, s)

	s.AddFrame(UnsynchronisedLyricsFrame{Language: "B", ContentDescriptor: "B"})
	testSequenceAB(t, s)
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
	got := f.UniqueIdentifier()[:1]
	if got != expected {
		t.Errorf("Expected frame with unique identifier %v, got %v", expected, got)
	}
}
