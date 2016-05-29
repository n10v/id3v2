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

func (uslt UnsynchronisedLyricsFrame) Form() []byte {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()

	b.WriteByte(util.NativeEncoding)
	b.WriteString(uslt.language)
	b.WriteString(uslt.ContentDescriptor())
	b.WriteByte(0)
	b.WriteString(uslt.Lyrics())

	bytesBufPool.Put(b)
	return b.Bytes()
}

func (uslt UnsynchronisedLyricsFrame) Language() string {
	return uslt.language
}

func (uslt *UnsynchronisedLyricsFrame) SetLanguage(lang string) {
	uslt.language = lang
}

func (uslt UnsynchronisedLyricsFrame) ContentDescriptor() string {
	return uslt.contentDescriptor.String()
}

func (uslt *UnsynchronisedLyricsFrame) SetContentDescriptor(cd string) {
	uslt.contentDescriptor.Reset()
	uslt.contentDescriptor.WriteString(cd)
}

func (uslt UnsynchronisedLyricsFrame) Lyrics() string {
	return uslt.lyrics.String()
}

func (uslt *UnsynchronisedLyricsFrame) SetLyrics(lyrics string) {
	uslt.lyrics.Reset()
	uslt.lyrics.WriteString(lyrics)
}
