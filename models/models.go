package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Basic struct {
	// BudgetPlaces - string, а не int, т.к не получается вырезать юникод в строке
	ID				primitive.ObjectID 	`bson:"_id"`
	Name 			string				`bson:"name"`
	Url  			string				`bson:"url"`
	Description 	string				`bson:"description"`
	Direction 		string				`bson:"direction"`
	Image     		string				`bson:"image"`
	Logo      		string				`bson:"logo"`
	Cost      		string				`bson:"cost"`
	BudgetScores		float64				`bson:"budget_score"`
	PaymentScores	float64				`bson:"payment_score"`
	BudgetPlaces	string 				`bson:"budget_places"`
	PaymentPlaces	string				`bson:"payment_places"`
}
type Vuz struct {
	ID				primitive.ObjectID	`bson:"_id"`
	VuzId 	string						`bson:"vuz_id"`
	Description   	string				`bson:"description"`
	Base          	Basic				`bson:"base"`
}
type Specialization struct {
	ID				primitive.ObjectID	`bson:"_id"`
	SpecId        	string				`bson:"spec_id"`
	VuzId 	string						`bson:"vuz_id"`
	Description		string				`bson:"description"`
	Base          	Basic				`bson:"base"`
}
type Program struct {
	ID				primitive.ObjectID	`bson:"_id"`
	ProgramId     	string				`bson:"program_id"`
	SpecId        	string				`bson:"spec_id"`
	VuzId 	string						`bson:"vuz_id"`
	Description   	string				`bson:"description"`
	Form          	string				`bson:"form"`
	Exams         	[]string			`bson:"exams"`
	Base          	Basic				`bson:"base"`
}
type Profession struct {
	ID				primitive.ObjectID	`bson:"_id"`
	ProgramId 		string				`bson:"program_id"`
	Name      		string				`bson:"name"`
	Image     		string				`bson:"image"`
}
type Contacts struct {
	ID				primitive.ObjectID	`bson:"_id"`
	VuzId   		string				`bson:"vuz_id"`
	WebSite 		string				`bson:"website"`
	Email   		string				`bson:"email"`
	Phone   		string				`bson:"phone"`
	Address 		string				`bson:"address"`
}
