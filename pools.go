// Copyright 2017 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

// bsPool is a pool of byte slices.
var bsPool = sync.Pool{
	New: func() interface{} { return nil },
}

// getByteSlice returns []byte with len >= minSize.
func getByteSlice(minSize int) []byte {
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

// putByteSlice puts b to pool.
func putByteSlice(b []byte) {
	bsPool.Put(b)
}

// bwPool is a pool of *bufio.Writer.
var bwPool = sync.Pool{
	New: func() interface{} { return bufio.NewWriter(nil) },
}

// getBufioWriter returns *bufio.Writer with specified w.
func getBufioWriter(w io.Writer) *bufio.Writer {
	bw := bwPool.Get().(*bufio.Writer)
	bw.Reset(w)
	return bw
}

// putBufioWriter puts bw back to pool.
func putBufioWriter(bw *bufio.Writer) {
	bwPool.Put(bw)
}

// lrPool is a pool of *io.LimitedReader.
var lrPool = sync.Pool{
	New: func() interface{} { return new(io.LimitedReader) },
}

// getLimitedReader returns *io.LimitedReader with specified rd and n from pool.
func getLimitedReader(rd io.Reader, n int64) *io.LimitedReader {
	r := lrPool.Get().(*io.LimitedReader)
	r.R = rd
	r.N = n
	return r
}

// putLimitedReader puts r back to pool.
func putLimitedReader(r *io.LimitedReader) {
	r.N = 0
	r.R = nil
	lrPool.Put(r)
}

// rdPool is a pool of *reader.
var rdPool = sync.Pool{
	New: func() interface{} { return newReader(nil) },
}

// getUtilReader returns *reader with specified rd.
func getUtilReader(rd io.Reader) *reader {
	reader := rdPool.Get().(*reader)
	reader.Reset(rd)
	return reader
}

// putUtilReader puts rd back to pool.
func putUtilReader(rd *reader) {
	rdPool.Put(rd)
}

// bbPool is a pool of *bytes.Buffer.
var bbPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

// getBytesBuffer returns empty *bytes.Buffer.
func getBytesBuffer() *bytes.Buffer {
	return bbPool.Get().(*bytes.Buffer)
}

// putBytesBuffer resets buf and puts to pool.
func putBytesBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bbPool.Put(buf)
}
