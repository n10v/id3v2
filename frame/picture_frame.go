package frame

import (
	"bytes"
	"errors"
	"github.com/bogem/id3v2/util"
	"os"
)

var (
	PictureTypes = [...]string{
		"Other",
		"File icon",
		"Other file icon",
		"Cover (front)",
		"Cover (back)",
		"Leaflet page",
		"Media",
		"Lead artist",
		"Artist",
		"Conductor",
		"Band/Orchestra",
		"Composer",
		"Lyricist/text writer",
		"Recording Location",
		"During recording",
		"During performance",
		"Movie/video screen capture",
		"A bright coloured fish",
		"Illustration",
		"Band/artist logotype",
		"Publisher/Studio logotype",
	}
)

type PictureFramer interface {
	Framer

	Description() string
	SetDescription(string)

	MimeType() string
	SetMimeType(string)

	Picture() []byte
	SetPicture([]byte) error

	PictureType() string
	SetPictureType(string) error
}

type PictureFrame struct {
	description string
	mimeType    string
	picture     bytes.Buffer
	pictureType byte
}

func (pf PictureFrame) Form() []byte {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()
	b.WriteByte(util.NativeEncoding)
	b.WriteString(pf.mimeType)
	b.WriteByte(0)
	b.WriteByte(pf.pictureType)
	b.WriteString(pf.description)
	b.WriteByte(0)
	b.Write(pf.Picture())
	bytesBufPool.Put(b)
	return b.Bytes()
}

func (pf PictureFrame) Description() string {
	return pf.description
}

func (pf *PictureFrame) SetDescription(desc string) {
	pf.description = desc
}

func (pf PictureFrame) MimeType() string {
	return pf.mimeType
}

func (pf *PictureFrame) SetMimeType(mt string) {
	pf.mimeType = mt
}

func (pf PictureFrame) Picture() []byte {
	return pf.picture.Bytes()
}

func (pf *PictureFrame) SetPicture(b []byte) error {
	pf.picture.Reset()
	if _, err := pf.picture.Write(b); err != nil {
		return err
	}
	return nil
}

func (pf *PictureFrame) SetPictureFromFile(file *os.File) error {
	pf.picture.Reset()
	if _, err := pf.picture.ReadFrom(file); err != nil {
		return err
	}
	return nil
}

func (pf PictureFrame) PictureType() string {
	return PictureTypes[pf.pictureType]
}

func (pf *PictureFrame) SetPictureType(pt string) error {
	for k, v := range PictureTypes {
		if v == pt {
			pf.pictureType = byte(k)
			return nil
		}
	}
	return errors.New("Unsupported picture type")
}
