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

func (tf TextFrame) Form() []byte {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()

	b.WriteByte(util.NativeEncoding)
	b.WriteString(tf.Text())

	bytesBufPool.Put(b)
	return b.Bytes()
}

func (tf *TextFrame) SetText(text string) {
	tf.textBuffer.Reset()
	tf.textBuffer.WriteString(text)
}

func (tf TextFrame) Text() string {
	return tf.textBuffer.String()
}
