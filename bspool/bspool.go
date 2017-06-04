// Copyright 2017 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package bspool is used for pooling bytes slices.
// This package is used only for internal usage in library id3v2. Users of
// library id3v2 must not use it.
package bspool

import (
	"sync"
)

var bsPool = sync.Pool{
	New: func() interface{} { return nil },
}

// Get returns []byte with len >= minSize.
func Get(minSize int) []byte {
	fromPool := bsPool.Get()
	if fromPool == nil {
		return make([]byte, minSize)
	}
	bs := fromPool.([]byte)
	if len(bs) < minSize {
		bs = make([]byte, minSize)
	}
	return bs
}

// Put puts b to pool.
func Put(b []byte) {
	bsPool.Put(b)
}
