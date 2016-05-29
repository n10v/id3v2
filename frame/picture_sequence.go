package frame

type PictureSequencer interface {
	Sequencer

	Picture(pt byte) PictureFramer
	AddPicture(PictureFramer)
}

// PictureSequence stores several pictures and implements interface Framer.
// Key for PictureSequence is a picture type code,
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

func (as PictureSequence) Picture(pt byte) PictureFramer {
	return as.sequence[int(pt)]
}

func (as *PictureSequence) AddPicture(pic PictureFramer) {
	pt := pic.PictureType()
	as.sequence[int(pt)] = pic
}
