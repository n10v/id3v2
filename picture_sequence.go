// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// pictureSequence stores several picture frames.
// Key for pictureSequence is a picture type code,
// so there is only one picture with the same picture type.
//
// It's only needed for internal operations. Users of library id3v2 should not
// use any sequence in no case.
type pictureSequence struct {
	sequence map[int]PictureFrame
}

func newPictureSequence() sequencer {
	return &pictureSequence{
		sequence: make(map[int]PictureFrame),
	}
}

func (ps pictureSequence) Frames() []Framer {
	frames := make([]Framer, 0, len(ps.sequence))
	for _, f := range ps.sequence {
		frames = append(frames, f)
	}
	return frames
}

func (ps *pictureSequence) AddFrame(f Framer) {
	pf := f.(PictureFrame)
	pt := pf.PictureType
	ps.sequence[int(pt)] = pf
}
