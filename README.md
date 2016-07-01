# id3v2

## IMHO
I think, **ID3** is a very overwhelmed standard: it does **more than it really should do**. There are a lot of aspects, which developer should take into consideration. And that's why it's pretty complicated to write a **good library**. So if you have some thoughts about writing a **new simply and elegant standard** for providing information for digital music tracks, just write me. I think, it's a good time to write an appropriate standard for it 😉

## Information
**Fast and stable ID3 writing library for Go**

This library can only set and write tags, but can't read them. So if you only want to set tags, it fits you. And if there is a tag of version 3 or 4, this library will just delete this tag, because it can't parse tags yet. If version of the tag is smaller than 3, this library will return an error.

What it **can** do:
* Set artist, album, year, genre, unsynchronised lyrics/text (USLT), comments and **attached pictures** (e.g. album covers) and write all to file
* Set several USLT, comments and attached pictures
* Work with all allowed encodings

What it **can't** do:
* Parse tags
* Work with extended header, flags, padding

**If you want some functionality, that library can't do, just write an issue. I will implement it as fast as I can**

**And of course, pull requests are welcome!**

## Benchmarks

All benchmarks run on **MacBook Air 13" (early 2013, 1,4GHz Intel Core i5, 4GB 1600MHz DDR3)**

#### Set title, artist, year and 50KB picture to 4,6 MP3:
```
BenchmarkSetCommonCase-4	     200	   9255995 ns/op	   38386 B/op	      32 allocs/op
```

#### Set title, artist, album, year, genre, unsynchronised lyrics, comment and 50KB picture to 4,6MB MP3:
```
BenchmarkSetManyTags-4  	     200	   9268306 ns/op	   41177 B/op	      46 allocs/op
```

## Installation
  	$ go get -u github.com/bogem/id3v2

## Example of Usage:
```go
package main

import (
  "github.com/bogem/id3v2"
  "log"
)

func main() {
  // Open file and find tag in it
  tag, err := id3v2.Open("file.mp3")
  if err != nil {
   log.Fatal("Error while opening mp3 file: ", err)
  }

  // Set tags
  tag.SetArtist("Artist")
  tag.SetTitle("Title")

  comment := id3v2.NewCommentFrame()
  comment.SetLanguage("eng")
  comment.SetDescription("Short description")
  comment.SetText("The actual text")
  tag.AddCommentFrame(comment)

  // Write it to file
  if err = tag.Flush(); err != nil {
    log.Fatal("Error while flushing a tag: ", err)
  }
}

```

## Documentation

You can find it here: https://godoc.org/github.com/bogem/id3v2

## TODO

- [ ] Parse tags
- [ ] Work with extended header, flags, padding ***(Does anyone really use it?)***
- [x] Documentation
- [x] Work with other encodings

## License
MIT
