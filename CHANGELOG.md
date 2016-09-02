# ID3v2 Changelog

## 0.8
* Major underhood changes: more perfomance and stability,
less bugs and memory consumption

## 0.7.1
* Returned Add(AttachedPicture/Comment/UnsynchronisedLyrics)Frame functions
for convenience and less errors.
See [GitHub issue](https://github.com/bogem/id3v2/issues/5) for details

## 0.7
* Now ID3v2 can parse tags! See documentation for details
* Some changes in API for convenice and fliexibility

## 0.6
* Huge update of API and documentation

## 0.5.1
* Noticeable improvement in performance
* Added panics in `util` package

## 0.5
* Added support of encodings
* Added documentation (https://godoc.org/github.com/bogem/id3v2)

## 0.4
* Added USLT and COMM frames (by request from https://github.com/bogem/id3v2/issues/3)
* Some memory improvements

## 0.3
* Noticeably decreased allocating memory space
* Some changes in setting pictures to picture frames
* Major bug fixes

## 0.2.2
* Some very small improvement in memory usage

## 0.2.1
* Added check of unsupported picture type while setting picture frame to tag

## 0.2
* Noticeably decreased allocating memory space

## 0.1.2
* Renamed tag.Close() to tag.Flush()

## 0.1.1
* Little improvement in performance

## 0.1
* First release!
