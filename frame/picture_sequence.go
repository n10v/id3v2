package frame

type PictureSequencer interface {
	Sequencer

	AddPicture(PictureFramer)
	Picture(pt byte) PictureFramer
}

// PictureSequence stores several APICs and implements interface Framer.
// Key for PictureSequence is a key for PictureType array,
// so there is only one picture with the same picture type
type PictureSequence struct {
	sequence map[int]PictureFramer
}

func NewPictureSequence() *PictureSequence {
	return &PictureSequence{
		sequence: make(map[int]PictureFramer),
	}
}

func (as PictureSequence) Frames() []Framer {
	frames := []Framer{}
	for _, f := range as.sequence {
		frames = append(frames, f)
	}
	return frames
}

func (as *PictureSequence) AddPicture(pic PictureFramer) {
	pt := pic.PictureType()
	as.sequence[int(pt)] = pic
}

func (as PictureSequence) Picture(pt byte) PictureFramer {
	return as.sequence[int(pt)]
}
