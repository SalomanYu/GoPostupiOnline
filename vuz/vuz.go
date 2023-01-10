package vuz

import (
	"fmt"
	"time"

	"github.com/SalomanYu/GoPostupiOnline/vuz/scraper"
	"github.com/gocolly/colly"
)

var (
	domain = "https://postupi.online/vuzi/"
	pageCount = 52
	defaultFormEducation = "Бакалавриат"
	tagNameForListBlocks = "list-unstyled.list-wrap"
)

func Start() {
	for i := 1; i <= pageCount; i++ {
		url := domain + fmt.Sprintf("?page_num=%d", i)
		c := colly.NewCollector()
		c.OnError(func(r *colly.Response, err error) {
			if r.StatusCode >= 500{
				return
			}
		})
		c.SetRequestTimeout(30 * time.Second)
		c.OnHTML("div.list-cover li.list", func(h *colly.HTMLElement) {
				scraper.ScrapeVuz(h)
		})
		err := c.Post(url, scraper.Headers)
		if err != nil {
			fmt.Println("Catched the error. Program stopped to sleep of 10 seconds.")
			time.Sleep(10*time.Second)
			err = c.Post(url, scraper.Headers)
		}
	}
}

func check_err(err error){
	if err != nil{
		panic(err)
	}
}


