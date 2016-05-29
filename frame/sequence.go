package frame

// Interface Sequence is used to manipulate with frames, which can be used more
// than one time (e.g. APIC, COMM, USLT, SYLT and etc.)
type Sequencer interface {
	Frames() []Framer
}
