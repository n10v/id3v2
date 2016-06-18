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
