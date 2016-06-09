package frame

type USLFSequencer interface {
	Sequencer

	USLF(language string, contentDescriptor string) UnsynchronisedLyricsFramer
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
