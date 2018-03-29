// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"testing"
)

var frontCoverPicture = mustReadFile(frontCoverPath)

func BenchmarkParseAllFrames(b *testing.B) {
	writeTag(b, EncodingUTF8)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tag, err := Open(mp3Path, parseOpts)
		if tag == nil || err != nil {
			b.Fatal("Error while opening mp3 file:", err)
		}
		if err = tag.Close(); err != nil {
			b.Error("Error while closing a tag:", err)
		}
	}
}

func BenchmarkParseAllFramesISO(b *testing.B) {
	writeTag(b, EncodingISO)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tag, err := Open(mp3Path, parseOpts)
		if tag == nil || err != nil {
			b.Fatal("Error while opening mp3 file:", err)
		}
		if err = tag.Close(); err != nil {
			b.Error("Error while closing a tag:", err)
		}
	}
}

func BenchmarkParseArtistAndTitle(b *testing.B) {
	writeTag(b, EncodingUTF8)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tag, err := Open(mp3Path, Options{Parse: true, ParseFrames: []string{"Artist", "Title"}})
		if tag == nil || err != nil {
			b.Fatal("Error while opening mp3 file:", err)
		}
		if err = tag.Close(); err != nil {
			b.Error("Error while closing a tag:", err)
		}
	}
}

func BenchmarkWrite(b *testing.B) {
	for n := 0; n < b.N; n++ {
		writeTag(b, EncodingUTF8)
	}
}

func BenchmarkWriteISO(b *testing.B) {
	for n := 0; n < b.N; n++ {
		writeTag(b, EncodingISO)
	}
}

func writeTag(b *testing.B, encoding Encoding) {
	tag, err := Open(mp3Path, Options{Parse: false})
	if tag == nil || err != nil {
		b.Fatal("Error while opening mp3 file:", err)
	}
	defer tag.Close()

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

	if err = tag.Save(); err != nil {
		b.Error("Error while saving a tag:", err)
	}
}
