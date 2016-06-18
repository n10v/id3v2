// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

// Primitive for encoding.
//
// If you are user of id3v2 library, all list of allowed encodings you can find
// in documentation in variables.
// For convenience, by default all frame constructors set UTF-8 encoding
// to frames.
// You can set by yourself encoding to frames via SetEncoding method.
// For example:
//	comm := id3v2.NewCommentFrame()
//	comm.SetEncoding(id3v2.ENISO)
//	comm.SetLanguage("eng")
//	comm.SetDescription(string([]byte{68, 101, 115, 99})) // "Desc" on ISO-8859-1
//	comm.SetText(string([]byte{84, 101, 120, 116})) // "Text" on ISO-8859-1
//	tag.AddCommentFrame(comm)
type Encoding struct {
	Key              byte
	TerminationBytes []byte
}
