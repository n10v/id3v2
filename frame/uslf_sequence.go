package frame

type USLFSequencer interface {
	Sequencer

	USLF(id string, contentDescriptor string) UnsynchronisedLyricsFrame
	AddUSLF(uslf UnsynchronisedLyricsFrame)
}

type USLFSequence struct {
	sequence map[string]UnsynchronisedLyricsFramer
}
