package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/SalomanYu/GoPostupiOnline/storages/mongo"
	"github.com/SalomanYu/GoPostupiOnline/models"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var currentVuzId, currentSpecId string
var formEducation = []string{"specialnosti/bakalavr/", "specialnosti/specialist/", "specialnosti/magistratura/"}
var Headers = map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
			"sec-ch-ua":  `Google Chrome";v="105", "Not)A;Brand";v="8", "Chromium";v="105"`,
			"accept":     "image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"cookie":     "yandexuid=6850906421666216763; yabs-sid=1696581601666216766; yuidss=6850906421666216763; ymex=1981576766.yrts.1666216766#1981576766.yrtsi.1666216766; gdpr=0; _ym_uid=1666216766168837185; _ym_d=1666216766; yandex_login=rosya-8; i=Peh4utbtslQvge42D7cbDtH7CwXIiDs5Yp6IXWYsxx/SEQD1HtUncw/qqJV7NXqNqOS81fsaJSedcq/Ds9+yOfVKCNQ=; is_gdpr=0; skid=6879224341667473690; ys=udn.cDrQr9GA0L7RgdC70LDQsg%3D%3D#c_chck.841052032; is_gdpr_b=CIyaHxCclAE=; Session_id=3:1668355426.5.0.1666216795333:P19ouQ:2f.1.2:1|711384492.0.2|3:10261113.753043.lm80KKusrHll2DmXDLpHMjsmBYY; sessionid2=3:1668355426.5.0.1666216795333:P19ouQ:2f.1.2:1|711384492.0.2|3:10261113.753043.fakesign0000000000000000000; _ym_isad=1; _ym_visorc=b",
		}


func scrapeBasic(h *colly.HTMLElement) (basic models.Basic) {
	if h.ChildText("p.list__score") == "" {
		return 
	}

	h.ForEach("div.list__info p", func(i int, e *colly.HTMLElement) {
		if i == 1{
			basic.Description = e.Text
		}
	})

	h.ForEach("div.list__score-wrap p", func(i int, e *colly.HTMLElement) {
		if strings.Contains(e.Text, "бал.бюджет"){
			scores, err := strconv.ParseFloat(e.ChildText("b"), 64)
			checkErr(err)
			basic.BudgetScores = scores
		}else if strings.Contains(e.Text, "бал.платно"){
			scores, err := strconv.ParseFloat(e.ChildText("b"), 64)
			checkErr(err)
			basic.PaymentScores = scores
		}else if strings.Contains(e.Text, "бюджетных мест") && !strings.Contains("нет", e.ChildText("b")){
			places := e.ChildText("b")
			basic.BudgetPlaces = places
		}else if strings.Contains(e.Text, "платных мест") && !strings.Contains("нет", e.ChildText("b")){
			places := e.ChildText("b")
			basic.PaymentPlaces = places
		}
	})

	basic.Url = h.ChildAttr("a", "href")
	basic.Image = h.ChildAttr("img", "data-dt")
	basic.Cost = h.ChildText("span.list__price b")
	basic.Name = h.ChildText("h2.list__h")
	basic.Logo = h.ChildAttr("img.list__img-sm", "src")
	basic.Direction = h.ChildText("p.list__pre")
	return 
}


func ScrapeVuz(h *colly.HTMLElement) {
	institution := models.Vuz{}
	institution.ID = primitive.NewObjectID()
	basic := scrapeBasic(h)
	if basic == (models.Basic{}){
		return 
	}
	c := colly.NewCollector()
	institution.VuzId = strings.Split(basic.Url, "/")[len(strings.Split(basic.Url, "/"))-2]
	currentVuzId = institution.VuzId
	institution.Base = basic

	c.OnHTML("h1.bg-nd__h", func(h *colly.HTMLElement) {
		name := h.Text
		if name != "" {
			institution.Base.Name = name
		}
	})
	c.OnHTML("div.descr-min", func(h *colly.HTMLElement) {
		institution.Description = h.Text
	})
	c.OnHTML("ul.facts-list-nd li", func(h *colly.HTMLElement) {
		facts := h.Text
		institution.Description += "\nФакты: " + facts
	})

	err := c.Post(basic.Url, Headers)
	checkErr(err)
	scrapeContacts(basic.Url + "contacts/")
	err = mongo.AddVuz(&institution)
	checkErr(err)
	log.Println("Parsed Vuz:")
	for _, form := range formEducation {
		scrapeManySpecializations(basic.Url + form)
	}
}

// $Env:GOOS = "linux"; $Env:GOARCH = "amd64"

func scrapeManySpecializations(url string) {	
	page := 1
	for {
		hasSpecs := false
		channelSpec := make(chan models.Specialization)
		countSpecs := 0

		c := colly.NewCollector()
		c.OnHTML("div.list-cover li", func(h *colly.HTMLElement) {
			hasSpecs = true
			countSpecs++
			go scrapeOneSpecialization(h, channelSpec)

		})
		err := c.Post(fmt.Sprintf("%s?page_num=%d", url, page), Headers)
		checkErr(err)

		for i:=0; i<countSpecs; i++ {
			spec, ok := <- channelSpec
			if ok == false{
				log.Println("ok == false в канале специализаций")
				break
			}
			err = mongo.AddSpecialization(&spec)
			checkErr(err)
			// log.Println("Specialization:")
			scrapeManyPrograms(spec.Base.Url, spec.SpecId)
		}
		close(channelSpec)
		log.Println("Закончили читать канал специализаций")
		if hasSpecs == false {
			break
		} else {
			page++
		}
	}
}


func scrapeOneSpecialization(h *colly.HTMLElement, channel chan models.Specialization){ 
	specialization := models.Specialization{}
	basic := scrapeBasic(h)
	if basic == (models.Basic{}){
		return
	}
	specialization.ID = primitive.NewObjectID()
	currentSpecId = strings.Split(basic.Url, "/")[len(strings.Split(basic.Url, "/"))-2]
	specialization.SpecId = currentSpecId
	specialization.Base = basic
	specialization.VuzId = currentVuzId
	
	re := regexp.MustCompile(`[а-яА-я]`)
	specialization.Base.Direction = re.FindString(h.ChildText("p.list__pre"))

	c := colly.NewCollector()
	c.OnHTML("div.descr-max", func(h *colly.HTMLElement) {
		specialization.Description = h.Text
	})
	
	err := c.Post(basic.Url, Headers)
	checkErr(err)
	channel <- specialization
}


func scrapeManyPrograms(url string, specId string) {
	page := 1
	for {
		hasPrograms := false
		channelProgram := make(chan models.Program)
		countPrograms := 0

		c := colly.NewCollector()
		c.OnHTML("div.list-cover li", func(h *colly.HTMLElement) {
			hasPrograms = true
			countPrograms++
			go scrapeOneProgram(h, channelProgram)

		})
		err := c.Post(fmt.Sprintf("%s?page_num=%d", url, page), Headers)
		checkErr(err)
		for i:=0; i<countPrograms; i++ {
			program, ok := <- channelProgram
			if ok == false{
				log.Println("ok==false в канале программ")
				break
			}
			program.SpecId = specId
			err = mongo.AddProgram(&program)
			checkErr(err)
			log.Println("Program: ")
			ScrapeProfessions(program.Base.Url)
		}
		close(channelProgram)
		log.Println("закончили работу в канале программ")
		if hasPrograms == false {
			break
		} else {
			page++
		}
	}
}

func scrapeOneProgram(h *colly.HTMLElement, channel chan models.Program) { 
	program := models.Program{}
	program.Base = scrapeBasic(h)
	if program.Base == (models.Basic{}){
		return
	}
	program.ID = primitive.NewObjectID()
	program.VuzId = currentVuzId
	program.ProgramId = strings.Split(program.Base.Url, "/")[len(strings.Split(program.Base.Url, "/"))-2]
	c := colly.NewCollector()

	// Забираем полное наименование программы 
	c.OnHTML("h1[id=prTitle]", func(h *colly.HTMLElement) {
		program.Base.Name = h.Text
	})

	// Выбираем уровень подготовки с помощью поиска подходящего варианта через регулярное выражение
	c.OnHTML("div.detail-box", func(h *colly.HTMLElement) {
		re := regexp.MustCompile("Бакалавриат|Специалитет|Магистратура|Подготовка специалистов среднего звена|Подготовка квалифицированных рабочих (служащих)")
		program.Form = re.FindString(h.Text)
	})

	// Парсим описание программы
	c.OnHTML("div.descr-max", func(h *colly.HTMLElement) {
		program.Description = h.Text
	})

	// Парсим предметы ЕГЭ
	c.OnHTML("div.score-box-wrap", func(h *colly.HTMLElement) {
		h.ForEach("div.score-box", func(i int, h *colly.HTMLElement) {
			if i == 1{
				h.ForEach("div.score-box__item", func(i int, h *colly.HTMLElement) {
					// score := h.ChildText("span.score-box__score") Не берем т.к без авторизации не показывается количество баллов
					exam := strings.Split(h.Text, "или")[0] //  Обрезаем строку и оставляем только матешу: Математика или другиеили Иностранный языкили Обществознание
					program.Exams = append(program.Exams, exam)
				})
			}
		})
	})

	err := c.Post(program.Base.Url, Headers)
	checkErr(err)
	channel <- program
}

func ScrapeProfessions(programUrl string) {
	
	c := colly.NewCollector()

	// Забираем айдишник программы из url професси. Берем 477 из https://msk.postupi.online/vuz/mip/programma/477/
	programId := strings.Split(programUrl, "/")[len(strings.Split(programUrl, "/"))-2]

	c.OnHTML("div.list-cover li", func(h *colly.HTMLElement) {
		newProfession := models.Profession{}
		newProfession.ID = primitive.NewObjectID()
		newProfession.Name = h.ChildText("h2")
		newProfession.Image = h.ChildAttr("img.img-load", "data-dt")
		newProfession.ProgramId = programId

		err := mongo.AddProfession(&newProfession)
		checkErr(err)
		// log.Println("Profession: ")
	})

	err := c.Post(programUrl + "professii/", Headers)
	checkErr(err)
}


func scrapeContacts(url string) {
	contact := models.Contacts{}
	contact.ID = primitive.NewObjectID()
	c := colly.NewCollector()

	// Ищем 4 span блока с контактами вуза
	c.OnHTML("section.section-box", func(h *colly.HTMLElement) {
		contactList := [4]string{}
		h.ForEach("span", func(i int, e *colly.HTMLElement) {
			contactList[i] = e.Text
		})

		if contactList[0] != "" {
			contact.WebSite = contactList[0]
			contact.Email = contactList[1]
			contact.Phone = contactList[2]
			contact.Address = contactList[3]
			contact.VuzId = currentVuzId
		}
	})
	err := c.Post(url, Headers)
	checkErr(err)
	err = mongo.AddContacts(&contact)
	checkErr(err)
	log.Println("Contact: ")
}

func checkErr(err error){
	if err != nil{
		log.Fatal(err)
		panic(err)
	}
}