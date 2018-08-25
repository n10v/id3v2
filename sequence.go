// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"sync"
)

// sequence is used to manipulate with frames, which can be in tag
// more than one (e.g. APIC, COMM, USLT and etc.)
type sequence struct {
	framers     map[string]Framer
	framesCache []Framer
}

func (s *sequence) AddFrame(f Framer) {
	s.framesCache = s.framesCache[:0]

	var id string
	if cf, ok := f.(CommentFrame); ok {
		id = cf.Language + cf.Description
	} else if pf, ok := f.(PictureFrame); ok {
		id = pf.Description
	} else if uslf, ok := f.(UnsynchronisedLyricsFrame); ok {
		id = uslf.Language + uslf.ContentDescriptor
	} else if udtf, ok := f.(UserDefinedTextFrame); ok {
		id = udtf.Description
	} else {
		panic("sequence: unknown type of Framer")
	}

	s.framers[id] = f
}

func (s *sequence) Count() int {
	return len(s.framers)
}

func (s *sequence) Frames() []Framer {
	if len(s.framesCache) == 0 {
		for _, f := range s.framers {
			s.framesCache = append(s.framesCache, f)
		}
	}
	return s.framesCache
}

var seqPool = sync.Pool{New: func() interface{} {
	return &sequence{framers: make(map[string]Framer)}
}}

func getSequence() *sequence {
	s := seqPool.Get().(*sequence)
	if s.Count() > 0 {
		s.framers = make(map[string]Framer)
		s.framesCache = s.framesCache[:0]
	}
	return s
}

func putSequence(s *sequence) {
	seqPool.Put(s)
}
