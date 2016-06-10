package frame

// Interface Sequence is used to manipulate with frames, which can be in tag
// more than one (e.g. APIC, COMM, USLT, SYLT and etc.)
type Sequencer interface {
	AddFrame(Framer)
	Frames() []Framer
}
