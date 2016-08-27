// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// sequencer is used to manipulate with frames, which can be in tag
// more than one (e.g. APIC, COMM, USLT, SYLT and etc.)
// Every object that implements sequencer interface must ensure the uniqueness
// of each frame.
// For example, in id3v2.4 documentation it says, that "there may be more than
// one comment frame in each tag, but only one with the same language and
// content descriptor.", so commentSequence must ensure that condition.
type sequencer interface {
	AddFrame(Framer)
	Frames() []Framer
}
