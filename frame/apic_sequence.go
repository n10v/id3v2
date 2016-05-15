package frame

import (
	"bytes"
	"errors"
)

type APICSequencer interface {
	Framer

	AddPicture(PictureFramer)
	Picture(picType string) (PictureFramer, error)
}

// APICSequence stores several APICs and implements interface Framer.
// Key for APICSequnce is a key for PictureType array,
// so there is only one picture with the same picture type
type APICSequence struct {
	sequence map[int]PictureFramer
	id       string
}

func NewAPICSequence() *APICSequence {
	return &APICSequence{
		sequence: make(map[int]PictureFramer),
	}
}

func (as APICSequence) Form() ([]byte, error) {
	var b bytes.Buffer
	for _, pf := range as.sequence {
		if frame, err := pf.Form(); err != nil {
			return nil, err
		} else {
			b.Write(frame)
		}
	}
	return b.Bytes(), nil
}

func (as APICSequence) ID() string {
	return as.id
}

func (as *APICSequence) SetID(id string) {
	as.id = id
}

func (as APICSequence) Size() (size uint32) {
	for _, pf := range as.sequence {
		size += pf.Size()
	}
	return
}

func (as *APICSequence) AddPicture(pic PictureFramer) {
	for k, v := range PictureTypes {
		if v == pic.PictureType() {
			as.sequence[k] = pic
			break
		}
	}
}

func (as APICSequence) Picture(picType string) (PictureFramer, error) {
	for k, v := range PictureTypes {
		if v == picType {
			return as.sequence[k], nil
		}
	}
	return &PictureFrame{}, errors.New("Unsupported picture type")
}
