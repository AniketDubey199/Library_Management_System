package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/AniketDubey199/JWT_auth/db"
	"github.com/AniketDubey199/JWT_auth/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

func Download(route fiber.Router, databse *db.MongoDB){
	route.Get("/" , func (c *fiber.Ctx)error{
		format := c.Query("format" , "json") // json by default
		// userinterface := c.Locals("userID")
		// user := userinterface.(string)

		// userID , err := primitive.ObjectIDFromHex(user)
		// if err != nil {
		// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 		"error": "cannot convert it into primitive object id",
		// 	})
		// }

		filter := bson.M{}

		books := new([]model.Book)
		booksCollection:= databse.Client.Database("librarydb").Collection("book")
		cursor , none := booksCollection.Find(context.Background(), filter)
		if none != nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username not found ",
			})
		}
		defer cursor.Close(context.Background())

		if err := cursor.All(context.Background(), books); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode books: " + err.Error(),
			})
		}

		var filename string
		switch format{
			case "json":
				filename = "books.json"
				// Create a json file
				file , err := os.Create(filename)
				if err != nil{
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Failed to create the JSON file",
					})
				}
				defer file.Close()

				encoder := json.NewEncoder(file)
				encoder.SetIndent("", " ")
				if err := encoder.Encode(books); err != nil{
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Failed to write the JSON file",
					})
				}

			case "csv":
				filename = "books.csv"

				//Create a csv file

				file , err := os.Create(filename)
				if err != nil{
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Failed to create the CSV file",
					})
				}
				defer file.Close()

				writer := csv.NewWriter(file)

				writer.Write([]string{"ID", "Title","Status","Author","Year"})

				for _, book := range *books {
					writer.Write([]string{
						book.ID.Hex(),
						book.Title,
						string(book.Status),
						book.Author,
						fmt.Sprintf("%d", book.Year),
					})
				}

				writer.Flush()
				if err := writer.Error(); err != nil{
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Failed to write the CSV file",
					})
				}

			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Invalid format . Use 'json' or 'csv' ",
				}) 
			}
		
		defer os.Remove(filename)
		return c.Download(filename)
	})
}