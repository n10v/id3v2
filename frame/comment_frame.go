package frame

import (
	"bytes"
	"github.com/bogem/id3v2/util"
)

type CommentFramer interface {
	Framer

	Language() string
	SetLanguage(string)

	Description() string
	SetDescription(string)

	Text() string
	SetText(string)
}

type CommentFrame struct {
	language    string
	description bytes.Buffer
	text        bytes.Buffer
}

func (cf CommentFrame) Form() []byte {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()

	b.WriteByte(util.NativeEncoding)
	b.WriteString(cf.language)
	b.WriteString(cf.Description())
	b.WriteByte(0)
	b.WriteString(cf.Text())

	bytesBufPool.Put(b)
	return b.Bytes()
}

func (cf CommentFrame) Language() string {
	return cf.language
}

func (cf *CommentFrame) SetLanguage(lang string) {
	cf.language = lang
}

func (cf CommentFrame) Description() string {
	return cf.description.String()
}

func (cf *CommentFrame) SetDescription(d string) {
	cf.description.Reset()
	cf.description.WriteString(d)
}

func (cf CommentFrame) Text() string {
	return cf.text.String()
}

func (cf *CommentFrame) SetText(text string) {
	cf.text.Reset()
	cf.text.WriteString(text)
}
