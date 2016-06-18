package frame

// Actual IDs for ID3v2.4
var (
	V24CommonIDs = map[string]string{
		"Title":            "TIT2",
		"Artist":           "TPE1",
		"Album":            "TALB",
		"Year":             "TYER", // Deprecated in ID3v2.4
		"Genre":            "TCON",
		"Attached picture": "APIC",
		"USLT":             "USLT",
		"Comment":          "COMM",
	}
)
