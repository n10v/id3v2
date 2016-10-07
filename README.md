[![Release](https://img.shields.io/github/release/bogem/id3v2.svg?maxAge=2592000)](https://github.com/bogem/id3v2/releases)
[![GoDoc](https://godoc.org/github.com/bogem/id3v2?status.svg)](https://godoc.org/github.com/bogem/id3v2)
[![Build Status](https://travis-ci.org/bogem/id3v2.svg?branch=master)](https://travis-ci.org/bogem/id3v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/bogem/id3v2)](https://goreportcard.com/report/github.com/bogem/id3v2)

# id3v2

## IMHO
I think, **ID3** is a very overwhelmed standard: it does **more than it really should do**.
There are a lot of aspects, which developer should take into consideration.
And that's why it's pretty complicated to write a **good library**.
So if you have some thoughts about writing a **new simply and elegant standard**
for providing information for digital music tracks, just write me.
I think, it's a good time to write an appropriate standard for it 😉

## Information
**Fast and stable ID3 parsing and writing library for Go**

It can:
* Support ID3v2.3 and ID3v2.4 tags
* Parse and write tags
* Set artist, album, year, genre, unsynchronised lyrics/text (USLT),
comments and attached pictures
* Set several USLTs, comments and attached pictures
* Work with all available encodings

**If you want some functionality, that library can't do,
or you have some questions, just write an issue.**

**And of course, pull requests are welcome!**

## Installation
  	$ go get -u github.com/bogem/id3v2

## Example of Usage:
```go
package main

import (
	"fmt"
	"log"

	"github.com/bogem/id3v2"
)

func main() {
	// Open file and find tag in it
	tag, err := id3v2.Open("file.mp3")
	if err != nil {
 		log.Fatal("Error while opening mp3 file: ", err)
 	}
	defer tag.Close()

	// Read tags
	fmt.Println(tag.Artist())
	fmt.Println(tag.Title())

  	// Set tags
	tag.SetArtist("New artist")
	tag.SetTitle("New title")

	comment := id3v2.CommentFrame{
		Encoding:   id3v2.ENUTF8,
		Language:   "eng",
		Desciption: "My opinion",
		Text:       "Very good song",
	}
	tag.AddCommentFrame(comment)

	// Write it to file
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}

```

## Documentation

You can find it here: https://godoc.org/github.com/bogem/id3v2

## TODO

- [ ] Unsynchronization (?)
- [ ] Work with extended header, flags, padding, footer ***(Does anyone really use it?)***
- [x] ID3v2.3 Support
- [x] Parse tags
- [x] Documentation
- [x] Work with other encodings

## License
MIT
