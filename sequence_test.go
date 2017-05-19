// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import "testing"

func TestSequenceCacheUpdate(t *testing.T) {
	t.Parallel()

	s := getSequence()
	defer putSequence(s)
	s.AddFrame(PictureFrame{Description: "A"})
	if len(s.Frames()) != 1 {
		t.Errorf("Expected %v frames, got %v", 1, len(s.Frames()))
	}
	s.AddFrame(PictureFrame{Description: "B"})
	if len(s.Frames()) != 2 {
		t.Errorf("Expected %v frames, got %v", 2, len(s.Frames()))
	}
}

func TestSequenceCommentFramesUniqueness(t *testing.T) {
	t.Parallel()

	s := getSequence()
	defer putSequence(s)
	s.AddFrame(CommentFrame{Language: "A", Description: "A"})
	if s.Count() != 1 {
		t.Errorf("Expected %v frames, got %v", 1, s.Count())
	}
	s.AddFrame(CommentFrame{Language: "B", Description: "B"})
	if s.Count() != 2 {
		t.Errorf("Expected %v frames, got %v", 2, s.Count())
	}
	s.AddFrame(CommentFrame{Language: "B", Description: "B"})
	if s.Count() != 2 {
		t.Errorf("Expected %v frames, got %v", 2, s.Count())
	}
}

func TestSequencePictureFramesUniqueness(t *testing.T) {
	t.Parallel()

	s := getSequence()
	defer putSequence(s)
	s.AddFrame(PictureFrame{Description: "A"})
	if s.Count() != 1 {
		t.Errorf("Expected %v frames, got %v", 1, s.Count())
	}
	s.AddFrame(PictureFrame{Description: "B"})
	if s.Count() != 2 {
		t.Errorf("Expected %v frames, got %v", 2, s.Count())
	}
	s.AddFrame(PictureFrame{Description: "B"})
	if s.Count() != 2 {
		t.Errorf("Expected %v frames, got %v", 2, s.Count())
	}
}

func TestSequenceUSLFsUniqueness(t *testing.T) {
	t.Parallel()

	s := getSequence()
	defer putSequence(s)
	s.AddFrame(UnsynchronisedLyricsFrame{Language: "A", ContentDescriptor: "A"})
	if s.Count() != 1 {
		t.Errorf("Expected %v frames, got %v", 1, s.Count())
	}
	s.AddFrame(UnsynchronisedLyricsFrame{Language: "B", ContentDescriptor: "B"})
	if s.Count() != 2 {
		t.Errorf("Expected %v frames, got %v", 2, s.Count())
	}
	s.AddFrame(UnsynchronisedLyricsFrame{Language: "B", ContentDescriptor: "B"})
	if s.Count() != 2 {
		t.Errorf("Expected %v frames, got %v", 2, s.Count())
	}
}
