package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()

	ticker := []string{
		"MSFT",
		"AAPL",
		"PEL",
	}

	type Stock struct {
		company, price, change string
	}

	stocks := []Stock{}
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error while getting data for ", r.Request.URL, err)
	})

	c.OnHTML("div#quote-header-info", func(h *colly.HTMLElement) {
		stock := Stock{}
		stock.company = h.ChildText("h1")
		fmt.Println("Company:", stock.company)
		stock.price = h.ChildText("fin-streamer[data-field='regularMarketPrice']")
		fmt.Println("Price:", stock.price)
		stock.change = h.ChildText("fin-streamer[data-field='regularMarketChangePercent']")
		fmt.Println("Change:", stock.change)
		stocks = append(stocks, stock)
	})

	for _, t := range ticker {
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")
	}
	fmt.Println(stocks)

	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatalln("Error opening file")
	}
	writer := csv.NewWriter(file)
	headers := []string{
		"Company",
		"Price",
		"Change",
	}
	writer.Write(headers)
	for _, t := range stocks {
		stock := []string{
			t.company,
			t.price,
			t.change,
		}
		writer.Write(stock)
	}
	defer writer.Flush()

}
