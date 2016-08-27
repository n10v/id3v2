// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import "os"
import "testing"

func BenchmarkSetCommonCase(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tag, err := Open(mp3Name)
		if tag == nil || err != nil {
			b.Error("Error while opening mp3 file: ", err)
		}
		defer tag.Close()
		tag.SetTitle("Title")
		tag.SetArtist("Artist")
		tag.SetYear("2016")

		// Setting front cover
		frontCover, err := os.Open(frontCoverName)
		if err != nil {
			b.Error("Error while opening front cover file")
		}
		defer frontCover.Close()
		pic := PictureFrame{
			Encoding:    ENUTF8,
			MimeType:    "image/jpeg",
			PictureType: PTFrontCover,
			Description: "Front cover",
			Picture:     frontCover,
		}
		tag.AddAttachedPicture(pic)

		if err = tag.Save(); err != nil {
			b.Error("Error while closing a tag: ", err)
		}
	}
}

func BenchmarkSetManyFrames(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tag, err := Open(mp3Name)
		if tag == nil || err != nil {
			b.Error("Error while opening mp3 file: ", err)
		}
		defer tag.Close()
		tag.SetTitle("Title")
		tag.SetArtist("Artist")
		tag.SetAlbum("Album")
		tag.SetYear("2016")
		tag.SetGenre("Genre")

		// Setting front cover
		frontCover, err := os.Open(frontCoverName)
		if err != nil {
			b.Error("Error while opening front cover file")
		}
		defer frontCover.Close()

		pic := PictureFrame{
			Encoding:    ENUTF8,
			MimeType:    "image/jpeg",
			PictureType: PTFrontCover,
			Description: "Front cover",
			Picture:     frontCover,
		}
		tag.AddAttachedPicture(pic)

		// Setting USLT
		uslt := UnsynchronisedLyricsFrame{
			Encoding:          ENUTF8,
			Language:          "eng",
			ContentDescriptor: "Content descriptor",
			Lyrics:            "bogem/id3v2",
		}
		tag.AddUnsynchronisedLyricsFrame(uslt)

		// Setting comment
		comm := CommentFrame{
			Encoding:    ENUTF8,
			Language:    "eng",
			Description: "Short description",
			Text:        "The actual text",
		}
		tag.AddCommentFrame(comm)

		if err = tag.Save(); err != nil {
			b.Error("Error while closing a tag: ", err)
		}
	}
}
