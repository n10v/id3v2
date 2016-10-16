// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// commentSequence stores several unique comment frames.
// Key for commentSequence is language and description,
// so there is only one comment frame with the same language and
// description.
//
// ID3v2 Documentation: "There may be more than one comment frame in each tag,
// but only one with the same language and content descriptor."
type commentSequence struct {
	sequence map[string]CommentFrame
}

func newCommentSequence() sequencer {
	return &commentSequence{
		sequence: make(map[string]CommentFrame),
	}
}

func (cs *commentSequence) AddFrame(f Framer) {
	cf := f.(CommentFrame)
	id := cf.Language + cf.Description
	cs.sequence[id] = cf
}

func (cs commentSequence) Frames() []Framer {
	frames := make([]Framer, 0, len(cs.sequence))
	for _, f := range cs.sequence {
		frames = append(frames, f)
	}
	return frames
}

func (cs commentSequence) Len() int {
	return len(cs.sequence)
}
