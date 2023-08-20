package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var searchPostcode string = "GU7%201EU" // gudulming station
var searchRadius string = "200"         // miles
var searchMinimumModelYear string = "2020"
var autoTraderBaseURL string = "https://www.autotrader.co.uk/"

type car struct {
	name        string
	keySpecs    []string
	price       string
	detailsLink string
}

func main() {
	cars := scrapeAutoTrader()
	fmt.Println(cars)
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

func scrapeAutoTrader() []car {
	currentPage := 1
	res := fetchHTML(autoTraderBaseURL + fmt.Sprintf("/car-search?year-from=%s&postcode=%s&radius=%s&page=%v", searchMinimumModelYear, searchPostcode, searchRadius, currentPage))
	defer res.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fullDoc, _ := goquery.OuterHtml(doc.Find("html"))
	saveToFile(fullDoc, "resp.html")

	var cars []car

	// Find the review items
	doc.Find("li[class='search-page__result']").Each(func(i int, s *goquery.Selection) {
		carDetails := s.Find("section[class='product-card-details']")
		carName := carDetails.Find("h3[class='product-card-details__title']")
		var keySpecs []string
		carDetails.Find("li").Each(func(i int, k *goquery.Selection) {
			keySpecs = append(keySpecs, k.Text())
		})

		carPrice := s.Find("div[class='product-card-pricing__price']")

		detailsLinkAnchorTag := s.Find("a[data-results-nav-fpa]")
		detailsLink, exists := detailsLinkAnchorTag.Attr("href")
		if exists != true {
			log.Fatal("No link found for car")
		}

		a, _ := goquery.OuterHtml(carPrice)
		fmt.Printf("Price %d: %s\n", i, a)

		car := car{name: carName.Text(), price: carPrice.Text(), detailsLink: detailsLink, keySpecs: keySpecs}
		cars = append(cars, car)
	})

	return cars

}
