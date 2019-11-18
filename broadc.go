package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// details in './details.go'

var topic string
var downloadDir = "/tmp/workspace"

func getData(audycja string) {

	// Request the HTML page
	res, err := http.Get(audycja)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s\n", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var dataMediaJson map[string]interface{}

	// Find the mp3
	doc.Find(".play.pr-media-play").Each(func(i int, s *goquery.Selection) {
		dataMedia, ok := s.Attr("data-media")
		if !ok {
			fmt.Println("no data-media attribute. continue...")
			return
		}

		err := json.Unmarshal([]byte(dataMedia), &dataMediaJson)
		if err != nil {
			fmt.Println("can't load json from html attribute. continue...")
			return
		}

		desc, _ := url.PathUnescape(dataMediaJson["desc"].(string))

		switch {
		case strings.Contains(desc, "ekonomiczny"):
			s := strings.Split(desc, " ")[2]
			desc = strings.Replace(s, ".", "-", -1)
			topic = "info"
		case strings.Contains(desc, "Winien"):
			s := strings.Split(desc, " ")[4]
			desc = strings.Replace(s, ".", "-", -1)
			topic = "winien"
		}

		fmt.Printf("curl -sSL -o %s/%s-%s.mp3 http:%v\n", downloadDir, desc, topic, dataMediaJson["file"])
	})
}

func makeMainDirectory() {
	// by default, os.ModePerm = 0777
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		log.Fatalf("Unable to create the directory. %v\n", err)
	}
	fmt.Printf("Directory created: %s\n", downloadDir)
}

func main() {

	winienima := "https://www.polskieradio.pl/9,Trojka/6253,Winien-i-ma"
	infekonomiczny := "https://www.polskieradio.pl/9/712"

	// make temp directory
	makeMainDirectory()

	// get broadcasts
	audycje := []string{winienima, infekonomiczny}

	for _, a := range audycje {
		getData(a)
	}
}
