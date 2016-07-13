// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// Framer provides a generic interface for frames.
// You can create your own frames. They must implement only this interface.
type Framer interface {
	// Body returns byte slice, that contains only frame body
	Body() []byte
}
