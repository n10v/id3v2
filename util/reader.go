package util

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
)

type Reader struct {
	buf *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{buf: bufio.NewReader(rd)}
}

// ReadAll reads from Reader until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from Reader until EOF, it does not treat an EOF from Read
// as an error to be reported.
func (r *Reader) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(r.buf)
}

// ReadByte reads and returns a single byte.
// If no byte is available, returns an error.
func (r *Reader) ReadByte() (byte, error) {
	return r.buf.ReadByte()
}

// ReadBytes reads n bytes.
func (r *Reader) ReadBytes(n int) ([]byte, error) {
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
// returning a slice containing the data up to the delimiter.
// If ReadTillDelim encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadTillDelim returns err != nil if and only if the returned data does not end in
// delim.
func (r *Reader) ReadTillDelim(delim byte) ([]byte, error) {
	read, err := r.buf.ReadBytes(delim)
	if read == nil || len(read) == 0 {
		return read, err
	}
	return read[:len(read)-1], err
}

// ReadTillDelims reads until the first occurrence of delims in the input,
// returning a slice containing the data up to the delimiters.
// If ReadTillDelims encounters an error before finding a delimiters,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadTillAndWithDelims returns err != nil if and only if the returned data does not end in
// delim.
func (r *Reader) ReadTillDelims(delims []byte) ([]byte, error) {
	if len(delims) == 0 {
		return r.ReadAll()
	}
	if len(delims) == 1 {
		return r.ReadTillDelim(delims[0])
	}

	buf := make([]byte, 0)

	for {
		b, err := r.buf.ReadByte()
		if err != nil {
			return buf, err
		}

		if b == delims[0] {
			peeked, err := r.buf.Peek(len(delims) - 1)
			if err != nil {
				return buf, err
			}
			if bytes.Equal(peeked, delims[1:]) {
				r.buf.Discard(len(peeked) - 1)
				break
			}
		}
		buf = append(buf, b)
	}

	return buf, nil
}
