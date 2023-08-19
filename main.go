package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
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

func saveToFile(input string, filename string) {
	err := ioutil.WriteFile(filename, []byte(input), 0666)
	if err != nil {
		log.Fatal(err)
	}
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
	fullDoc, _ := goquery.OuterHtml(doc.Find("html"))
	fmt.Printf("Full doc: %s", fullDoc)
	saveToFile(fullDoc, "resp.html")
	// Find the review items
	doc.Find(".search-page__result").Each(func(i int, s *goquery.Selection) {
		// for each print the html
		html, _ := goquery.OuterHtml(s)
		fmt.Printf("HTML %d: %s\n", i, html)

	})

}
