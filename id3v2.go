package id3v2

import (
	"bytes"
	"github.com/bogem/id3v2/frame"
	"github.com/bogem/id3v2/util"
	"os"
	"sync"
)

const (
	// Picture types
	PTOther                   = 0
	PTFileIcon                = 1
	PTOtherFileIcon           = 2
	PTFrontCover              = 3
	PTBackCover               = 4
	PTLeafletPage             = 5
	PTMedia                   = 6
	PTLeadArtistSoloist       = 7
	PTArtistPerformer         = 8
	PTConductor               = 9
	PTBandOrchestra           = 10
	PTComposer                = 11
	PTLyricistTextWriter      = 12
	PTRecordingLocation       = 13
	PTDuringRecording         = 14
	PTDuringPerformance       = 15
	PTMovieScreenCapture      = 16
	PTBrightColouredFish      = 17
	PTIllustration            = 18
	PTBandArtistLogotype      = 19
	PTPublisherStudioLogotype = 20
)

var (
	// Encodings
	ENISO = util.Encoding{
		Key:              0,
		TerminationBytes: []byte{0},
	}
	ENUTF16 = util.Encoding{
		Key:              1,
		TerminationBytes: []byte{0, 0},
	}
	ENUTF16BE = util.Encoding{
		Key:              2,
		TerminationBytes: []byte{0, 0},
	}
	ENUTF8 = util.Encoding{
		Key:              3,
		TerminationBytes: []byte{0},
	}
)

var bytesBufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func Open(name string) (*Tag, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return ParseTag(file)
}

func NewAttachedPicture() *frame.PictureFrame {
	pf := new(frame.PictureFrame)
	pf.SetEncoding(ENUTF8)
	return pf
}

func NewCommentFrame() *frame.CommentFrame {
	cf := new(frame.CommentFrame)
	cf.SetEncoding(ENUTF8)
	return cf
}

func NewTextFrame(text string) *frame.TextFrame {
	tf := new(frame.TextFrame)
	tf.SetEncoding(ENUTF8)
	tf.SetText(text)
	return tf
}

func NewUnsynchronisedLyricsFrame() *frame.UnsynchronisedLyricsFrame {
	uslf := new(frame.UnsynchronisedLyricsFrame)
	uslf.SetEncoding(ENUTF8)
	return uslf
}
