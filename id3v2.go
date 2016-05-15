package id3v2

import (
	"bufio"
	"bytes"
	"github.com/bogem/id3v2/frame"
	"io/ioutil"
	"os"
)

func Open(name string) (*Tag, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return ParseTag(file)
}

func NewAttachedPicture() *frame.PictureFrame {
	return frame.NewPictureFrame()
}

func (t *Tag) Close() error {
	t.Header.FramesSize = t.CalcualteSizeOfAllFrames()

	f := t.File
	defer f.Close()

	if _, err := f.Seek(0, os.SEEK_SET); err != nil {
		return err
	}

	originalFile := bufio.NewReader(f)
	if _, err := originalFile.Discard(int(t.OriginalSize)); err != nil {
		return err
	}

	newFile, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer os.Remove(newFile.Name())
	newFileBuf := bufio.NewWriter(newFile)

	var newTag bytes.Buffer
	if header, err := FormTagHeader(*t.Header); err != nil {
		return err
	} else {
		newTag.Write(header)
	}
	if frames, err := t.FormFrames(); err != nil {
		return err
	} else {
		newTag.Write(frames)
	}

	if _, err = newFileBuf.ReadFrom(&newTag); err != nil {
		return err
	}
	if _, err = newFileBuf.ReadFrom(originalFile); err != nil {
		return err
	}
	if err = newFileBuf.Flush(); err != nil {
		return err
	}

	os.Rename(newFile.Name(), f.Name())
	return nil
}
