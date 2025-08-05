package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Status string 

const (
	Read Status = "read"
	Reading Status = "reading"
	To_Read Status = "to_read"
)	

type User struct {
	ID primitive.ObjectID `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type Book struct{
	ID primitive.ObjectID `json:"id"`
	Title string `json:"title"`
	Status Status `json:"status"`
	Author string `json:"author"`
	Year int 	`json:"year"`
	UserID int 	`json:"userid"`
}