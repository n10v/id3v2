package id3v2

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func prepareTestFile() (*os.File, error) {
	src, err := os.Open("./testdata/test.mp3")
	if err != nil {
		return nil, err
	}
	defer src.Close()

	tmpFile, err := ioutil.TempFile("", "chapter_test")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(tmpFile, src)
	if err != nil {
		return nil, err
	}
	return tmpFile, nil
}

func TestAddChapterFrame(t *testing.T) {
	tmpFile, err := prepareTestFile()
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpFile.Name())

	tag, err := Open(tmpFile.Name(), Options{Parse: true})
	if tag == nil || err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}

	chap := ChapterFrame{
		ElementID:   "chap0",
		StartTime:   0,
		EndTime:     time.Duration(1000 * nanosInMillis),
		StartOffset: 0,
		EndOffset:   0,
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
	if frame.StartTime != 0 {
		t.Errorf("expected: %d, but got %s", 0, frame.StartTime)
	}
	if frame.EndTime.Seconds()*1000 != 1000 {
		t.Errorf("expected: %d, but got %d", 1000, int(frame.EndTime.Seconds()*1000))
	}
}

func TestAddChapterFrameWithTitle(t *testing.T) {
	tmpFile, err := prepareTestFile()
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpFile.Name())

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
		Title: &TextFrame{
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
	if frame.Title.Text != "chapter 0" {
		t.Errorf("expected: %s, but got %s", "chapter 0", frame.Title)
	}
}

func TestAddChapterFrameWithSubTitle(t *testing.T) {
	tmpFile, err := prepareTestFile()
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpFile.Name())

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
		SubTitle: &TextFrame{
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
	if frame.SubTitle.Text != "chapter 0" {
		t.Errorf("expected: %s, but got %s", "chapter 0", frame.Title)
	}
}
