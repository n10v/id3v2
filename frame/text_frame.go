package frame

import (
	"bytes"
	"github.com/bogem/id3v2/util"
)

type TextFramer interface {
	Framer

	Text() string
	SetText(string)
}

type TextFrame struct {
	textBuffer bytes.Buffer
}

func (tf TextFrame) Bytes() ([]byte, error) {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()
	defer bytesBufPool.Put(b)

	b.WriteByte(util.NativeEncoding)
	b.WriteString(tf.Text())

	return b.Bytes(), nil
}

func (tf *TextFrame) SetText(text string) {
	tf.textBuffer.Reset()
	tf.textBuffer.WriteString(text)
}

func (tf TextFrame) Text() string {
	return tf.textBuffer.String()
}
