# id3v2 [![GoDoc](https://godoc.org/github.com/bogem/id3v2?status.svg)](https://godoc.org/github.com/bogem/id3v2) [![Build Status](https://travis-ci.org/bogem/id3v2.svg?branch=master)](https://travis-ci.org/bogem/id3v2)

**Fast, simple and powerful ID3 parsing and writing library for Go, based only on standard library.**

id3v2 can:

* ✅ support of ID3v2.3 and ID3v2.4 tags;
* ✅ parse and write tags;
* ✅ work with all available encodings;
* ✅ set and read all text frames, unsynchronised lyrics/text (USLT),
  comments, attached pictures, UFID and TXXX frames;
* ✅ set and read frames, that can be used multiple times in tag;
* ✅ be used in multiple goroutines.

id3v2 can't:

* ❌ work with unsynchronization, extended header, flags, padding, footer.

If you want some functionality, that library can't do,
or you have some questions, just write an issue. **And of course, pull requests are welcome!**

## Installation

```
go get -u github.com/bogem/id3v2
```

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

## Work with encodings

id3v2 can encode and decode text of avaialble encodings (ISO-8859-1,
UTF-16 with BOM, UTF-16BE without BOM, UTF-8). Just set [the desired encoding](https://github.com/bogem/id3v2/blob/master/encoding.go#L27) in `Encoding` field of frame and set UTF-8 encoded text to corresponding text fields.

For example, if you set comment frame with custom encoding and write it:

```go
tag := id3v2.NewEmptyTag()
comment := id3v2.CommentFrame{
	Encoding:    id3v2.EncodingUTF16,
	Language:    "ger",
	Description: "Tier",
	Text:        "Der Löwe",
}
tag.AddCommentFrame(comment)

_, err := tag.WriteTo(w)
if err != nil {
	log.Fatal(err)
}
```

`Text` field will be automatically encoded with UTF-16BE with BOM and written to w.

By default, if version of tag is 4 then UTF-8 is used for methods like
`SetArtist`, `SetTitle`, `SetGenre` and etc, otherwise ISO-8859-1.

## Benchmarks

See benchmarks in the description of [last release](https://github.com/bogem/id3v2/releases).

## Reference

https://godoc.org/github.com/bogem/id3v2
