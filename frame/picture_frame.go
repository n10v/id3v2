package frame

import (
	"bytes"
	"errors"
	"github.com/bogem/id3v2/util"
	"io"
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

	Picture() io.Reader
	SetPicture(io.Reader)

	PictureType() string
	SetPictureType(string) error
}

type PictureFrame struct {
	description string
	mimeType    string
	picture     io.Reader
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

	b.ReadFrom(pf.picture)
	if v, ok := pf.picture.(*os.File); ok {
		v.Close()
	}

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

func (pf PictureFrame) Picture() io.Reader {
	return pf.picture
}

func (pf *PictureFrame) SetPicture(rd io.Reader) {
	pf.picture = rd
}

func (pf *PictureFrame) SetPictureFromFile(name string) error {
	if pf.picture != nil {
		if v, ok := pf.picture.(*os.File); ok {
			v.Close()
		}
	}
	pictureFile, err := os.Open(name)
	pf.picture = pictureFile
	return err
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
