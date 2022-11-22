package mongo

import (
	"context"
	"fmt"

	"github.com/SalomanYu/GoPostupiOnline/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	collectionVuz 			*mongo.Collection
	collectionSpec 			*mongo.Collection
	collectionProgram 		*mongo.Collection
	collectionContacts 		*mongo.Collection
	collectionProfession 	*mongo.Collection
	ctx 				 =  context.TODO()
	vuzCount int

)

func init(){
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	check_err(err)
	
	err = client.Ping(ctx, nil)
	check_err(err)

	collectionVuz = client.Database("PostupiOnline").Collection("Vuz")
	collectionSpec = client.Database("PostupiOnline").Collection("Specialization")
	collectionProgram = client.Database("PostupiOnline").Collection("Program")
	collectionProfession = client.Database("PostupiOnline").Collection("Profession")
	collectionContacts = client.Database("PostupiOnline").Collection("Contacts")
}

func AddVuz(vuz *models.InstitutionInfo) error{
	_, err := collectionVuz.InsertOne(ctx, vuz)
	vuzCount++
	fmt.Printf("%d. Vuz: %s\n", vuzCount, vuz.InstitutionId)
	return err
}
func AddSpecialization(spec *models.SpecializationInfo) error{
	_, err := collectionSpec.InsertOne(ctx, spec)
	return err
}
func AddProgram(program *models.ProgramInfo) error{
	_, err := collectionProgram.InsertOne(ctx, program)
	return err
}
func AddProfession(profession *models.ProfessionInfo) error{
	_, err := collectionProfession.InsertOne(ctx, profession)
	return err
}
func AddContacts(contacts *models.ContactsInfo) error{
	_, err := collectionContacts.InsertOne(ctx, contacts)
	return err
}

func check_err(err error){
	if err != nil{
		panic(err)
	}
}