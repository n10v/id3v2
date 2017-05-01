// Copyright 2017 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package lrpool is used for pooling io.LimitedReaders.
// This package is used only for internal usage in library id3v2. Users of
// library id3v2 must not use it.
package lrpool

import (
	"io"
	"sync"
)

var lrPool = sync.Pool{
	New: func() interface{} { return new(io.LimitedReader) },
}

// Get returns *io.LimitedReader with specified rd and n from pool.
func Get(rd io.Reader, n int64) *io.LimitedReader {
	r := lrPool.Get().(*io.LimitedReader)
	r.R = rd
	r.N = n
	return r
}

// Put puts r back to pool.
func Put(r *io.LimitedReader) {
	r.N = 0
	r.R = nil
	lrPool.Put(r)
}
