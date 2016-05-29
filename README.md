# ID3v2

## IMHO
I think, **ID3v2** is a very overwhelmed standard: it does **more than it really should do**. There are a lot of aspects, which developer should take into consideration. And that's why it's pretty complicated to write a **good library**. So if you have some thoughts about writing a **new simply and elegant standard** for providing information for digital music tracks, just write me. I think, it's a good time to write an appropriate standard for it 😉

## Information
**Stable and fast ID3v2 Library written in Go**

This library can only set and write tags, but can't read them. So if you only want to set tags, it fits you. And if there is a tag of version 3 or 4, this library will just delete it.

What it **can** do:
* Set artist, album, year, genre and **attached pictures** (e.g. album covers) and write all to file
* Set several attached pictures

What it **can't** do:
* Parse tags
* Work with extended header, flags, padding
* Work with encodings, except UTF-8

**If you want some functionality, that library can't do, just write an issue. I will implement it as fast as I can**

**And of course, pull requests are welcome!**

## Benchmarks

All benchmarks run on **MacBook Air 13" (early 2013, 1,4GHz Intel Core i5, 4GB 1600MHz DDR3)**

#### Set title, artist, album, year, genre and 50KB picture to 4,6MB MP3:
```
BenchmarkSet-4	     100	  15646308 ns/op	   13777 B/op	      43 allocs/op
```

## Installation
  	$ go get -u github.com/bogem/id3v2

## Usage
#### Example:
```go
package main

import (
  "github.com/bogem/id3v2"
  "log"
)

func main() {
  tag, err := id3v2.Open("file.mp3")
  if err != nil {
   log.Fatal("Error while opening mp3 file: ", err)
  }

  tag.SetArtist("Artist")
  tag.SetTitle("Title")
  tag.SetYear("2016")
  ...


  if err = tag.Flush(); err != nil {
    log.Fatal("Error while closing a tag: ", err)
  }
}

```

#### Available functions for setting text frames:
```go
tag.SetTitle(title string)
tag.SetArtist(artist string)
tag.SetAlbum(album string)
tag.SetYear(year string)
tag.SetGenre(genre string)
```

#### Setting a picture:
```go
package main

import (
  "github.com/bogem/id3v2"
  "log"
  "os"
)

func main() {
  tag, err := id3v2.Open("file.mp3")
  if err != nil {
    log.Fatal("Error while opening mp3 file: ", err)
  }

  pic := id3v2.NewAttachedPicture()
  pic.SetMimeType("image/jpeg")
  pic.SetDescription("Cover")
  pic.SetPictureType(id3v2.PTFrontCover)
  if err = pic.SetPictureFromFile("artwork.jpg"); err != nil {
    log.Fatal("Error while setting a picture from file: ", err)
  }

  if err = tag.SetAttachedPicture(pic); err != nil {
		log.Fatal("Error while setting a picture frame to tag: ", err)
  }

  if err = tag.Flush(); err != nil {
    log.Fatal("Error while closing a tag: ", err)
  }
}
```

**Available picture types:**
* `PTOther`
* `PTFileIcon`
* `PTOtherFileIcon`
* `PTFrontCover`
* `PTBackCover`
* `PTLeafletPage`
* `PTMedia`
* `PTLeadArtistSoloist`
* `PTArtistPerformer`
* `PTConductor`
* `PTBandOrchestra`
* `PTComposer`
* `PTLyricistTextWriter`
* `PTRecordingLocation`
* `PTDuringRecording`
* `PTDuringPerformance`
* `PTMovieScreenCapture`
* `PTBrightColouredFish`
* `PTIllustration`
* `PTBandArtistLogotype`
* `PTPublisherStudioLogotype`

## TODO

- [ ] Parse tags
- [ ] Work with other encodings
- [ ] Work with extended header, flags, padding ***(Does somebody really use it?)***

## License
MIT
