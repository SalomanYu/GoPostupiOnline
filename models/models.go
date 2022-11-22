package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type BasicInfo struct {
	// BudgetPlaces - string, а не int, т.к не получается вырезать юникод в строке
	ID				primitive.ObjectID 	`bson:"_id"`
	Name 			string				`bson:"name"`
	Url  			string				`bson:"url"`
	Description 	string				`bson:"description"`
	Direction 		string				`bson:"direction"`
	Image     		string				`bson:"image"`
	Logo      		string				`bson:"logo"`
	Cost      		string				`bson:"cost"`
	BudgetScore		float64				`bson:"budget_score"`
	PaymentScore	float64				`bson:"payment_score"`
	BudgetPlaces	string 				`bson:"budget_places"`
	PaymentPlaces	string				`bson:"payment_places"`
}
type InstitutionInfo struct {
	ID				primitive.ObjectID	`bson:"_id"`
	InstitutionId 	string				`bson:"vuz_id"`
	Description   	string				`bson:"description"`
	Base          	BasicInfo			`bson:"base"`
}
type SpecializationInfo struct {
	ID				primitive.ObjectID	`bson:"_id"`
	SpecId        	string				`bson:"spec_id"`
	InstitutionId 	string				`bson:"vuz_id"`
	Description		string				`bson:"description"`
	Base          	BasicInfo			`bson:"base"`
}
type ProgramInfo struct {
	ID				primitive.ObjectID	`bson:"_id"`
	ProgramId     	string				`bson:"program_id"`
	SpecId        	string				`bson:"spec_id"`
	InstitutionId 	string				`bson:"vuz_id"`
	Description   	string				`bson:"description"`
	Form          	string				`bson:"form"`
	Exams         	[]string			`bson:"exams"`
	Base          	BasicInfo			`bson:"base"`
}
type ProfessionInfo struct {
	ID				primitive.ObjectID	`bson:"_id"`
	ProgramId 		string				`bson:"program_id"`
	Name      		string				`bson:"name"`
	Image     		string				`bson:"image"`
}
type ContactsInfo struct {
	ID				primitive.ObjectID	`bson:"_id"`
	VuzId   		string				`bson:"vuz_id"`
	WebSite 		string				`bson:"website"`
	Email   		string				`bson:"email"`
	Phone   		string				`bson:"phone"`
	Address 		string				`bson:"address"`
}
