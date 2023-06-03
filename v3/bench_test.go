// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var frontCoverPicture = mustReadFile(frontCoverPath)

func BenchmarkParseAllFrames(b *testing.B) {
	writeTag(b, EncodingUTF8)
	musicContent := mustReadFile(mp3Path)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tag, err := ParseReader(bytes.NewReader(musicContent), parseOpts)
		if tag == nil || err != nil {
			b.Fatal("Error while opening mp3 file:", err)
		}
	}
}

func BenchmarkParseAllFramesISO(b *testing.B) {
	writeTag(b, EncodingISO)
	musicContent := mustReadFile(mp3Path)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tag, err := ParseReader(bytes.NewReader(musicContent), parseOpts)
		if tag == nil || err != nil {
			b.Fatal("Error while opening mp3 file:", err)
		}
	}
}

func BenchmarkParseArtistAndTitle(b *testing.B) {
	writeTag(b, EncodingUTF8)
	musicContent := mustReadFile(mp3Path)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		opts := Options{Parse: true, ParseFrames: []string{"Artist", "Title"}}
		tag, err := ParseReader(bytes.NewReader(musicContent), opts)
		if tag == nil || err != nil {
			b.Fatal("Error while opening mp3 file:", err)
		}
	}
}

func BenchmarkWrite(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchWrite(b, EncodingUTF8)
	}
}

func BenchmarkWriteISO(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchWrite(b, EncodingISO)
	}
}

func benchWrite(b *testing.B, encoding Encoding) {
	tag := NewEmptyTag()
	setFrames(tag, encoding)
	if _, err := tag.WriteTo(ioutil.Discard); err != nil {
		b.Error("Error while writing a tag:", err)
	}
}

func writeTag(b *testing.B, encoding Encoding) {
	tag, err := Open(mp3Path, Options{Parse: false})
	if tag == nil || err != nil {
		b.Fatal("Error while opening mp3 file:", err)
	}
	defer tag.Close()

	setFrames(tag, encoding)

	if err = tag.Save(); err != nil {
		b.Error("Error while saving a tag:", err)
	}
}

func setFrames(tag *Tag, encoding Encoding) {
	tag.SetTitle("Title")
	tag.SetArtist("Artist")
	tag.SetAlbum("Album")
	tag.SetYear("2016")
	tag.SetGenre("Genre")

	pic := PictureFrame{
		Encoding:    encoding,
		MimeType:    "image/jpeg",
		PictureType: PTFrontCover,
		Description: "Front cover",
		Picture:     frontCoverPicture,
	}
	tag.AddAttachedPicture(pic)

	uslt := UnsynchronisedLyricsFrame{
		Encoding:          encoding,
		Language:          "eng",
		ContentDescriptor: "Content descriptor",
		Lyrics:            "bogem/id3v2",
	}
	tag.AddUnsynchronisedLyricsFrame(uslt)

	comm := CommentFrame{
		Encoding:    encoding,
		Language:    "eng",
		Description: "Short description",
		Text:        "The actual text",
	}
	tag.AddCommentFrame(comm)
}
