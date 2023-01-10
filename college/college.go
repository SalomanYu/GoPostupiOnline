package college

import (
	"fmt"
	"time"

	"github.com/SalomanYu/GoPostupiOnline/college/scraper"
	"github.com/gocolly/colly"
)

var (
	domain = "https://postupi.online/ssuzy/" // Переменная
	pageCount = 62 // Переменная
	defaultFormEducation = "Подготовка квалифицированных рабочих (служащих)" // Переменная
	tagNameForListBlocks = "list-wrap" // Переменная
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


