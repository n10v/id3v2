package main

import (
	"fmt"
	"log"

	"github.com/bogem/id3v2/v3"
)

func main() {
	tag, err := id3v2.Open("in.mp3", id3v2.Options{Parse: true})
	if err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	defer tag.Close()

	// Read tags
	fmt.Println(tag.Artist())
	fmt.Println(tag.Title())

	// Write tag to file.mp3
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}
