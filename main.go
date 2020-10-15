package main

import (
	"encoding/json"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const (
	settingsConfig = "config.json"
	scrapingJSON   = "scraping.json"
	outputJSON     = "output.json"
)

func readSettingsJSON() {
	data, err := ioutil.ReadFile(settingsConfig)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Contents of file:", string(data))
	readScrapingJSON()
}

func readScrapingJSON() {
	data, err := ioutil.ReadFile(scrapingJSON)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Contents of file:", string(data))
	scraper()
}

func scraper() {
	fmt.Println("Test")
}

func writeToFile() {
	message := []byte("Hello, Gophers!")
	err := ioutil.WriteFile(outputJSON, message, 0644)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	readSettingsJSON()
}
