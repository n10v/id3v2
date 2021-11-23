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
	type fields struct {
		ElementID   string
		StartTime   time.Duration
		EndTime     time.Duration
		StartOffset uint32
		EndOffset   uint32
		Title       *TextFrame
		Description *TextFrame
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "element id only",
			fields: fields{
				ElementID:   "chap0",
				StartTime:   0,
				EndTime:     time.Duration(1000 * nanosInMillis),
				StartOffset: 0,
				EndOffset:   0,
			},
		},
		{
			name: "with title",
			fields: fields{
				ElementID:   "chap0",
				StartTime:   0,
				EndTime:     time.Duration(1000 * nanosInMillis),
				StartOffset: 0,
				EndOffset:   0,
				Title: &TextFrame{
					Encoding: EncodingUTF8,
					Text:     "chapter 0",
				},
			},
		},
		{
			name: "with description",
			fields: fields{
				ElementID:   "chap0",
				StartTime:   0,
				EndTime:     time.Duration(1000 * nanosInMillis),
				StartOffset: 0,
				EndOffset:   0,
				Description: &TextFrame{
					Encoding: EncodingUTF8,
					Text:     "chapter 0",
				},
			},
		},
		{
			name: "with title and description",
			fields: fields{
				ElementID:   "chap0",
				StartTime:   0,
				EndTime:     time.Duration(1000 * nanosInMillis),
				StartOffset: 0,
				EndOffset:   0,
				Title: &TextFrame{
					Encoding: EncodingUTF8,
					Text:     "chapter 0 title",
				},
				Description: &TextFrame{
					Encoding: EncodingUTF8,
					Text:     "chapter 0 description",
				},
			},
		},
		{
			name: "non-zero time and offset",
			fields: fields{
				ElementID:   "chap0",
				StartTime:   time.Duration(1000 * nanosInMillis),
				EndTime:     time.Duration(1000 * nanosInMillis),
				StartOffset: 10,
				EndOffset:   10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := prepareTestFile()
			if err != nil {
				t.Error(err)
			}
			defer os.Remove(tmpFile.Name())

			tag, err := Open(tmpFile.Name(), Options{Parse: true})
			if tag == nil || err != nil {
				log.Fatal("Error while opening mp3 file: ", err)
			}

			cf := ChapterFrame{
				ElementID:   tt.fields.ElementID,
				StartTime:   tt.fields.StartTime,
				EndTime:     tt.fields.EndTime,
				StartOffset: tt.fields.StartOffset,
				EndOffset:   tt.fields.EndOffset,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
			}
			tag.AddChapterFrame(cf)

			if err := tag.Save(); err != nil {
				t.Error(err)
			}
			tag.Close()

			tag, err = Open(tmpFile.Name(), Options{Parse: true})
			if tag == nil || err != nil {
				log.Fatal("Error while opening mp3 file: ", err)
			}
			frame := tag.GetLastFrame("CHAP").(ChapterFrame)
			if frame.ElementID != tt.fields.ElementID {
				t.Errorf("Expected element ID: %s, but got %s", tt.fields.ElementID, frame.ElementID)
			}
			if tt.fields.Title != nil && frame.Title.Text != tt.fields.Title.Text {
				t.Errorf("Expected title: %s, but got %s", tt.fields.Title.Text, frame.Title)
			}
			if tt.fields.Description != nil && frame.Description.Text != tt.fields.Description.Text {
				t.Errorf("Expected description: %s, but got %s", tt.fields.Description.Text, frame.Description.Text)
			}
			if frame.StartTime != tt.fields.StartTime {
				t.Errorf("Expected start time: %s, but got %s", tt.fields.StartTime, frame.StartTime)
			}
			if frame.EndTime != tt.fields.EndTime {
				t.Errorf("Expected end time: %s, but got %s", tt.fields.EndTime, frame.EndTime)
			}
			if frame.StartOffset != tt.fields.StartOffset {
				t.Errorf("Expected start offset: %d, but got %d", tt.fields.StartOffset, frame.StartOffset)
			}
			if frame.EndOffset != tt.fields.EndOffset {
				t.Errorf("Expected end offset: %d, but got %d", tt.fields.EndOffset, frame.EndOffset)
			}
		})
	}
}
