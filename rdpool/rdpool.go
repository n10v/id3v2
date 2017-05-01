// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package rdpool is used only for internal usage in library id3v2. Users of
// library id3v2 must not use it.
package rdpool

import (
	"io"
	"sync"

	"github.com/bogem/id3v2/util"
)

var readerPool = sync.Pool{
	New: func() interface{} { return util.NewReader(nil) },
}

// Get returns *util.Reader with specified rd.
func Get(rd io.Reader) *util.Reader {
	reader := readerPool.Get().(*util.Reader)
	reader.Reset(rd)
	return reader
}

// Put puts rd back to pool.
func Put(rd *util.Reader) {
	readerPool.Put(rd)
}
