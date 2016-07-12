// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// Framer provides a generic interface for frames.
type Framer interface {
	Bytes() ([]byte, error)
}
