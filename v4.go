// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// Actual IDs for ID3v2.4
var (
	V24CommonIDs = map[string]string{
		"Title":                      "TIT2",
		"Artist":                     "TPE1",
		"Album":                      "TALB",
		"Year":                       "TYER", // Deprecated in ID3v2.4
		"Genre":                      "TCON",
		"Attached picture":           "APIC",
		"Unsynchronised lyrics/text": "USLT",
		"Comment":                    "COMM",
	}
)
