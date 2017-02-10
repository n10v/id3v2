// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// uslfSequence stores several unique USLT frames.
// Key for uslfSequence is language and content descriptor,
// so there is only one USLT frame with the same language and
// content descriptor.
//
// ID3v2 Documentation: "There may be more than one 'Unsynchronised
// lyrics/text transcription' frame in each tag, but only one with the
// same language and content descriptor."
type uslfSequence struct {
	sequence    map[string]UnsynchronisedLyricsFrame
	framesCache []Framer
}

func newUSLFSequence() sequencer {
	return &uslfSequence{
		sequence: make(map[string]UnsynchronisedLyricsFrame),
	}
}

func (us *uslfSequence) AddFrame(f Framer) {
	uslf := f.(UnsynchronisedLyricsFrame)
	id := uslf.Language + uslf.ContentDescriptor
	us.sequence[id] = uslf
}

func (us uslfSequence) Count() int {
	return len(us.sequence)
}

func (us *uslfSequence) Frames() []Framer {
	cache := us.framesCache
	if cache == nil || len(cache) == 0 {
		cache = make([]Framer, 0, len(us.sequence))
		for _, f := range us.sequence {
			cache = append(cache, f)
		}
		us.framesCache = cache
	}
	return cache
}
