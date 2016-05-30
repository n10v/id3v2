package id3v2

import "testing"

func BenchmarkSetCommonCase(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t, err := Open("test.mp3")
		if t == nil || err != nil {
			b.Error("Error while opening mp3 file: ", err)
		}
		t.SetTitle("Title")
		t.SetArtist("Artist")
		t.SetYear("2016")

		// Setting front cover
		pic := NewAttachedPicture()
		pic.SetMimeType("image/jpeg")
		pic.SetPictureType(PTFrontCover)
		if err = pic.SetPictureFromFile("front_cover.jpg"); err != nil {
			b.Error("Error while setting a picture from file")
		}
		t.AddAttachedPicture(pic)

		if err = t.Flush(); err != nil {
			b.Error("Error while closing a tag: ", err)
		}
	}
}

func BenchmarkSetManyTags(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t, err := Open("test.mp3")
		if t == nil || err != nil {
			b.Error("Error while opening mp3 file: ", err)
		}
		t.SetTitle("Title")
		t.SetArtist("Artist")
		t.SetAlbum("Album")
		t.SetYear("2016")
		t.SetGenre("Genre")

		// Setting front cover
		pic := NewAttachedPicture()
		pic.SetMimeType("image/jpeg")
		pic.SetDescription("Cover")
		pic.SetPictureType(PTFrontCover)
		if err = pic.SetPictureFromFile("front_cover.jpg"); err != nil {
			b.Error("Error while setting a picture from file")
		}
		t.AddAttachedPicture(pic)

		// Setting USLT
		uslt := NewUnsynchronisedLyricsFrame()
		uslt.SetLanguage("eng")
		uslt.SetContentDescriptor("Content descriptor")
		uslt.SetLyrics("bogem/id3v2")
		t.AddUnsynchronisedLyricsFrame(uslt)

		// Setting comment
		cf := NewCommentFrame()
		cf.SetLanguage("eng")
		cf.SetDescription("Short description")
		cf.SetText("The actual text")
		t.AddCommentFrame(cf)

		if err = t.Flush(); err != nil {
			b.Error("Error while closing a tag: ", err)
		}
	}
}
