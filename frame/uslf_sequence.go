package frame

type USLFSequencer interface {
	Sequencer

	USLF(language string, contentDescriptor string) UnsynchronisedLyricsFramer
	AddUSLF(UnsynchronisedLyricsFramer)
}

type USLFSequence struct {
	sequence map[string]UnsynchronisedLyricsFramer
}

func NewUSLFSequence() *USLFSequence {
	return &USLFSequence{
		sequence: make(map[string]UnsynchronisedLyricsFramer),
	}
}

func (us USLFSequence) Frames() []Framer {
	frames := []Framer{}
	for _, f := range us.sequence {
		frames = append(frames, f)
	}
	return frames
}

func (us USLFSequence) USLF(language string, contentDescriptor string) UnsynchronisedLyricsFramer {
	return us.sequence[language+contentDescriptor]
}

func (us *USLFSequence) AddUSLF(uslf UnsynchronisedLyricsFramer) {
	id := uslf.Language() + uslf.ContentDescriptor()
	us.sequence[id] = uslf
}
