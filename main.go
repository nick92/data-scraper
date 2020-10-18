package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

const (
	settingsConfig = "settings.json"
	scrapingJSON   = "stackover-flow.json"
	outputJSON     = "output.json"
)

var (
	config = &tls.Config{
		InsecureSkipVerify: true,
	}
	transport = &http.Transport{
		TLSClientConfig: config,
	}
	netClient = &http.Client{
		Transport: transport,
	}
)

type Selectors struct {
	ID              string
	Type            string
	ParentSelectors []string
	Selector        string
	Multiple        bool
	Regex           string
	Delay           int
}

type Scraping struct {
	StartUrl  []string
	ID        string `json:"_id,omitempty"`
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

func readSiteMap() *Scraping {
	// open the file and read the file
	data, err := ioutil.ReadFile(scrapingJSON)

	var scrape Scraping
	err = json.Unmarshal(data, &scrape)

	// log any errors
	if err != nil {
		log.Println(err)
	}

	return &scrape
}

func SelectorText(doc *goquery.Document, selector *Selectors) []string {
	// Find the review items
	// fmt.Println(selector.Selector)
	var text []string
	var matchText string
	doc.Find(selector.Selector).Each(func(i int, s *goquery.Selection) {

		if selector.Regex != "" {
			re := regexp.MustCompile(selector.Regex)
			matchText = re.FindString(s.Text())
		} else {
			matchText = s.Text()
		}
		text = append(text, matchText)
	})
	return text
}

func SelectorLink(doc *goquery.Document, selector *Selectors, baseURL string) []string {
	// Find the review items
	// fmt.Println(selector.Selector)
	var links []string
	doc.Find(selector.Selector).Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			fmt.Printf("href not found")
		}

		links = append(links, toFixedURL(href, baseURL))
		// if selector.Multiple == false {
		// 	return false
		// }
		// return true
	})
	return links
}

func crawlURL(href string) *goquery.Document {
	fmt.Printf("Crawling Url -> %v \n", href)
	response, err := netClient.Get(href)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func toFixedURL(href, baseURL string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	toFixedURI := base.ResolveReference(uri)

	return toFixedURI.String()
}

func getSiteMap(startURL []string, selector *Selectors) *Scraping {

	baseSiteMap := readSiteMap()
	newSiteMap := new(Scraping)
	newSiteMap.ID = selector.ID
	newSiteMap.StartUrl = startURL
	newSiteMap.Selectors = baseSiteMap.Selectors
	return newSiteMap
}

func scraper(siteMap *Scraping, parent string) interface{} {

	output := make(map[string]interface{})
	for _, startURL := range siteMap.StartUrl {
		linkOutput := make(map[string]interface{})
		for _, selector := range siteMap.Selectors {
			if parent == selector.ParentSelectors[0] {
				doc := crawlURL(startURL)
				if selector.Type == "SelectorText" {
					resultText := SelectorText(doc, &selector)
					linkOutput[selector.ID] = resultText
				} else if selector.Type == "SelectorLink" {
					links := SelectorLink(doc, &selector, startURL)
					fmt.Println(links)

					newSiteMap := getSiteMap(links, &selector)
					result := scraper(newSiteMap, selector.ID)
					linkOutput[selector.ID] = result
				}
			}
		}
		output[startURL] = linkOutput
	}
	return output
}

func main() {
	// readScrapingJSON()
	siteMap := readSiteMap()
	finalOutput := scraper(siteMap, "_root")

	file, err := json.MarshalIndent(finalOutput, "", " ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(file))
	_ = ioutil.WriteFile("output.json", file, 0644)

}
