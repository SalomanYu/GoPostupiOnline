package excel

import (
	"fmt"
	"strings"

	"github.com/SalomanYu/GoPostupiOnline/college/models"

	"github.com/xuri/excelize/v2"
)

const (
	VuzFile = "Data/Vuzes.xlsx"
	ContactFile = "Data/Contact.xlsx"
	SpecializationFile = "Data/Specialization.xlsx"
	ProgramFile = "Data/Program.xlsx"
	ProfessionFile = "Data/Profession.xlsx"
)

var (
	VuzBook 			*excelize.File
	ContactBook 		*excelize.File
	SpecializationBook 	*excelize.File
	ProgramBook 		*excelize.File
	ProfessionBook 		*excelize.File
	VuzRow, SpecRow, ProgramRow, ProfessionsRow int
)


func init(){
	VuzBook = 			excelize.NewFile()
	ContactBook = 		excelize.NewFile()
	SpecializationBook =excelize.NewFile()
	ProgramBook = 		excelize.NewFile()
	ProfessionBook = 	excelize.NewFile()

	defer VuzBook.Close()
	defer ContactBook.Close()
	defer SpecializationBook.Close()
	defer ProgramBook.Close()
	defer ProfessionBook.Close()
	
}


func AddVuz(vuz *models.Vuz){
	VuzRow++
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("A%d", VuzRow), VuzRow)
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("B%d", VuzRow), vuz.VuzId)
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("C%d", VuzRow), vuz.Base.Name)
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("D%d", VuzRow), vuz.Base.Description)
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("E%d", VuzRow), vuz.Base.Direction)
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("F%d", VuzRow), vuz.Base.Cost)
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("G%d", VuzRow), vuz.Base.Scores.PlacesBudget)
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("H%d", VuzRow), vuz.Base.Scores.PlacesPayment)
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("I%d", VuzRow), vuz.Base.Scores.PointsBudget)
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("J%d", VuzRow), vuz.Base.Scores.PointsPayment)
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("K%d", VuzRow), vuz.Base.Image)
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("L%d", VuzRow), vuz.Base.Logo)
	VuzBook.SetCellValue("Sheet1", fmt.Sprintf("M%d", VuzRow), vuz.Base.Url)

	fmt.Println("VuzRow", VuzRow)
	if err := VuzBook.SaveAs(VuzFile); err != nil{
		panic(err)
	}
}

func AddContacts(contact *models.Contacts){
	ContactBook.SetCellValue("Sheet1", fmt.Sprintf("A%d", VuzRow), VuzRow)
	ContactBook.SetCellValue("Sheet1", fmt.Sprintf("B%d", VuzRow), contact.VuzId)
	ContactBook.SetCellValue("Sheet1", fmt.Sprintf("C%d", VuzRow), contact.WebSite)
	ContactBook.SetCellValue("Sheet1", fmt.Sprintf("D%d", VuzRow), contact.Email)
	ContactBook.SetCellValue("Sheet1", fmt.Sprintf("E%d", VuzRow), contact.Phone)
	ContactBook.SetCellValue("Sheet1", fmt.Sprintf("F%d", VuzRow), contact.Address)

	if err := ContactBook.SaveAs(ContactFile); err != nil{
		panic(err)
	}
}

func AddSpecialization(spec *models.Specialization){
	SpecRow++
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("A%d", SpecRow), SpecRow)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("B%d", SpecRow), spec.VuzId)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("C%d", SpecRow), spec.SpecId)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("D%d", SpecRow), spec.Base.Name)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("E%d", SpecRow), spec.Base.Description)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("F%d", SpecRow), spec.Base.Direction)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("G%d", SpecRow), spec.Base.Cost)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("H%d", SpecRow), spec.Base.Scores.PlacesBudget)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("I%d", SpecRow), spec.Base.Scores.PlacesPayment)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("J%d", SpecRow), spec.Base.Scores.PointsBudget)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("K%d", SpecRow), spec.Base.Scores.PointsPayment)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("L%d", SpecRow), spec.Base.Image)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("M%d", SpecRow), spec.Base.Logo)
	SpecializationBook.SetCellValue("Sheet1", fmt.Sprintf("N%d", SpecRow), spec.Base.Url)

	fmt.Println("SpecRow", SpecRow)
	if err := SpecializationBook.SaveAs(SpecializationFile); err != nil{
		panic(err)
	}
}

func AddProgram(program *models.Program){
	ProgramRow++
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("A%d", ProgramRow), ProgramRow)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("B%d", ProgramRow), program.VuzId)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("C%d", ProgramRow), program.SpecId)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("D%d", ProgramRow), program.ProgramId)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("E%d", ProgramRow), program.Base.Name)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("F%d", ProgramRow), strings.Join([]string{program.Base.Description, program.Description}, "\n"))
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("G%d", ProgramRow), program.Base.Direction)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("H%d", ProgramRow), program.Base.Cost)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("I%d", ProgramRow), program.Base.Scores.PlacesBudget)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("J%d", ProgramRow), program.Base.Scores.PlacesPayment)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("K%d", ProgramRow), program.Base.Scores.PointsBudget)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("L%d", ProgramRow), program.Base.Scores.PointsPayment)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("M%d", ProgramRow), program.Base.Image)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("N%d", ProgramRow), program.Base.Logo)
	ProgramBook.SetCellValue("Sheet1", fmt.Sprintf("O%d", ProgramRow), program.Base.Url)

	fmt.Println("Program", ProgramRow)
	if err := ProgramBook.SaveAs(ProgramFile); err != nil{
		panic(err)
	}
}

func AddProfession(prof *models.Profession){
	ProfessionsRow++
	ProfessionBook.SetCellValue("Sheet1", fmt.Sprintf("A%d", ProfessionsRow), ProfessionsRow)
	ProfessionBook.SetCellValue("Sheet1", fmt.Sprintf("B%d", ProfessionsRow), prof.ProgramId)
	ProfessionBook.SetCellValue("Sheet1", fmt.Sprintf("C%d", ProfessionsRow), prof.Name)
	ProfessionBook.SetCellValue("Sheet1", fmt.Sprintf("D%d", ProfessionsRow), prof.Image)

	fmt.Printf("[%d]Profession %d\n", VuzRow, ProfessionsRow)
	if err := ProfessionBook.SaveAs(ProfessionFile); err != nil{
		panic(err)
	}
}