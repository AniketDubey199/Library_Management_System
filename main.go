package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {

	db,err := InitalizeDB()
	if err != nil{
		fmt.Print("Error");
	}

	app := fiber.New(fiber.Config{
		AppName: "Library API",
	})
	// Authentication routes
	Authentication(app.Group("/auth") , db)

	app.Listen(":3000")
}