// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import "io"

// Framer provides a generic interface for frames.
// You can create your own frames. They must implement only this interface.
type Framer interface {
	// Body returns byte slice, that contains only frame body. Used in tests.
	Body() []byte

	// Size returns the size of frame body.
	Size() int

	// WriteTo writes body slice into io.Writer.
	WriteTo(w io.Writer) (n int64, err error)
}
