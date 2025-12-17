package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AniketDubey199/JWT_auth/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	Collection *mongo.Collection
}
func InitalizeDB() (*MongoDB ,error) {
	ctx , cancel :=context.WithTimeout(context.Background() , 10*time.Second)
	defer cancel()

	var err error
	clienoptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client , err := mongo.Connect(ctx , clienoptions)
	if err!=nil{
		return nil , fmt.Errorf("request cannot be approved %v" , err)
	}
	db_name := os.Getenv("DB_NAME")
	collection := client.Database(db_name).Collection("library")

	log.Println("Connected to Mongodb")
	return &MongoDB{
		Client: client,
		Collection: collection,
	},nil
}

func CreateUser(client *mongo.Client , user *model.User) (error){
	collection := client.Database("librarydb").Collection("library")
	ctx , cancel := context.WithTimeout(context.Background() , 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx , user)
	return err
}

func AddBook(client *mongo.Client , book *model.Book) error{
	collection := client.Database("librarydb").Collection("book")
	ctx , cancel := context.WithTimeout(context.Background() , 5*time.Second)
	defer cancel()

	_,err := collection.InsertOne(ctx,book)
	return err
}
