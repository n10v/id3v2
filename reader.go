// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

// reader is used for convenient parsing of frames.
type reader struct {
	buf *bufio.Reader
}

// NewReader returns *Reader with specified rd.
func newReader(rd io.Reader) *reader {
	return &reader{buf: bufio.NewReader(rd)}
}

// Discard skips the next n bytes, returning the number of bytes discarded.
// If Discard skips fewer than n bytes, it also returns an error.
func (r *reader) Discard(n int) (discarded int, err error) {
	return r.buf.Discard(n)
}

// Read reads data into p.
// It returns the number of bytes read into p.
// The bytes are taken from at most one Read on the underlying Reader,
// hence n may be less than len(p).
// At EOF, the count will be zero and err will be io.EOF.
func (r *reader) Read(p []byte) (n int, err error) {
	return r.buf.Read(p)
}

// ReadByte reads and returns a single byte.
// If no byte is available, returns an error.
func (r *reader) ReadByte() (byte, error) {
	return r.buf.ReadByte()
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (r *reader) Next(n int) ([]byte, error) {
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
func (r *reader) ReadTillDelim(delim byte) ([]byte, error) {
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
func (r *reader) ReadTillDelims(delims []byte) ([]byte, error) {
	if len(delims) == 0 {
		return nil, errors.New("delims is empty")
	}
	if len(delims) == 1 {
		return r.ReadTillDelim(delims[0])
	}

	result := make([]byte, 0)

	for {
		read, err := r.ReadTillDelim(delims[0])
		if err != nil {
			return result, err
		}
		result = append(result, read...)

		peeked, err := r.buf.Peek(len(delims))
		if err != nil {
			return result, err
		}

		if bytes.Equal(peeked, delims) {
			break
		}

		b, err := r.ReadByte()
		if err != nil {
			return result, err
		}
		result = append(result, b)
	}

	return result, nil
}

// Reset discards any buffered data, resets all state,
// and switches the buffered reader to read from r.
func (r *reader) Reset(rd io.Reader) {
	r.buf.Reset(rd)
}
