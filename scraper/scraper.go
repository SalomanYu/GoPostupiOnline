package scraper

import (
	"strings"
	"strconv"
	"regexp"
	"log"
	"fmt"

	
	"github.com/SalomanYu/GoPostupiOnline/excel"
	"github.com/SalomanYu/GoPostupiOnline/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gocolly/colly"
)

var (
	currentInsitutionId    string
	currentSpecId          string
	formEducation       = []string{"specialnosti/bakalavr/", "specialnosti/specialist/", "specialnosti/magistratura/"}
	Headers 			= map[string]string{
							"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
							"sec-ch-ua":  `Google Chrome";v="105", "Not)A;Brand";v="8", "Chromium";v="105"`,
							"accept":     "image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
							"cookie":     "yandexuid=6850906421666216763; yabs-sid=1696581601666216766; yuidss=6850906421666216763; ymex=1981576766.yrts.1666216766#1981576766.yrtsi.1666216766; gdpr=0; _ym_uid=1666216766168837185; _ym_d=1666216766; yandex_login=rosya-8; i=Peh4utbtslQvge42D7cbDtH7CwXIiDs5Yp6IXWYsxx/SEQD1HtUncw/qqJV7NXqNqOS81fsaJSedcq/Ds9+yOfVKCNQ=; is_gdpr=0; skid=6879224341667473690; ys=udn.cDrQr9GA0L7RgdC70LDQsg%3D%3D#c_chck.841052032; is_gdpr_b=CIyaHxCclAE=; Session_id=3:1668355426.5.0.1666216795333:P19ouQ:2f.1.2:1|711384492.0.2|3:10261113.753043.lm80KKusrHll2DmXDLpHMjsmBYY; sessionid2=3:1668355426.5.0.1666216795333:P19ouQ:2f.1.2:1|711384492.0.2|3:10261113.753043.fakesign0000000000000000000; _ym_isad=1; _ym_visorc=b",
						}
)


func scrapeBasicInfo(h *colly.HTMLElement) (basic models.BasicInfo) {
	// Если нет краткой инфы о вузе/специальности/программе, то возвращаем пустую структуру
	score_list := h.ChildText("p.list__score")
	if score_list == "" {
		return 
	}

	// Парсим краткое описание
	h.ForEach("div.list__info p", func(i int, e *colly.HTMLElement) {
		if i == 1{
			basic.Description = e.Text
		}
	})

	// Парсим блок и с инфой о бюджетных и платных местах и баллах для поступления
	h.ForEach("div.list__score-wrap p", func(i int, e *colly.HTMLElement) {
		if strings.Contains(e.Text, "бал.бюджет"){
			score, err := strconv.ParseFloat(e.ChildText("b"), 64)
			check_err(err)
			basic.BudgetScore = score
		}else if strings.Contains(e.Text, "бал.платно"){
			score, err := strconv.ParseFloat(e.ChildText("b"), 64)
			check_err(err)
			basic.PaymentScore = score
		}else if strings.Contains(e.Text, "бюджетных мест") && !strings.Contains("нет", e.ChildText("b")){
			places := e.ChildText("b")
			basic.BudgetPlaces = places
		}else if strings.Contains(e.Text, "платных мест") && !strings.Contains("нет", e.ChildText("b")){
			places := e.ChildText("b")
			basic.PaymentPlaces = places
		}
	})

	// Наполняем оставшим контентом структуру
	basic.Url = h.ChildAttr("a", "href")
	basic.Image = h.ChildAttr("img", "data-dt")
	basic.Cost = h.ChildText("span.list__price b")
	basic.Name = h.ChildText("h2.list__h")
	basic.Logo = h.ChildAttr("img.list__img-sm", "src")

	re := regexp.MustCompile(`[а-яА-я]`)
	basic.Direction = re.FindString(h.ChildText("p.list__pre"))
	return 
}


func ScrapeInstitution(h *colly.HTMLElement) {
	institution := models.InstitutionInfo{}
	institution.ID = primitive.NewObjectID()
	basic := scrapeBasicInfo(h)
	if basic == (models.BasicInfo{}){
		return 
	}
	c := colly.NewCollector()

	// Меняем действующий айди вуза на актуальный
	institution.InstitutionId = strings.Split(basic.Url, "/")[len(strings.Split(basic.Url, "/"))-2]
	currentInsitutionId = institution.InstitutionId
	institution.Base = basic

	// Пробуем достать полное название Вуза и поменять его с тем, что предлагает метод scrapemodels.BasicInfo 
	c.OnHTML("h1.bg-nd__h", func(h *colly.HTMLElement) {
		name := h.Text
		if name != "" {
			institution.Base.Name = name
		}
	})

	// Достаем описание вуза 
	c.OnHTML("div.descr-min", func(h *colly.HTMLElement) {
		institution.Description = h.Text
	})

	// Парсим факты вуза и добавляем их к описанию Вуза
	c.OnHTML("ul.facts-list-nd li", func(h *colly.HTMLElement) {
		facts := h.Text
		institution.Description += "\nФакты: " + facts
	})

	c.Post(basic.Url, Headers)
	scrapeContacts(basic.Url + "contacts/")
	
	// Сохраняем вуз
	excel.AddVuz(&institution)
	log.Println("Parsed Vuz:")

	// Забрали инфу о Вузе и его контактах, теперь парсим специальности этого вуза
	for _, form := range formEducation {
		scrapeSpecializations(basic.Url + form)
	}
}

// $Env:GOOS = "linux"; $Env:GOARCH = "amd64"

func scrapeSpecializations(url string) {	
	page := 1
	for {
		hasSpecs := false
		channelSpec := make(chan models.SpecializationInfo)
		countSpecs := 0

		c := colly.NewCollector()
		c.OnHTML("div.list-cover li", func(h *colly.HTMLElement) {
			hasSpecs = true
			countSpecs++
			go scrapeOneSpecialization(h, channelSpec)

		})
		c.Post(fmt.Sprintf("%s?page_num=%d", url, page), Headers)
		
		for i:=0; i<countSpecs; i++ {
			spec, ok := <- channelSpec
			if ok == false{
				break
			}
			excel.AddSpecialization(&spec)
			log.Println("Specialization:")
			scrapePrograms(spec.Base.Url, spec.SpecId)
		}
		// Проверяем, есть ли на странице специализации или мы вышли за рамки пагинации
		if hasSpecs == false {
			break
		} else {
			page++
		}
	}
}


func scrapeOneSpecialization(h *colly.HTMLElement, channel chan models.SpecializationInfo){ 
	specialization := models.SpecializationInfo{}
	basic := scrapeBasicInfo(h)
	if basic == (models.BasicInfo{}){
		return
	}
	specialization.ID = primitive.NewObjectID()
	currentSpecId = strings.Split(basic.Url, "/")[len(strings.Split(basic.Url, "/"))-2]
	specialization.SpecId = currentSpecId
	specialization.Base = basic
	specialization.InstitutionId = currentInsitutionId

	c := colly.NewCollector()
	// Переходим на страницу специализации и пытаемся спарсить подробное описание, если его нет. То будет краткое описание в specialization.Base.Description
	c.OnHTML("div.descr-max", func(h *colly.HTMLElement) {
		specialization.Description = h.Text
	})
	
	c.Post(basic.Url, Headers)
	channel <- specialization
}


func scrapePrograms(url string, specId string) {
	page := 1
	for {
		hasPrograms := false
		channelProgram := make(chan models.ProgramInfo)
		countPrograms := 0

		c := colly.NewCollector()
		c.OnHTML("div.list-cover li", func(h *colly.HTMLElement) {
			hasPrograms = true
			countPrograms++
			go scrapeOneProgram(h, channelProgram)

		})
		c.Post(fmt.Sprintf("%s?page_num=%d", url, page), Headers)
		
		for i:=0; i<countPrograms; i++ {
			program, ok := <- channelProgram
			if ok == false{
				break
			}
			program.SpecId = specId
			excel.AddProgram(&program)
			log.Println("Program: ")
			ScrapeProfession(program.Base.Url)
		}

		// Проверяем, есть ли на странице программы или мы вышли за рамки пагинации
		if hasPrograms == false {
			break
		} else {
			page++
		}
	}
}

func scrapeOneProgram(h *colly.HTMLElement, channel chan models.ProgramInfo) { 
	program := models.ProgramInfo{}
	program.Base = scrapeBasicInfo(h)
	if program.Base == (models.BasicInfo{}){
		return
	}
	program.ID = primitive.NewObjectID()
	program.InstitutionId = currentInsitutionId
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

	c.Post(program.Base.Url, Headers)
	channel <- program
}

func ScrapeProfession(programUrl string) {
	newProfession := models.ProfessionInfo{}
	newProfession.ID = primitive.NewObjectID()
	c := colly.NewCollector()

	// Забираем айдишник программы из url професси. Берем 477 из https://msk.postupi.online/vuz/mip/programma/477/
	programId := strings.Split(programUrl, "/")[len(strings.Split(programUrl, "/"))-2]

	// Наполняем контентом нашу профессию
	c.OnHTML("li.list-col", func(h *colly.HTMLElement) {
		newProfession.Name = h.ChildText("h2")
		newProfession.Image = h.ChildAttr("img.img-load", "data-dt")
		newProfession.ProgramId = programId

		// Сохраняем ее
		excel.AddProression(&newProfession)
		log.Println("Profession: ")
	})

	// Отправляем запрос на страницу с профессиями конкретной программы 
	c.Post(programUrl + "professii/", Headers)
}


func scrapeContacts(url string) {
	contact := models.ContactsInfo{}
	contact.ID = primitive.NewObjectID()
	c := colly.NewCollector()

	// Ищем 4 span блока с контактами вуза
	c.OnHTML("section.section-box", func(h *colly.HTMLElement) {
		contactList := [4]string{}
		h.ForEach("span", func(i int, e *colly.HTMLElement) {
			contactList[i] = e.Text
		})

		// Наполняем контентом контакты
		if contactList[0] != "" {
			contact.WebSite = contactList[0]
			contact.Email = contactList[1]
			contact.Phone = contactList[2]
			contact.Address = contactList[3]
			contact.VuzId = currentInsitutionId
		}
	})
	c.Post(url, Headers)
	excel.AddContacts(&contact)
	log.Println("Contact: ")
}

func check_err(err error){
	if err != nil{
		log.Fatal(err)
		panic(err)
	}
}