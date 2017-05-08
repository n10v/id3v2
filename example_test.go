package id3v2_test

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/bogem/id3v2"
)

func Example() {
	// Open file and parse tag in it.
	tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
	if err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	defer tag.Close()

	// Read frames.
	fmt.Println(tag.Artist())
	fmt.Println(tag.Title())

	// Set simple text frames.
	tag.SetArtist("Artist")
	tag.SetTitle("Title")

	// Set comment frame.
	comment := id3v2.CommentFrame{
		Encoding:    id3v2.ENUTF8,
		Language:    "eng",
		Description: "My opinion",
		Text:        "Very good song",
	}
	tag.AddCommentFrame(comment)

	// Write tag to file.
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}

func ExampleTag_AddAttachedPicture() {
	tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
	if tag == nil || err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}

	artwork, err := ioutil.ReadFile("artwork.jpg")
	if err != nil {
		log.Fatal("Error while reading artwork file", err)
	}

	pic := id3v2.PictureFrame{
		Encoding:    id3v2.ENUTF8,
		MimeType:    "image/jpeg",
		PictureType: id3v2.PTFrontCover,
		Description: "Front cover",
		Picture:     artwork,
	}
	tag.AddAttachedPicture(pic)
}

func ExampleTag_AddCommentFrame() {
	tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
	if tag == nil || err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}

	comment := id3v2.CommentFrame{
		Encoding:    id3v2.ENUTF8,
		Language:    "eng",
		Description: "My opinion",
		Text:        "Very good song",
	}
	tag.AddCommentFrame(comment)
}

func ExampleTag_AddUnsynchronisedLyricsFrame() {
	tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
	if tag == nil || err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}

	uslt := id3v2.UnsynchronisedLyricsFrame{
		Encoding:          id3v2.ENUTF8,
		Language:          "ger",
		ContentDescriptor: "Deutsche Nationalhymne",
		Lyrics:            "Einigkeit und Recht und Freiheit...",
	}
	tag.AddUnsynchronisedLyricsFrame(uslt)
}

func ExampleTag_GetFrames() {
	tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
	if tag == nil || err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}

	pictures := tag.GetFrames(tag.CommonID("Attached picture"))
	if pictures != nil {
		for _, f := range pictures {
			pic, ok := f.(id3v2.PictureFrame)
			if !ok {
				log.Fatal("Couldn't assert picture frame")
			}
			// Do something with picture frame.
			// For example, print description of picture frame:
			fmt.Println(pic.Description)
		}
	}
}

func ExampleTag_GetLastFrame() {
	tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
	if tag == nil || err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}

	bpmFramer := tag.GetLastFrame(tag.CommonID("BPM"))
	if bpmFramer != nil {
		bpm, ok := bpmFramer.(id3v2.TextFrame)
		if !ok {
			log.Fatal("Couldn't assert bpm frame")
		}
		fmt.Println(bpm.Text)
	}
}
