// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package bwpool is used for pooling bufio.Writers.
// This package is used only for internal usage in library id3v2. Users of
// library id3v2 must not use it.
package bwpool

import (
	"bufio"
	"io"
	"sync"
)

var bwPool = sync.Pool{
	New: func() interface{} { return bufio.NewWriter(nil) },
}

// Get returns *bufio.Writer with specified w.
func Get(w io.Writer) *bufio.Writer {
	bw := bwPool.Get().(*bufio.Writer)
	bw.Reset(w)
	return bw
}

// Put puts bw back to pool.
func Put(bw *bufio.Writer) {
	bwPool.Put(bw)
}
