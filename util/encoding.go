// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

// Encoding is a primitive for encodings.
//
// If you are user of id3v2 library, all list of allowed encodings you can find
// in documentation in variables.
// You can set encoding by yourself like this:
//
//	comment := id3v2.CommentFrame{
//		Encoding:    id3v2.EncodingISO,
//		Language:    "eng",
//		Description: string([]byte{68, 101, 115, 99}),
//		Text:        string([]byte{84, 101, 120, 116}),
//	}
//	tag.AddCommentFrame(comment)
type Encoding struct {
	Key              byte
	TerminationBytes []byte
}
