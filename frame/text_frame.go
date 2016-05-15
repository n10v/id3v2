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
	id         string
	textBuffer bytes.Buffer
}

func (tf TextFrame) Form() ([]byte, error) {
	var b bytes.Buffer
	if header, err := FormFrameHeader(&tf); err != nil {
		return nil, err
	} else {
		b.Write(header)
	}
	b.WriteByte(util.NativeEncoding)
	b.WriteString(tf.Text())
	return b.Bytes(), nil
}

func (tf TextFrame) ID() string {
	return tf.id
}

func (tf *TextFrame) SetID(id string) {
	tf.id = id
}

func (tf TextFrame) Size() uint32 {
	encodingSize := 1
	size := uint32(encodingSize + tf.textBuffer.Len())
	return size
}

func (tf *TextFrame) SetText(text string) {
	tf.textBuffer.Reset()
	tf.textBuffer.WriteString(text)
}

func (tf TextFrame) Text() string {
	return tf.textBuffer.String()
}
