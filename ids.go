// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// Actual IDs for ID3v2.4
var (
	V24IDs = map[string]string{
		// Identification frames
		"Content group description":          "TIT1",
		"Title/Songname/Content description": "TIT2",
		"Subtitle/Description refinement":    "TIT3",
		"Album/Movie/Show title":             "TALB",
		"Original album/movie/show title":    "TOAL",
		"Track number/Position in set":       "TRCK",
		"Part of a set":                      "TPOS",
		"Set subtitle":                       "TSST",
		"ISRC":                               "TSRC",

		// Involved persons frames
		"Lead artist/Lead performer/Soloist/Performing group": "TPE1",
		"Band/Orchestra/Accompaniment":                        "TPE2",
		"Conductor":                                           "TPE3",
		"Interpreted, remixed, or otherwise modified by": "TPE4",
		"Original artist/performer":                      "TOPE",
		"Lyricist/Text writer":                           "TEXT",
		"Original lyricist/text writer":                  "TOLY",
		"Composer":                                       "TCOM",
		"Musician credits list":                          "TMCL",
		"Involved people list":                           "TIPL",
		"Encoded by":                                     "TENC",

		// Derived and subjective properties frames
		"BPM":          "TBPM",
		"Length":       "TLEN",
		"Initial key":  "TKEY",
		"Language":     "TLAN",
		"Content type": "TCON",
		"File type":    "TFLT",
		"Media type":   "TMED",
		"Mood":         "TMOO",

		// Rights and license frames
		"Copyright message":            "TCOP",
		"Produced notice":              "TPRO",
		"Publisher":                    "TPUB",
		"File owner/licensee":          "TOWN",
		"Internet radio station name":  "TRSN",
		"Internet radio station owner": "TRSO",

		// Other text frames
		"Original filename":                                "TOFN",
		"Playlist delay":                                   "TDLY",
		"Encoding time":                                    "TDEN",
		"Original release time":                            "TDOR",
		"Recording time":                                   "TDRC",
		"Release time":                                     "TDRL",
		"Tagging time":                                     "TDTG",
		"Software/Hardware and settings used for encoding": "TSSE",
		"Album sort order":                                 "TSOA",
		"Performer sort order":                             "TSOP",
		"Title sort order":                                 "TSOT",

		"Attached picture": "APIC",
		"Comments":         "COMM",
		"Unsynchronised lyrics/text transcription": "USLT",
	}
)
