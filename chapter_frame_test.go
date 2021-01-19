package id3v2

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestAddChapterFrame(t *testing.T) {
	src, err := os.Open("./testdata/test.mp3")
	if err != nil {
		t.Error(err)
	}
	defer src.Close()

	tmpFile, err := ioutil.TempFile("", "chapter_test")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, src)
	if err != nil {
		t.Error(err)
	}

	tag, err := Open(tmpFile.Name(), Options{Parse: true})
	if tag == nil || err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}

	chap := ChapterFrame{
		ElementID:   "chap0",
		StartTime:   0,
		EndTime:     0,
		StartOffset: 0,
		EndOffset:   0,
		TIT2SubFrame: &TextFrame{
			Encoding: EncodingUTF8,
			Text:     "chapter 0",
		},
	}
	tag.AddChapterFrame(chap)

	if err := tag.Save(); err != nil {
		t.Error(err)
	}
	tag.Close()

	tag, err = Open(tmpFile.Name(), Options{Parse: true})
	if tag == nil || err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	frame := tag.GetLastFrame("CHAP").(ChapterFrame)
	if frame.ElementID != "chap0" {
		t.Error(err)
	}
	if frame.Title != "chapter 0" {
		t.Errorf("expected: %s, but got %s", "chapter 0", frame.Title)
	}
}
