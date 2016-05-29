package frame

import "bytes"

type APICSequencer interface {
	Framer

	AddPicture(PictureFramer)
	Picture(pt byte) PictureFramer
}

// APICSequence stores several APICs and implements interface Framer.
// Key for APICSequnce is a key for PictureType array,
// so there is only one picture with the same picture type
type APICSequence struct {
	sequence map[int]PictureFramer
}

func NewAPICSequence() *APICSequence {
	return &APICSequence{
		sequence: make(map[int]PictureFramer),
	}
}

func (as APICSequence) Form() []byte {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()
	for _, pf := range as.sequence {
		frame := pf.Form()
		b.Write(frame)
	}
	bytesBufPool.Put(b)
	return b.Bytes()
}

func (as *APICSequence) AddPicture(pic PictureFramer) {
	pt := pic.PictureType()
	as.sequence[int(pt)] = pic
}

func (as APICSequence) Picture(pt byte) PictureFramer {
	return as.sequence[int(pt)]
}
