// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

// Reader is used for convenient parsing of frames.
type Reader struct {
	buf      *bufio.Reader
	bytesBuf *bytes.Buffer // Need for intermediate calculations
}

// NewReader returns *Reader with specified rd.
func NewReader(rd io.Reader) *Reader {
	return &Reader{buf: bufio.NewReader(rd)}
}

func (r *Reader) initBytesBuf() {
	if r.bytesBuf == nil {
		r.bytesBuf = new(bytes.Buffer)
	}
	r.bytesBuf.Reset()
}

// Discard skips the next n bytes, returning the number of bytes discarded.
// If Discard skips fewer than n bytes, it also returns an error.
func (r *Reader) Discard(n int) (discarded int, err error) {
	return r.buf.Discard(n)
}

// Read reads data into p.
// It returns the number of bytes read into p.
// The bytes are taken from at most one Read on the underlying Reader,
// hence n may be less than len(p).
// At EOF, the count will be zero and err will be io.EOF.
func (r *Reader) Read(p []byte) (n int, err error) {
	return r.buf.Read(p)
}

// ReadByte reads and returns a single byte.
// If no byte is available, returns an error.
func (r *Reader) ReadByte() (byte, error) {
	return r.buf.ReadByte()
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (r *Reader) Next(n int) ([]byte, error) {
	if n == 0 {
		return nil, nil
	}

	peeked, err := r.buf.Peek(n)
	if err != nil {
		return nil, err
	}

	if _, err := r.buf.Discard(n); err != nil {
		return nil, err
	}

	return peeked, nil
}

// ReadTillDelim reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and NOT including the delimiter.
// If ReadTillDelim encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadTillDelim returns err != nil if and only if ReadTillDelim didn't find
// delim.
func (r *Reader) ReadTillDelim(delim byte) ([]byte, error) {
	read, err := r.buf.ReadBytes(delim)
	if err != nil || len(read) == 0 {
		return read, err
	}
	err = r.buf.UnreadByte()
	return read[:len(read)-1], err
}

// ReadTillDelims reads until the first occurrence of delims in the input,
// returning a slice containing the data up to and NOT including the delimiters.
// If ReadTillDelims encounters an error before finding a delimiters,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadTillDelims returns err != nil if and only if ReadTillDelims didn't find
// delims.
func (r *Reader) ReadTillDelims(delims []byte) ([]byte, error) {
	if len(delims) == 0 {
		return nil, errors.New("delims is empty")
	}
	if len(delims) == 1 {
		return r.ReadTillDelim(delims[0])
	}

	r.initBytesBuf()

	for {
		read, err := r.ReadTillDelim(delims[0])
		if err != nil {
			return r.bytesBuf.Bytes(), err
		}
		r.bytesBuf.Write(read)

		peeked, err := r.buf.Peek(len(delims))
		if err != nil {
			return r.bytesBuf.Bytes(), err
		}

		if bytes.Equal(peeked, delims) {
			break
		}

		b, err := r.ReadByte()
		if err != nil {
			return r.bytesBuf.Bytes(), err
		}
		r.bytesBuf.WriteByte(b)
	}

	return r.bytesBuf.Bytes(), nil
}

// String returns the contents of the unread portion of the buffered data
// as a string. It returns error if there was an error during read.
func (r *Reader) String() (string, error) {
	r.initBytesBuf()
	_, err := r.bytesBuf.ReadFrom(r)
	return r.bytesBuf.String(), err
}

// Reset discards any buffered data, resets all state,
// and switches the buffered reader to read from r.
func (r *Reader) Reset(rd io.Reader) {
	r.buf.Reset(rd)
}
