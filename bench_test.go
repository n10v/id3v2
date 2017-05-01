// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"io/ioutil"
	"testing"
)

func BenchmarkParseAllFrames(b *testing.B) {
	if err := resetMP3Tag(); err != nil {
		b.Fatal("Error while reseting mp3 file:", err)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tag, err := Open(mp3Name, defaultOpts)
		if tag == nil || err != nil {
			b.Fatal("Error while opening mp3 file:", err)
		}
		if err = tag.Close(); err != nil {
			b.Error("Error while closing a tag:", err)
		}
	}
}

func BenchmarkParseAndWriteCommonCase(b *testing.B) {
	frontCover, err := ioutil.ReadFile(frontCoverName)
	if err != nil {
		b.Error("Error while reading front cover file")
	}

	// We use b.N+1, because in first iteration we just resetting tag
	// and setting common frames. Also timer will be resetted.
	for n := 0; n < b.N+1; n++ {
		tag, err := Open(mp3Name, defaultOpts)
		if tag == nil || err != nil {
			b.Fatal("Error while opening mp3 file:", err)
		}

		// Delete all frames in first iteration
		if n == 0 {
			tag.DeleteAllFrames()
		}

		tag.SetTitle("Title")
		tag.SetArtist("Artist")
		tag.SetYear("2016")

		// Set front cover
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

		// Reset timer because we just resetting tag in first iteration
		if n == 0 {
			b.ResetTimer()
		}
	}
}

func BenchmarkParseAndWriteManyFrames(b *testing.B) {
	frontCover, err := ioutil.ReadFile(frontCoverName)
	if err != nil {
		b.Error("Error while reading front cover file")
	}

	// We use b.N+1, because in first iteration we just resetting tag
	// and setting many frames. Also timer will be resetted.
	for n := 0; n < b.N+1; n++ {
		tag, err := Open(mp3Name, defaultOpts)
		if tag == nil || err != nil {
			b.Fatal("Error while opening mp3 file:", err)
		}

		// Delete all frames in first iteration
		if n == 0 {
			tag.DeleteAllFrames()
		}

		tag.SetTitle("Title")
		tag.SetArtist("Artist")
		tag.SetAlbum("Album")
		tag.SetYear("2016")
		tag.SetGenre("Genre")

		// Set front cover
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

		// Reset timer because we just resetting tag in first iteration
		if n == 0 {
			b.ResetTimer()
		}
	}
}
