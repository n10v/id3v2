// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package frame

// CommentSequence stores several comment frames.
// Key for CommentSequence is language and description,
// so there is only one comment frame with the same language and
// description.
//
// It's only needed for internal operations. Users of library id3v2 should not
// use any sequence in no case.
type CommentSequence struct {
	sequence map[string]CommentFramer
}

func NewCommentSequence() Sequencer {
	return &CommentSequence{
		sequence: make(map[string]CommentFramer),
	}
}

func (cs CommentSequence) Frames() []Framer {
	var (
		i      = 0
		frames = make([]Framer, len(cs.sequence))
	)

	for _, v := range cs.sequence {
		frames[i] = v
		i++
	}
	return frames
}

func (cs CommentSequence) Comment(language string, description string) CommentFramer {
	return cs.sequence[language+description]
}

func (cs *CommentSequence) AddFrame(f Framer) {
	cf := f.(CommentFramer)
	id := cf.Language() + cf.Description()
	cs.sequence[id] = cf
}
