package id3v2

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"
)

const (
	testChapterTocSampleTitle = "Chapter TOC title"
)

func newChapterFrames(noOfChapters int) []ChapterFrame {
	var start time.Duration
	offset := time.Duration(1000 * nanosInMillis)

	chapters := make([]ChapterFrame, noOfChapters)

	for i := 0; i < noOfChapters; i++ {
		end := start + offset

		chapters[i] = ChapterFrame{
			ElementID:   fmt.Sprintf("ch%d", i),
			StartTime:   start,
			EndTime:     end,
			StartOffset: IgnoredOffset,
			EndOffset:   IgnoredOffset,
			Title: &TextFrame{
				Encoding: EncodingUTF8,
				Text:     fmt.Sprintf("Chapter %d", i),
			},
		}

		start = end
	}

	return chapters
}

func TestAddChapterTocFrame(t *testing.T) {
	const noOfChapters = 5
	buf := &bytes.Buffer{}
	tag := NewEmptyTag()

	chapters := newChapterFrames(noOfChapters)

	chapterIds := make([]string, len(chapters))
	for i, c := range chapters {
		tag.AddChapterFrame(c)

		chapterIds[i] = c.ElementID
	}

	chapterToc := ChapterTocFrame{
		ElementID:  "Main TOC",
		TopLevel:   true,
		Ordered:    true,
		ChapterIds: chapterIds,
		Description: &TextFrame{
			Encoding: EncodingUTF8,
			Text:     testChapterTocSampleTitle,
		},
	}

	tag.AddChapterTocFrame(chapterToc)
	tag.WriteTo(buf)

	// Read back

	tagBack, err := ParseReader(buf, Options{Parse: true})
	if err != nil {
		log.Fatal("Error parsing mp3 content: ", err)
	}

	if !tagBack.HasFrames() {
		log.Fatal("No tags in content in mp3 content")
	}

	chapterTocBackFrame := tag.GetLastFrame("CTOC")
	if chapterTocBackFrame == nil {
		log.Fatal("Error getting chapter TOC frame: ", err)
	}

	chapterTocBack, ok := chapterTocBackFrame.(ChapterTocFrame)
	if !ok {
		log.Fatal("Error casting chapter TOC frame")
	}

	if chapterToc.ElementID != chapterTocBack.ElementID {
		t.Errorf("Expected element ID: %s, but got %s", chapterToc.ElementID, chapterTocBack.ElementID)
	}

	if chapterToc.TopLevel != chapterTocBack.TopLevel {
		t.Errorf("Expected top level: %v, but got %v", chapterToc.TopLevel, chapterTocBack.TopLevel)
	}

	if chapterToc.Ordered != chapterTocBack.Ordered {
		t.Errorf("Expected ordered: %v, but got %v", chapterToc.Ordered, chapterTocBack.Ordered)
	}

	if expected, actual := len(chapterToc.ChapterIds), len(chapterTocBack.ChapterIds); expected != actual {
		t.Errorf("Expected ordered: %v, but got %v", expected, actual)
	}

	for i := 0; i < len(chapterToc.ChapterIds); i++ {
		if expected, actual := chapterToc.ChapterIds[i], chapterTocBack.ChapterIds[i]; expected != actual {
			t.Errorf("Expected chapter reference at index: %d: %s, but got %s", i, expected, actual)
		}
	}

	if chapterToc.Description != nil && chapterToc.Description.Text != chapterTocBack.Description.Text {
		t.Errorf("Expected description: %s, but got %s", chapterToc.Description.Text, chapterTocBack.Description.Text)
	}
}
