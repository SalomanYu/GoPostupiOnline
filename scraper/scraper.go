package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/SalomanYu/GoPostupiOnline/models"
	"github.com/SalomanYu/GoPostupiOnline/storages/mongo"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var currentVuzId, currentSpecId string
// var formEducations = []string{"specialnosti/bakalavr/", "specialnosti/specialist/", "specialnosti/magistratura/"}


var Headers = map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
			"sec-ch-ua":  `Google Chrome";v="105", "Not)A;Brand";v="8", "Chromium";v="105"`,
			"accept":     "image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"cookie":     "yandexuid=6850906421666216763; yabs-sid=1696581601666216766; yuidss=6850906421666216763; ymex=1981576766.yrts.1666216766#1981576766.yrtsi.1666216766; gdpr=0; _ym_uid=1666216766168837185; _ym_d=1666216766; yandex_login=rosya-8; i=Peh4utbtslQvge42D7cbDtH7CwXIiDs5Yp6IXWYsxx/SEQD1HtUncw/qqJV7NXqNqOS81fsaJSedcq/Ds9+yOfVKCNQ=; is_gdpr=0; skid=6879224341667473690; ys=udn.cDrQr9GA0L7RgdC70LDQsg%3D%3D#c_chck.841052032; is_gdpr_b=CIyaHxCclAE=; Session_id=3:1668355426.5.0.1666216795333:P19ouQ:2f.1.2:1|711384492.0.2|3:10261113.753043.lm80KKusrHll2DmXDLpHMjsmBYY; sessionid2=3:1668355426.5.0.1666216795333:P19ouQ:2f.1.2:1|711384492.0.2|3:10261113.753043.fakesign0000000000000000000; _ym_isad=1; _ym_visorc=b",
			"Content-Type": "text/html",
		}


func ScrapeVuz(h *colly.HTMLElement) {
	// Надо еще рефачить и делить метод
	basic := scrapeBasic(h)
	if !hasBasicInfo(basic){
		return
	}
	institution := models.Vuz{}
	institution.ID = primitive.NewObjectID()
	institution.VuzId = strings.Split(basic.Url, "/")[len(strings.Split(basic.Url, "/"))-2]
	currentVuzId = institution.VuzId
	institution.Base = basic
	bodyHTML, err := getBodyCodeFromUrl(institution.Base.Url)
	if err != nil{
		fmt.Println("Catched the error. Program stopped to sleep of 10 seconds.")
		time.Sleep(10*time.Second)
		html, err := getBodyCodeFromUrl(institution.Base.Url)
		checkErr(err)
		bodyHTML = html
	
	}
	institution.Description = getMiniDescription(bodyHTML)
	fullname := getFullDescription(bodyHTML)
	if fullname != ""{
		institution.Base.Name = fullname
	}
	facts := getFacts(bodyHTML)
	if facts != ""{
		institution.Description += facts
	}
	err = scrapeContacts(basic.Url + "contacts/")
	if err != nil{
		fmt.Println("Catched the error. Program stopped to sleep of 10 seconds.")
		time.Sleep(10*time.Second)
		err = scrapeContacts(basic.Url + "contacts/")
		checkErr(err)
	}
	err = mongo.AddVuz(&institution)
	checkErr(err)
	log.Printf("Vuz:%s", institution.VuzId)
	var formEducations = []string{"specialnosti/spo/", "specialnosti/npo"}
	// var formEducations = []string{"specialnosti/bakalavr/", "specialnosti/specialist/", "specialnosti/magistratura/"}
	for _, form := range formEducations {
		scrapeVuzSpecializations(basic.Url + form)
	}
}

func scrapeVuzSpecializations(url string) {	
	page := 1
	for {
		specsBlocks, err := findHtmlBlocks(fmt.Sprintf("%s?page_num=%d", url, page))
		if err != nil{
			fmt.Println("Catched the error. Program stopped to sleep of 10 seconds.")
			time.Sleep(10*time.Second)
			blocks, err := findHtmlBlocks(fmt.Sprintf("%s?page_num=%d", url, page))
			checkErr(err)
			specsBlocks = blocks
		} 
		if len(specsBlocks) == 0 {
			break
		}
		var wg sync.WaitGroup
		wg.Add(len(specsBlocks))
		for _, spec := range specsBlocks{
			go scrapeSpecialization(spec, &wg)
		}
		wg.Wait()
		page++
	}
}

func scrapeSpecialization(h *colly.HTMLElement, wg *sync.WaitGroup){ 
	// надо дорефачить
	basic := scrapeBasic(h)
	if !hasBasicInfo(basic){
		return
	}
	specialization := models.Specialization{}
	specialization.ID = primitive.NewObjectID()
	currentSpecId = strings.Split(basic.Url, "/")[len(strings.Split(basic.Url, "/"))-2]
	specialization.SpecId = currentSpecId
	specialization.Base = basic
	specialization.VuzId = currentVuzId
	bodyHTML, err := getBodyCodeFromUrl(specialization.Base.Url)
	if err != nil{
		fmt.Println("Catched the error. Program stopped to sleep of 10 seconds.")
		time.Sleep(10*time.Second)
		html, err := getBodyCodeFromUrl(specialization.Base.Url)
		checkErr(err)
		bodyHTML = html
	}
	specialization.Base.Direction = getSpecializationDirection(h, specialization.SpecId)
	specialization.Description = getSpecializationDescription(bodyHTML)
	
	err = mongo.AddSpecialization(&specialization)
	checkErr(err)
	log.Printf("Specialization:%s", specialization.SpecId)
	scrapeSpecializationPrograms(specialization.Base.Url, specialization.SpecId)
	wg.Done()
}

func scrapeSpecializationPrograms(url string, specId string) {
	page := 1
	for {
		programsBlocks, err := findHtmlBlocks(fmt.Sprintf("%s?page_num=%d", url, page))
		if err != nil{
			fmt.Println("Catched the error. Program stopped to sleep of 10 seconds.")
			time.Sleep(10*time.Second)
			blocks, err := findHtmlBlocks(fmt.Sprintf("%s?page_num=%d", url, page))
			checkErr(err)
			programsBlocks = blocks
		}
		if len(programsBlocks) == 0{
			break
		}
		var wg sync.WaitGroup
		wg.Add(len(programsBlocks))
		for _, item := range programsBlocks{
			go scrapeProgram(item, specId, &wg)
		}
		wg.Wait()
		page++
		
	}
}

func scrapeProgram(h *colly.HTMLElement, specId string, wg *sync.WaitGroup){ 
	basic := scrapeBasic(h)
	if !hasBasicInfo(basic){
		return
	}
	program := models.Program{}
	program.Base = basic
	program.ID = primitive.NewObjectID()
	program.VuzId = currentVuzId
	program.ProgramId = strings.Split(program.Base.Url, "/")[len(strings.Split(program.Base.Url, "/"))-2]
	program.SpecId = specId
	bodyHTMl, err := getBodyCodeFromUrl(program.Base.Url)
	if err != nil{
		fmt.Println("Catched the error. Program stopped to sleep of 10 seconds.")
		time.Sleep(10*time.Second)
		html, err := getBodyCodeFromUrl(program.Base.Url)
		checkErr(err)
		bodyHTMl = html
	}
	program.Description = getSpecializationDescription(bodyHTMl)
	program.Base.Direction = getProgramDirection(h, specId)
	program.Form = getFormEducation(bodyHTMl)
	program.Exams = getSubjects(bodyHTMl)
	program.HasProfessions, err = ScrapeProfessions(program.Base.Url)
	if err != nil {
		fmt.Println("Catched the error. Program stopped to sleep of 10 seconds.")
		time.Sleep(10*time.Second)
		program.HasProfessions, err = ScrapeProfessions(program.Base.Url)
		checkErr(err)
	}
	fullname := getFullName(bodyHTMl)
	if fullname != ""{
		program.Base.Name = strings.Split(fullname, ":")[0] 
		if program.Form == "Магистратура"{ // Названия программ у магистратуры выглядят следующим образом: Профиль магистратуры "Управление свойствами нетканых материалов" РГУ им. А.Н. Косыгина, Москва
			re := regexp.MustCompile(`"(.*?)"`)
			program.Base.Name = re.FindString(program.Base.Name)
		}
	}
	err = mongo.AddProgram(&program)
	checkErr(err)
	log.Printf("Program:%s", program.ProgramId)
	wg.Done()
}

func ScrapeProfessions(programUrl string) (programHasProfessions bool, err error){
	programHasProfessions = false
	c := colly.NewCollector()
	badGateway := checkBadGateway(c)
	if badGateway{
		return
	}
	c.SetRequestTimeout(30 * time.Second)
	// Забираем айдишник программы из url професси. Берем 477 из https://msk.postupi.online/vuz/mip/programma/477/
	programId := strings.Split(programUrl, "/")[len(strings.Split(programUrl, "/"))-2]
	c.OnHTML("div.list-cover li", func(h *colly.HTMLElement) {
		programHasProfessions = true
		newProfession := models.Profession{}
		newProfession.ID = primitive.NewObjectID()
		newProfession.Name = h.ChildText("h2")
		newProfession.Image = h.ChildAttr("img.img-load", "data-dt")
		newProfession.ProgramId = programId

		err := mongo.AddProfession(&newProfession)
		checkErr(err)
	})

	err = c.Post(programUrl + "professii/", Headers)
	return 
}


func scrapeContacts(url string) (err error){
	contact := models.Contacts{}
	contact.ID = primitive.NewObjectID()
	c := colly.NewCollector()
	badGateway := checkBadGateway(c)
	if badGateway{
		return
	}
	c.SetRequestTimeout(30 * time.Second)

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
	err = c.Post(url, Headers)
	if err != nil {
		return
	}
	err = mongo.AddContacts(&contact)
	checkErr(err)
	log.Printf("Contact:%s", contact.VuzId)
	return
}

