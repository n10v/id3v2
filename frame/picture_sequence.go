package frame

// PictureSequence stores several picture frames.
// Key for PictureSequence is a picture type code,
// so there is only one picture with the same picture type.
//
// It's only needed for internal operations. Users of library id3v2 should not
// use any sequence in no case.
type PictureSequence struct {
	sequence map[int]PictureFramer
}

func NewPictureSequence() Sequencer {
	return &PictureSequence{
		sequence: make(map[int]PictureFramer),
	}
}

func (ps PictureSequence) Frames() []Framer {
	var (
		i      = 0
		frames = make([]Framer, len(ps.sequence))
	)

	for _, v := range ps.sequence {
		frames[i] = v
		i++
	}
	return frames
}

func (ps PictureSequence) Picture(pt byte) PictureFramer {
	return ps.sequence[int(pt)]
}

func (ps *PictureSequence) AddFrame(f Framer) {
	pf := f.(PictureFramer)
	pt := pf.PictureType()
	ps.sequence[int(pt)] = pf
}
