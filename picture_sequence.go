// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// pictureSequence stores several unique picture frames.
// Key for pictureSequence is a description,
// so there is only one picture with the same description.
//
// ID3v2 Documentation: "There may be several pictures attached to one file,
// each in their individual "APIC" frame, but only one with the same content
// descriptor.(TODO:) There may only be one picture with the picture type
// declared as picture type $01 and $02 respectively."
type pictureSequence struct {
	sequence    map[string]PictureFrame
	framesCache []Framer
}

func newPictureSequence() sequencer {
	return &pictureSequence{
		sequence: make(map[string]PictureFrame),
	}
}

func (ps *pictureSequence) AddFrame(f Framer) {
	ps.framesCache = nil

	pf := f.(PictureFrame)
	ps.sequence[pf.Description] = pf
}

func (ps pictureSequence) Count() int {
	return len(ps.sequence)
}

func (ps *pictureSequence) Frames() []Framer {
	cache := ps.framesCache
	if len(cache) == 0 {
		cache = make([]Framer, 0, len(ps.sequence))
		for _, f := range ps.sequence {
			cache = append(cache, f)
		}
		ps.framesCache = cache
	}
	return cache
}
