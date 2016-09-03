// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"io"

	"github.com/bogem/id3v2/util"
)

type UnknownFrame struct {
	body io.Reader
}

func (uk UnknownFrame) Body() []byte {
	bufRd := util.NewReader(uk.body)
	body, err := bufRd.ReadAll()
	if err != nil {
		panic(err)
	}
	return body
}

func parseUnknownFrame(rd io.Reader) (Framer, error) {
	return UnknownFrame{body: rd}, nil
}
