// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// CommentSequence stores several comment frames.
// Key for CommentSequence is language and description,
// so there is only one comment frame with the same language and
// description.
//
// It's only needed for internal operations. Users of library id3v2 should not
// use any sequence in no case.
type CommentSequence struct {
	sequence map[string]CommentFrame
}

func NewCommentSequence() Sequencer {
	return &CommentSequence{
		sequence: make(map[string]CommentFrame),
	}
}

func (cs CommentSequence) Frames() []Framer {
	frames := make([]Framer, 0, len(cs.sequence))
	for _, f := range cs.sequence {
		frames = append(frames, f)
	}
	return frames
}

func (cs *CommentSequence) AddFrame(f Framer) {
	cf := f.(CommentFrame)
	id := cf.Language + cf.Description
	cs.sequence[id] = cf
}
