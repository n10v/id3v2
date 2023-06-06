// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"bytes"
	"io"
)

// TextFrame is used to work with all text frames
// (all T*** frames like TIT2 (title), TALB (album) and so on).
type TextFrame struct {
	encoding         Encoding
	informationBytes []byte
}

func NewTextFrame(encoding Encoding, informationBytes []byte) *TextFrame {
	return &TextFrame{encoding: encoding, informationBytes: informationBytes}
}

// TODO: Docs
func (tf *TextFrame) GetEncoding() Encoding {
	return tf.encoding
}

func (tf *TextFrame) GetInformation() (string, error) {
	return decodeBytes(tf.informationBytes, tf.encoding)
}

func (tf *TextFrame) GetInformationBytes() []byte {
	return tf.informationBytes
}

// Information must be valid UTF-8 string.
func (tf *TextFrame) SetInformationFromString(information string) {
	tf.encoding = EncodingUTF8
	tf.informationBytes = []byte(information)
}

func (tf *TextFrame) SetInformationFromBytes(informationBytes []byte, encoding Encoding) {
	tf.encoding = encoding
	tf.informationBytes = informationBytes
}

func (tf *TextFrame) EncodeWith(encoding Encoding) error {
	var err error
	tf.informationBytes, err = encodeBytes(tf.informationBytes, encoding)
	if err != nil {
		return err
	}

	tf.encoding = encoding

	return nil
}

func (tf TextFrame) Size() int {
	encodingKeySize := 1

	terminationBytesLength := 0
	if tf.shouldAddTerminationBytes() {
		terminationBytesLength = len(tf.encoding.TerminationBytes)
	}

	return encodingKeySize + len(tf.informationBytes) + terminationBytesLength
}

func (tf TextFrame) UniqueIdentifier() string {
	return ""
}

func (tf TextFrame) WriteTo(w io.Writer) (int64, error) {
	bw := bufio.NewWriter(w)

	if err := bw.WriteByte(tf.encoding.Key); err != nil {
		return 0, err
	}

	if _, err := bw.Write(tf.informationBytes); err != nil {
		return 0, err
	}

	if tf.shouldAddTerminationBytes() {
		if _, err := bw.Write(tf.encoding.TerminationBytes); err != nil {
			return 0, err
		}
	}

	if err := bw.Flush(); err != nil {
		return 0, err
	}

	return int64(bw.Size()), nil
}

func (tf TextFrame) shouldAddTerminationBytes() bool {
	return !bytes.HasSuffix(tf.informationBytes, tf.encoding.TerminationBytes)
}

func parseTextFrame(body []byte) (Framer, error) {
	tf := TextFrame{
		encoding:         getEncoding(body[0]),
		informationBytes: body[1:],
	}

	return tf, nil
}
