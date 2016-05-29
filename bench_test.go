package id3v2

import "testing"

func BenchmarkSet(b *testing.B) {
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

		pic := NewAttachedPicture()
		pic.SetMimeType("image/jpeg")
		pic.SetDescription("Cover")
		pic.SetPictureType(PTFrontCover)
		if err = pic.SetPictureFromFile("front_cover.jpg"); err != nil {
			b.Error("Error while setting a picture from file")
		}
		t.SetAttachedPicture(pic)

		if err = t.Flush(); err != nil {
			b.Error("Error while closing a tag: ", err)
		}
	}
}
