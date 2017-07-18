[![GoDoc](https://godoc.org/github.com/bogem/id3v2?status.svg)](https://godoc.org/github.com/bogem/id3v2)
[![Build Status](https://travis-ci.org/bogem/id3v2.svg?branch=master)](https://travis-ci.org/bogem/id3v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/bogem/id3v2)](https://goreportcard.com/report/github.com/bogem/id3v2)

# id3v2

## IMHO
I think, **ID3** is a very overwhelmed standard: it does **more than it really should do**.
There are a lot of aspects, which developer should take into consideration.
And that's why it's pretty complicated to write a **good library**.
So if you have some thoughts and ideas or if you want to support me about writing a **new simple and elegant standard**
for providing information for digital music tracks, just write me an e-mail.
I think, it's a good time to write an appropriate standard for it 😉

## Description
**Fast and stable ID3 parsing and writing library for Go, based only on standard library and without any third-party dependency.**

It can:
* support ID3v2.3 and ID3v2.4 tags;
* parse and write tags;
* set artist, album, year, genre, unsynchronised lyrics/text (USLT),
comments and attached pictures;
* set several USLTs, comments and attached pictures;
* be used in multiple goroutines.

It can't:
* do unsyncronization;
* work with extended header, flags, padding, footer.

**id3v2 is still in beta. Until version 1.0 the API may be changed.**

If you have **issues with encoding** or you don't know, how to **set
encoding to frame**, please write new issue!

If you want some **functionality, that library can't do**,
or you have some **questions**, just write an issue.

**And of course, pull requests are welcome!**

## Installation
  	$ go get -u github.com/bogem/id3v2

## Example of usage
```go
package main

import (
	"fmt"
	"log"

	"github.com/bogem/id3v2"
)

func main() {
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
	tag.SetArtist("New artist")
	tag.SetTitle("New title")

	// Set comment frame.
	comment := id3v2.CommentFrame{
		Encoding:    id3v2.EncodingUTF8,
		Language:    "eng",
		Description: "My opinion",
		Text:        "Very good song",
	}
	tag.AddCommentFrame(comment)

	// Write it to file.
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}
```

## Read multiple frames
```go
pictures := tag.GetFrames(tag.CommonID("Attached picture"))
for _, f := range pictures {
	pic, ok := f.(id3v2.PictureFrame)
	if !ok {
		log.Fatal("Couldn't assert picture frame")
	}

	// Do something with picture frame.
	// For example, print the description:
	fmt.Println(pic.Description)
}
```

## Options
```go
// Options influence on processing the tag.
type Options struct {
	// Parse defines, if tag will be parsed.
	Parse bool

	// ParseFrames defines, that frames do you only want to parse. For example,
	// `ParseFrames: []string{"Artist", "Title"}` will only parse artist
	// and title frames. You can specify IDs ("TPE1", "TIT2") as well as
	// descriptions ("Artist", "Title"). If ParseFrame is blank or nil,
	// id3v2 will parse all frames in tag. It works only if Parse is true.
	//
	// It's very useful for performance, so for example
	// if you want to get only some text frames,
	// id3v2 will not parse huge picture or unknown frames.
	ParseFrames []string
}
```

## Documentation

https://godoc.org/github.com/bogem/id3v2
