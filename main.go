package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SalomanYu/GoPostupiOnline/scraper"

	"github.com/gocolly/colly"
)

var (
	domain    = "https://postupi.online/vuzi/"
	pageCount = 52

)


func main() {
	start := time.Now().Unix()
	for i := 1; i <= pageCount; i++ {
		url := domain + fmt.Sprintf("?page_num=%d", i)
		c := colly.NewCollector()
		c.OnHTML("div.list-cover li", func(h *colly.HTMLElement) {
			scraper.ScrapeInstitution(h)
		})
		err := c.Post(url, scraper.Headers)
		check_err(err)
	}
	var a string
	log.Printf("\n\nTime: %d", time.Now().Unix()-start)
	fmt.Println("Program stoped.")
	fmt.Scan((&a))
}

func check_err(err error){
	if err != nil{
		panic(err)
	}
}


