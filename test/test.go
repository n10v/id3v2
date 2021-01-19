package main

import (
	"fmt"
	"github.com/bogem/id3v2"
	"log"
	"time"
)

func main() {
	tag, err := id3v2.Open("/Users/r_takaishi/Projects/kimagurefm/kimagurefm-047/kimagurefm-047.mp3", id3v2.Options{Parse: true})
	if err != nil {
		log.Fatal(err)
	}
	defer tag.Close()

	for _, frame := range tag.GetFrames("CHAP") {
		cf := frame.(id3v2.ChapterFrame)
		st := cf.StartTime
		st = st.Round(time.Second)
		h := st / time.Hour
		st -= h * time.Hour
		m := st / time.Minute
		st -= m * time.Minute
		s := st / time.Second

		fmt.Printf("%02d:%02d:%02d %s\n", h, m, s, cf.Title.Text)

	}
}
