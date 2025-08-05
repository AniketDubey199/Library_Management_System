package main

import (
	"context"
	"fmt"
	"log"
	"time"

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
	clienoptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client , err := mongo.Connect(ctx , clienoptions)
	if err!=nil{
		return nil , fmt.Errorf("request cannot be approved %v" , err)
	}
	collection := client.Database("librarydb").Collection("library")

	log.Println("Connected to Mongodb")
	return &MongoDB{
		Client: client,
		Collection: collection,
	},nil
}

func CreateUser(client *mongo.Client , user User) (error){
	collection := client.Database("librarydb").Collection("library")
	ctx , cancel := context.WithTimeout(context.Background() , 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx , user)
	return err
}
