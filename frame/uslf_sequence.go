// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package frame

// USLFSequence stores several USLT frames.
// Key for USLFSequence is language and content descriptor,
// so there is only one USLT frame with the same language and
// content descriptor.
//
// It's only needed for internal operations. Users of library id3v2 should not
// use any sequence in no case.
type USLFSequence struct {
	sequence map[string]UnsynchronisedLyricsFramer
}

func NewUSLFSequence() Sequencer {
	return &USLFSequence{
		sequence: make(map[string]UnsynchronisedLyricsFramer),
	}
}

func (us USLFSequence) Frames() []Framer {
	var (
		i      = 0
		frames = make([]Framer, len(us.sequence))
	)

	for _, v := range us.sequence {
		frames[i] = v
		i++
	}
	return frames
}

func (us USLFSequence) USLF(language string, contentDescriptor string) UnsynchronisedLyricsFramer {
	return us.sequence[language+contentDescriptor]
}

func (us *USLFSequence) AddFrame(f Framer) {
	uslf := f.(UnsynchronisedLyricsFramer)
	id := uslf.Language() + uslf.ContentDescriptor()
	us.sequence[id] = uslf
}
