package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

var searchPostcode string = "GU7%201EU" // gudulming station
var searchRadius string = "100"         // miles
var searchMinimumModelYear string = "2020"
var autoTraderBaseURL string = "https://www.autotrader.co.uk/"

func main() {
	scrapeAutoTrader()
}

func fetchHTML(url string) *http.Response {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	return res
}

func scrapeAutoTrader() {
	currentPage := 1
	res := fetchHTML(autoTraderBaseURL + fmt.Sprintf("/car-search?year-from=%s&postcode=%s&radius=%s&page=%v", searchMinimumModelYear, searchPostcode, searchRadius, currentPage))
	defer res.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("li").Find("section").Each(func(i int, s *goquery.Selection) {
		// for each get the car name
		carName := s.Find("h3").Text()
		// For each item found, get the title
		fmt.Printf("Car %d: %s\n", i, carName)
	})
}
