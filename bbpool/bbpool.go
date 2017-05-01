// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package bbpool is used for pooling bytes.Buffers.
// This package is used only for internal usage in library id3v2. Users of
// library id3v2 must not use it.
package bbpool

import (
	"bytes"
	"sync"
)

var bbPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

// Get returns *bytes.Buffer from pool.
func Get() *bytes.Buffer {
	b := bbPool.Get().(*bytes.Buffer)
	return b
}

// Put puts b back to pool.
func Put(b *bytes.Buffer) {
	b.Reset()
	bbPool.Put(b)
}
