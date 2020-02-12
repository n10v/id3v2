// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"io"
	"math/rand"
	"strconv"
	"time"
)

// UnknownFrame is used for frames, which id3v2 so far doesn't know how to
// parse and write it. It just contains an unparsed byte body of the frame.
type UnknownFrame struct {
	Body []byte
}

var uidGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

func (uf UnknownFrame) UniqueIdentifier() string {
	// All unknown frames should have unique identifier, because we actually don't know their real identifiers.
	return strconv.Itoa(uidGenerator.Int())
}

func (uf UnknownFrame) Size() int {
	return len(uf.Body)
}

func (uf UnknownFrame) WriteTo(w io.Writer) (n int64, err error) {
	i, err := w.Write(uf.Body)
	return int64(i), err
}

func parseUnknownFrame(br *bufReader) (Framer, error) {
	body := br.ReadAll()
	return UnknownFrame{Body: body}, br.Err()
}
