// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// sequence is used to manipulate with frames, which can be in tag
// more than one (e.g. APIC, COMM, USLT and etc.)
type sequence struct {
	framers     map[string]Framer
	framesCache []Framer
}

func newSequence() *sequence {
	return &sequence{
		framers: make(map[string]Framer),
	}
}

func (s *sequence) AddFrame(f Framer) {
	s.framesCache = nil

	var id string
	if cf, ok := f.(CommentFrame); ok {
		id = cf.Language + cf.Description
	} else if pf, ok := f.(PictureFrame); ok {
		id = pf.Description
	} else if uslf, ok := f.(UnsynchronisedLyricsFrame); ok {
		id = uslf.Language + uslf.ContentDescriptor
	} else {
		panic("sequence: unknown type of Framer")
	}

	s.framers[id] = f
}

func (s *sequence) Count() int {
	return len(s.framers)
}

func (s *sequence) Frames() []Framer {
	cache := s.framesCache
	if len(cache) == 0 {
		cache = make([]Framer, 0, len(s.framers))
		for _, f := range s.framers {
			cache = append(cache, f)
		}
		s.framesCache = cache
	}
	return cache
}
