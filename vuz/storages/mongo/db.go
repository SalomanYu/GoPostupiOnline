package mongo

import (
	"context"
	"fmt"

	"github.com/SalomanYu/GoPostupiOnline/vuz/models"

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
	dbName = "vuzes"
)

func init(){
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	check_err(err)
	
	err = client.Ping(ctx, nil)
	check_err(err)

	collectionVuz = client.Database(dbName).Collection("Vuz")
	collectionSpec = client.Database(dbName).Collection("Specialization")
	collectionProgram = client.Database(dbName).Collection("Program")
	collectionProfession = client.Database(dbName).Collection("Profession")
	collectionContacts = client.Database(dbName).Collection("Contacts")
}

func AddVuz(vuz *models.Vuz) error{
	_, err := collectionVuz.InsertOne(ctx, vuz)
	vuzCount++
	fmt.Printf("%d. Vuz: %s\n", vuzCount, vuz.VuzId)
	return err
}
func AddSpecialization(spec *models.Specialization) error{
	_, err := collectionSpec.InsertOne(ctx, spec)
	return err
}
func AddProgram(program *models.Program) error{
	_, err := collectionProgram.InsertOne(ctx, program)
	return err
}
func AddProfession(profession *models.Profession) error{
	_, err := collectionProfession.InsertOne(ctx, profession)
	return err
}
func AddContacts(contacts *models.Contacts) error{
	_, err := collectionContacts.InsertOne(ctx, contacts)
	return err
}

func check_err(err error){
	if err != nil{
		panic(err)
	}
}