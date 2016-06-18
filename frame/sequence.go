// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package frame

// Sequencer is used to manipulate with frames, which can be in tag
// more than one (e.g. APIC, COMM, USLT, SYLT and etc.)
//
// It's only needed for internal operations. Users of library id3v2 must not
// use any sequence in no case
type Sequencer interface {
	AddFrame(Framer)
	Frames() []Framer
}
