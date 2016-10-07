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
	New: func() interface{} { return nil },
}

func Get(rd io.Reader) *util.Reader {
	var reader *util.Reader

	ireader := readerPool.Get()
	if ireader == nil {
		reader = util.NewReader(rd)
	} else {
		reader = ireader.(*util.Reader)
		reader.Reset(rd)
	}

	return reader
}

func Put(b *util.Reader) {
	readerPool.Put(b)
}
