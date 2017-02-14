package id3v2

import "testing"

func TestCommentSequenceFramesUniqueness(t *testing.T) {
	s := newCommentSequence()
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

func TestCommentSequenceCacheUpdate(t *testing.T) {
	s := newCommentSequence()
	s.AddFrame(CommentFrame{Language: "A", Description: "A"})
	if len(s.Frames()) != 1 {
		t.Errorf("Expected %v frames, got %v", 1, len(s.Frames()))
	}
	s.AddFrame(CommentFrame{Language: "B", Description: "B"})
	if len(s.Frames()) != 2 {
		t.Errorf("Expected %v frames, got %v", 2, len(s.Frames()))
	}
}

func TestPictureSequenceFramesUniqueness(t *testing.T) {
	s := newPictureSequence()
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

func TestPictureSequenceCacheUpdate(t *testing.T) {
	s := newPictureSequence()
	s.AddFrame(PictureFrame{Description: "A"})
	if len(s.Frames()) != 1 {
		t.Errorf("Expected %v frames, got %v", 1, len(s.Frames()))
	}
	s.AddFrame(PictureFrame{Description: "B"})
	if len(s.Frames()) != 2 {
		t.Errorf("Expected %v frames, got %v", 2, len(s.Frames()))
	}
}

func TestUSLFSequenceFramesUniqueness(t *testing.T) {
	s := newUSLFSequence()
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

func TestUSLFSequenceCacheUpdate(t *testing.T) {
	s := newUSLFSequence()
	s.AddFrame(UnsynchronisedLyricsFrame{Language: "A", ContentDescriptor: "A"})
	if len(s.Frames()) != 1 {
		t.Errorf("Expected %v frames, got %v", 1, len(s.Frames()))
	}
	s.AddFrame(UnsynchronisedLyricsFrame{Language: "B", ContentDescriptor: "B"})
	if len(s.Frames()) != 2 {
		t.Errorf("Expected %v frames, got %v", 2, len(s.Frames()))
	}
}
