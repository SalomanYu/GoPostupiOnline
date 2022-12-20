package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SalomanYu/GoPostupiOnline/models"

	"github.com/gocolly/colly"
)

func scrapeBasic(h *colly.HTMLElement) (basic models.Basic) {
	if !BlockContainsImportantInfo(h){
		return
	}
	basic.Description = getFullDescription(h)
	basic.Scores = getScores(h)
	basic.Url = h.ChildAttr("a", "href")
	basic.Image = h.ChildAttr("img", "data-dt")
	basic.Cost = h.ChildText("span.list__price b")
	basic.Name = h.ChildText("h2.list__h")
	basic.Logo = h.ChildAttr("img.list__img-sm", "src")
	basic.Direction = h.ChildText("p.list__pre")
	return 
}

func hasBasicInfo(basic models.Basic) bool {
	if basic == (models.Basic{}){
		return false
	}
	return true
}
func getFullName(html *colly.HTMLElement) (name string){
	name = html.ChildText("h1[id=prTitle]")
	return
}

func getBodyCodeFromUrl(url string) (body *colly.HTMLElement, err error) {
	c := colly.NewCollector()
	badGateway := checkBadGateway(c)
	if badGateway{
		return
	}
	c.SetRequestTimeout(30 * time.Second)
	c.OnHTML("body", func(h *colly.HTMLElement) {
		body = h
	})
	err = c.Post(url, Headers)
	return
}

func getFacts(body *colly.HTMLElement) (facts string){
	var facts_list []string
	body.ForEach("ul.facts-list-nd li", func(i int, h *colly.HTMLElement) {
		facts_list = append(facts_list, h.Text)
	})
	facts = "Факты\n:" + strings.Join(facts_list, ".")
	return
}

func findHtmlBlocks(blocksUrl string) (blocks []*colly.HTMLElement, err error){
	c := colly.NewCollector()
	badGateway := checkBadGateway(c)
	if badGateway{
		return
	}
	c.SetRequestTimeout(30 * time.Second)
	c.OnHTML("div.list-cover", func(h *colly.HTMLElement) {
			h.ForEach("li.list", func(i int, h *colly.HTMLElement) {
				blocks = append(blocks, h)
			})
		})
	err = c.Post(blocksUrl, Headers)
	checkErr(err)
	return
}

func getFullDescription(html *colly.HTMLElement) (description string){
	html.ForEach("div.list__info p", func(i int, e *colly.HTMLElement) {
		if i == 1{
			description = e.Text
		}
	})
	return
}
func getMiniDescription(body *colly.HTMLElement) (description string){
	description = body.ChildText("div.descr-min")
	return
}

func BlockContainsImportantInfo(html *colly.HTMLElement) bool {
	if html.ChildText("p.list__score") == "" {
		return false 
	}
	return true
}

func getScores(html *colly.HTMLElement)(scores models.Scores){
	html.ForEach("div.list__score-wrap p", func(i int, e *colly.HTMLElement) {
		switch true{
		case strings.Contains(e.Text, "бал.бюджет"):
			digit, err := strconv.ParseFloat(e.ChildText("b"), 64)
			checkErr(err)
			scores.PointsBudget = digit
		case strings.Contains(e.Text, "бал.платно"):
			digit, err := strconv.ParseFloat(e.ChildText("b"), 64)
			checkErr(err)
			scores.PointsPayment = digit
		case strings.Contains(e.Text, "бюджетных мест") && !strings.Contains("нет", e.ChildText("b")):
			scores.PlacesBudget = e.ChildText("b")
		case strings.Contains(e.Text, "платных мест") && !strings.Contains("нет", e.ChildText("b")):
			scores.PlacesPayment = e.ChildText("b")
		}
	})
	return
}

func getSpecializationDirection(body *colly.HTMLElement, specId string) (direction string){
	direction = strings.ReplaceAll(body.ChildText("p.list__pre"), specId, "")
	return
}
func getProgramDirection(html *colly.HTMLElement, specId string) (direction string){
	direction = strings.Split(html.ChildText("p.list__pre"), specId)[1]
	return
} 

func getSpecializationDescription(body *colly.HTMLElement)(description string){
	description = body.ChildText("div.descr-max")
	return
}

func getFormEducation(body *colly.HTMLElement) (form string) {
	re := regexp.MustCompile("Бакалавриат|Специалитет|Магистратура|Подготовка специалистов среднего звена|Подготовка квалифицированных рабочих (служащих)")
	form = re.FindString(body.ChildText("div.detail-box"))
	return
}

func getSubjects(body *colly.HTMLElement) (subjects []string){
	body.ForEach("div.score-box-wrap div.score-box", func(i int, h *colly.HTMLElement) {
		if i == 1{
			h.ForEach("div.score-box__item", func(i int, h *colly.HTMLElement) {
				// score := h.ChildText("span.score-box__score") Не берем т.к без авторизации не показывается количество баллов
				exam := strings.Split(h.Text, "или")[0] //  Обрезаем строку и оставляем только матешу: Математика или другиеили Иностранный языкили Обществознание
				subjects = append(subjects, exam)
			})
		}
	})
	return
}

func checkBadGateway(c *colly.Collector) (badGateway bool){
	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode >= 500{
			badGateway = true
			return
		}
	})
	return
}

func checkErr(err error){
	if err != nil{
		// log.Fatal(err)
		panic(err)
	}
}