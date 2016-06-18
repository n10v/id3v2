package frame

import (
	"bytes"
	"github.com/bogem/id3v2/util"
)

// TextFramer is used to work with all text frames
// (all T*** frames like TIT2, TALB and so on).
type TextFramer interface {
	Framer

	Encoding() util.Encoding
	SetEncoding(util.Encoding)

	Text() string
	SetText(string)
}

type TextFrame struct {
	encoding   util.Encoding
	textBuffer bytes.Buffer
}

func (tf TextFrame) Bytes() ([]byte, error) {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()
	defer bytesBufPool.Put(b)

	b.WriteByte(tf.encoding.Key)
	b.WriteString(tf.Text())

	return b.Bytes(), nil
}

func (tf TextFrame) Encoding() util.Encoding {
	return tf.encoding
}

func (tf *TextFrame) SetEncoding(e util.Encoding) {
	tf.encoding = e
}

func (tf TextFrame) Text() string {
	return tf.textBuffer.String()
}

func (tf *TextFrame) SetText(text string) {
	tf.textBuffer.Reset()
	tf.textBuffer.WriteString(text)
}
