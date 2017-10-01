// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"errors"
	"io"
)

var ErrInvalidLanguageLength = errors.New("language code must consist of three letters according to ISO 639-2")

// Framer provides a generic interface for frames.
// You can create your own frames. They must implement only this interface.
type Framer interface {
	// Size returns the size of frame body.
	Size() int

	// WriteTo writes body slice into w.
	WriteTo(w io.Writer) (n int64, err error)
}
