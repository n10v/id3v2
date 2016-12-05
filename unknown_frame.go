// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"io"

	"github.com/bogem/id3v2/rdpool"
)

type UnknownFrame struct {
	body []byte
}

func (uk UnknownFrame) Body() []byte {
	return uk.body
}

func (uk UnknownFrame) Size() int {
	return len(uk.body)
}

func (uk UnknownFrame) WriteTo(w io.Writer) (n int, err error) {
	return w.Write(uk.body)
}

func parseUnknownFrame(rd io.Reader) (Framer, error) {
	bufRd := rdpool.Get(rd)
	defer rdpool.Put(bufRd)

	body, err := bufRd.ReadAll()

	return UnknownFrame{body: body}, err
}
