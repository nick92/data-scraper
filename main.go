package main

import (
	"encoding/json"
	"fmt"
	// "golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const (
	settingsConfig = "settings.json"
	scrapingJSON   = "scraping.json"
	outputJSON     = "output.json"
)

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

func readScrapingJSON() {
	// open the file and read the file
	data, err := ioutil.ReadFile(scrapingJSON)
	// define data struture
	type Scraping struct {
		startUrl []string
		id    string
		parentSelectors []string
	}
	// define data struture
	var scrape Scraping
	err = json.Unmarshal(data, &scrape)
	// log any errors
	if err != nil {
		log.Println(err)
	}
	// lets just print for now
	fmt.Println("id", scrape.startUrl)
	// run the scraper and start scraping.
	// scraper()
}

func scraper() {
	response, err := http.Get("https://www.devdungeon.com")
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	// Read response data in to memory
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	// Create a regular expression to find comments
	re := regexp.MustCompile("<!--(.|\n)*?-->")
	comments := re.FindAllString(string(body), -1)
	if comments == nil {
		fmt.Println("No matches.")
	} else {
		for _, comment := range comments {
			fmt.Println(comment)
		}
	}
}

func writeToFile() {
	message := []byte("Hello, Gophers!")
	err := ioutil.WriteFile(outputJSON, message, 0644)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	readScrapingJSON()
}
