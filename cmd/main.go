package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func main() {
	output := flag.String("o", "output.txt", "output file")
	flag.Parse()

	outputFile, err := os.OpenFile(*output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	c := colly.NewCollector(colly.MaxDepth(4))
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	c.OnHTML(".paginator > ul > li", func(element *colly.HTMLElement) {
		href, ok := element.DOM.Children().First().Attr("href")
		if ok {
			ok, err := element.Request.HasVisited(href)
			if err != nil {
				fmt.Println(err)
				return
			}

			if !ok {
				element.Request.Visit(href)
			}
		}
	})

	c.OnHTML(".message-body", func(element *colly.HTMLElement) {
		text := element.DOM.Contents().FilterFunction(func(_ int, selection *goquery.Selection) bool {
			isMessageQuote := selection.HasClass("quote") && selection.Children().First().HasClass("quote-nick")
			return !isMessageQuote
		}).Text()

		outputFile.WriteString(text)
		outputFile.WriteString("\n")
	})

	c.OnHTML(".thread", func(element *colly.HTMLElement) {
		href, ok := element.DOM.Children().First().Attr("href")
		if ok {
			ok, err := element.Request.HasVisited(href)
			if err != nil {
				fmt.Println(err)
				return
			}

			if !ok {
				element.Request.Visit(href + "/1")
			}
		}
	})

	c.OnHTML(".main_left .level-three a", func(element *colly.HTMLElement) {
		href, ok := element.DOM.Attr("href")
		if ok {
			ok, err := element.Request.HasVisited(href)
			if err != nil {
				fmt.Println(err)
				return
			}

			if !ok {
				element.Request.Visit(href)
			}
		}
	})

	c.Visit("https://www.familjeliv.se/forum/hot/1")
	c.Wait()
}
