package main

import (
	"fmt"
	"os"
	"log"

	"github.com/AniketDubey199/JWT_auth/auth"
	"github.com/AniketDubey199/JWT_auth/auth/authmiddleware"
	"github.com/AniketDubey199/JWT_auth/book"
	"github.com/AniketDubey199/JWT_auth/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	data,err := db.InitalizeDB()
	if err != nil{
		fmt.Print("Error");
	}

	app := fiber.New(fiber.Config{
		AppName: "Library API",
	})
	// Authentication routes
	auth.Authentication(app.Group("/auth") , data)
	
	Download(app.Group("/download"),data)
	// Verify the JWT . If Valid , it will set the UserID 
	protected := app.Use(authmiddleware.AuthMiddleware(data))
	
	// Book handler
	book.Bookhandler(protected.Group("/book"),data)
	
	// Download kar rhe hai 
	// Download(protected.Group("/download"),data)
	port:= os.Getenv("PORT")
	if port == ""{
		port = ":3000"
	}
	
	app.Listen("0.0.0.0:"+port)
}	