// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package bbpool is used only for internal usage. Users of
// library id3v2 must not use it.
package bbpool

import (
	"bytes"
	"sync"
)

var bytesBufferPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func Get() *bytes.Buffer {
	b := bytesBufferPool.Get().(*bytes.Buffer)
	b.Reset()
	return b
}

func Put(b *bytes.Buffer) {
	bytesBufferPool.Put(b)
}
