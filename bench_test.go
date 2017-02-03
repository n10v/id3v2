// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"io/ioutil"
	"testing"
)

func BenchmarkParseAllFrames(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tag, err := Open(mp3Name)
		if tag == nil || err != nil {
			b.Fatal("Error while opening mp3 file:", err)
		}

		tag.AllFrames()

		if err = tag.Close(); err != nil {
			b.Error("Error while closing a tag:", err)
		}
	}
}

func BenchmarkSetCommonCase(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tag, err := Open(mp3Name)
		if tag == nil || err != nil {
			b.Fatal("Error while opening mp3 file:", err)
		}

		tag.DeleteAllFrames()

		tag.SetTitle("Title")
		tag.SetArtist("Artist")
		tag.SetYear("2016")

		// Set front cover
		frontCover, err := ioutil.ReadFile(frontCoverName)
		if err != nil {
			b.Error("Error while reading front cover file")
		}

		pic := PictureFrame{
			Encoding:    ENUTF8,
			MimeType:    "image/jpeg",
			PictureType: PTFrontCover,
			Description: "Front cover",
			Picture:     frontCover,
		}
		tag.AddAttachedPicture(pic)

		if err = tag.Save(); err != nil {
			b.Error("Error while saving a tag:", err)
		}
		if err = tag.Close(); err != nil {
			b.Error("Error while closing a tag:", err)
		}
	}
}

func BenchmarkSetManyFrames(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tag, err := Open(mp3Name)
		if tag == nil || err != nil {
			b.Fatal("Error while opening mp3 file:", err)
		}

		tag.DeleteAllFrames()

		tag.SetTitle("Title")
		tag.SetArtist("Artist")
		tag.SetAlbum("Album")
		tag.SetYear("2016")
		tag.SetGenre("Genre")

		// Set front cover
		frontCover, err := ioutil.ReadFile(frontCoverName)
		if err != nil {
			b.Error("Error while reading front cover file")
		}

		pic := PictureFrame{
			Encoding:    ENUTF8,
			MimeType:    "image/jpeg",
			PictureType: PTFrontCover,
			Description: "Front cover",
			Picture:     frontCover,
		}
		tag.AddAttachedPicture(pic)

		// Set USLT
		uslt := UnsynchronisedLyricsFrame{
			Encoding:          ENUTF8,
			Language:          "eng",
			ContentDescriptor: "Content descriptor",
			Lyrics:            "bogem/id3v2",
		}
		tag.AddUnsynchronisedLyricsFrame(uslt)

		// Set comment
		comm := CommentFrame{
			Encoding:    ENUTF8,
			Language:    "eng",
			Description: "Short description",
			Text:        "The actual text",
		}
		tag.AddCommentFrame(comm)

		if err = tag.Save(); err != nil {
			b.Error("Error while saving a tag:", err)
		}
		if err = tag.Close(); err != nil {
			b.Error("Error while closing a tag:", err)
		}
	}
}
