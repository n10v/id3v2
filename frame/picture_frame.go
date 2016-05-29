package frame

import (
	"bytes"
	"github.com/bogem/id3v2/util"
	"io"
	"os"
)

type PictureFramer interface {
	Framer

	Description() string
	SetDescription(string)

	MimeType() string
	SetMimeType(string)

	Picture() io.Reader
	SetPicture(io.Reader)

	PictureType() byte
	SetPictureType(byte)
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

func (pf PictureFrame) PictureType() byte {
	return pf.pictureType
}

func (pf *PictureFrame) SetPictureType(pt byte) {
	pf.pictureType = pt
}
