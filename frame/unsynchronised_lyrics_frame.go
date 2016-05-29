package frame

import (
	"bytes"
	"github.com/bogem/id3v2/util"
)

type UnsynchronisedLyricsFramer interface {
	Framer

	Language() string
	SetLanguage(string)

	ContentDescriptor() string
	SetContentDescriptor(string)

	Lyrics() string
	SetLyrics(string)
}

type UnsynchronisedLyricsFrame struct {
	language          string
	contentDescriptor bytes.Buffer
	lyrics            bytes.Buffer
}

func (uslf UnsynchronisedLyricsFrame) Form() []byte {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()

	b.WriteByte(util.NativeEncoding)
	b.WriteString(uslf.language)
	b.WriteString(uslf.ContentDescriptor())
	b.WriteByte(0)
	b.WriteString(uslf.Lyrics())

	bytesBufPool.Put(b)
	return b.Bytes()
}

func (uslf UnsynchronisedLyricsFrame) Language() string {
	return uslf.language
}

func (uslf *UnsynchronisedLyricsFrame) SetLanguage(lang string) {
	uslf.language = lang
}

func (uslf UnsynchronisedLyricsFrame) ContentDescriptor() string {
	return uslf.contentDescriptor.String()
}

func (uslf *UnsynchronisedLyricsFrame) SetContentDescriptor(cd string) {
	uslf.contentDescriptor.Reset()
	uslf.contentDescriptor.WriteString(cd)
}

func (uslf UnsynchronisedLyricsFrame) Lyrics() string {
	return uslf.lyrics.String()
}

func (uslf *UnsynchronisedLyricsFrame) SetLyrics(lyrics string) {
	uslf.lyrics.Reset()
	uslf.lyrics.WriteString(lyrics)
}
