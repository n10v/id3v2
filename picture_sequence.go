// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// PictureSequence stores several picture frames.
// Key for PictureSequence is a picture type code,
// so there is only one picture with the same picture type.
//
// It's only needed for internal operations. Users of library id3v2 should not
// use any sequence in no case.
type PictureSequence struct {
	sequence map[int]PictureFrame
}

func NewPictureSequence() Sequencer {
	return &PictureSequence{
		sequence: make(map[int]PictureFrame),
	}
}

func (ps PictureSequence) Frames() []Framer {
	frames := make([]Framer, 0, len(ps.sequence))
	for _, f := range ps.sequence {
		frames = append(frames, f)
	}
	return frames
}

func (ps *PictureSequence) AddFrame(f Framer) {
	pf := f.(PictureFrame)
	pt := pf.PictureType
	ps.sequence[int(pt)] = pf
}
