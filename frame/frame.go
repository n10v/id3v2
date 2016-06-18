// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package frame

import (
	"bytes"
	"sync"
)

var bytesBufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

// Framer provides a generic interface for frames.
type Framer interface {
	Bytes() ([]byte, error)
}
