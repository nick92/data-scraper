package main

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	// "encoding/csv"
	// "encoding/xml"
	"fmt"
	// "golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	settingsConfig = "settings.json"
	scrapingJSON   = "scraping.json"
	outputJSON     = "output.json"
)

type Selectors struct {
	Id              string
	Type            string
	ParentSelectors []string
	Selector        string
	Multiple        bool
	Regex           string
	Delay           int
}

type Scraping struct {
	StartUrl  []string
	Id        string
	Selectors []Selectors
}

func readSettingsJSON() {
	// open the file and read the file
	data, err := ioutil.ReadFile(settingsConfig)
	// define data struture
	type Config struct {
		JavaScript    bool
		Proxy         bool
		ProxyLists    []string
		RotatingProxy bool
		Export        string
	}
	// json data
	var settings Config
	err = json.Unmarshal(data, &settings)
	// log any errors
	if err != nil {
		log.Println(err)
	}
	// just priting it
	fmt.Println("JavaScript: ", settings.JavaScript)
	fmt.Println("Proxy: ", settings.Proxy)
	fmt.Println("ProxyLists: ", settings.ProxyLists)
	fmt.Println("RotatingProxy: ", settings.RotatingProxy)
	fmt.Println("Export: ", settings.Export)
}

func readScrapingJSON() *Scraping {
	// open the file and read the file
	data, err := ioutil.ReadFile(scrapingJSON)

	var scrape Scraping
	err = json.Unmarshal(data, &scrape)
	// fmt.Printf("id = %v\n", scrape.Id)
	// fmt.Printf("url = %v\n", scrape.StartUrl)
	// fmt.Printf("PS = %v\n\n", scrape.Selectors)
	// log any errors
	if err != nil {
		log.Println(err)
	}

	return &scrape
	// lets just print for now
	// fmt.Println("id", scrape.startUrl)
	// run the scraper and start scraping.
	// scraper()
}

func scraper() {
	scrapingJson := readScrapingJSON()
	fmt.Printf("url = %v\n", scrapingJson.StartUrl)
	response, err := http.Get(scrapingJson.StartUrl[0])
	if err != nil {
		log.Println(err)
	}

	defer response.Body.Close()

	// Get the response body as a string
	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)

	for _, selector := range scrapingJson.Selectors {
		if selector.Type == "SelectorText" {
			SelectorText(&pageContent, &selector)
		} else if selector.Type == "SelectorImage" {
			SelectorImage(&pageContent, &selector)
		}

	}
}

func SelectorText(pageContent *string, selector *Selectors) {
	titleStartIndex := 1
	// Find a substr
	// fmt.Println(pageContent)
	titleStartIndex = strings.Index(*pageContent, "<"+selector.Selector+">")
	if titleStartIndex == -1 {
		fmt.Printf("%s : No element found", selector.Id)
		os.Exit(0)
	}
	// The start index of the title is the index of the first
	// character, the < symbol. We don't want to include
	// <title> as part of the final value, so let's offset
	// the index by the number of characers in <title>
	titleStartIndex += len(selector.Selector) + 2

	// Find the index of the closing tag
	titleEndIndex := strings.Index(*pageContent, "</"+selector.Selector+">")
	if titleEndIndex == -1 {
		fmt.Printf("%s : No closing tag for title found.", selector.Id)
		os.Exit(0)
	}

	// (Optional)
	// Copy the substring in to a separate variable so the
	// variables with the full document data can be garbage collected
	pageTitle := []byte((*pageContent)[titleStartIndex:titleEndIndex])

	// Print out the result
	fmt.Printf("%s : %s\n", selector.Id, pageTitle)
}

func SelectorImage(pageContent *string, selector *Selectors) {

	re := regexp.MustCompile("<img([\\w\\W]+?)>")
	comments := re.FindAllString(*pageContent, -1)
	re1 := regexp.MustCompile("src=([\\w\\W]+?) ")
	comments = re1.FindAllString(comments[0], -1)
	if comments == nil {
		fmt.Printf("%s : No element found", selector.Id)
		os.Exit(0)
	}

	// Print out the result
	fmt.Printf("%s : %s\n", selector.Id, comments)
}
func writeToFile() {
	message := []byte("Hello, Gophers!")
	err := ioutil.WriteFile(outputJSON, message, 0644)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	// readScrapingJSON()
	scraper()
}
