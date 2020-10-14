package main

import (
	// import standard libraries
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	// import third party libraries
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
)

func readSettingsJSON() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Contents of file:", string(data))
	readScrapingJSON()
}

func readScrapingJSON() {
	data, err := ioutil.ReadFile("scraping.json")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Contents of file:", string(data))
	scraper()
}

func scraper() {
	fmt.Println("Test")
}

func main() {
	readSettingsJSON()
}
